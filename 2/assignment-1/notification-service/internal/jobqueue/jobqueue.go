package jobqueue

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	idempotencyTTL                   = time.Duration(24) * time.Hour
	maxAttempts                      = 3
	notificationIdempotencyKeyFormat = "notification:idempotency:%s"
)

var backoffs = []time.Duration{
	1 * time.Second,
	2 * time.Second,
	4 * time.Second,
}

type Job struct {
	IdempotencyKey string `json:"idempotency_key"`
	AppointmentID  string `json:"appointment_id"`
	DoctorID       string `json:"doctor_id"`
	OccurredAt     string `json:"occurred_at"`
	Channel        string `json:"channel"`
	Recipient      string `json:"recipient"`
	Message        string `json:"message"`
}

type Queue struct {
	jobs       chan Job
	workers    int
	gatewayUrl string
	client     *http.Client
	rdb        *redis.Client
	wg         sync.WaitGroup
}

func NewQueue(rdb *redis.Client, gatewayUrl string, workers int, buffer int) *Queue {
	if workers <= 0 {
		workers = 3
	}

	if buffer <= 0 {
		buffer = 64
	}

	return &Queue{
		jobs:       make(chan Job, buffer),
		workers:    workers,
		gatewayUrl: gatewayUrl,
		client: &http.Client{
			Timeout: time.Duration(5) * time.Second,
		},
		rdb: rdb,
	}
}

func BuildIdempotencyKey(eventType string, id string, occurredAt string) string {
	sum := sha256.Sum256([]byte(eventType + id + occurredAt))
	return hex.EncodeToString(sum[:])
}

func BuildJob(appointmentID string, doctorID string, occurredAt string, idempotencyKey string) Job {
	return Job{
		IdempotencyKey: idempotencyKey,
		AppointmentID:  appointmentID,
		DoctorID:       doctorID,
		OccurredAt:     occurredAt,
		Channel:        "email",
		Recipient:      "patient@clinic.kz",
		Message:        fmt.Sprintf("Your appointment %s with doctor %s is complete.", appointmentID, doctorID),
	}
}

func (q *Queue) Start(ctx context.Context) {
	for i := 0; i < q.workers; i++ {
		q.wg.Add(1)

		go q.worker(ctx)
	}
}

func (q *Queue) Stop() {
	close(q.jobs)

	q.wg.Wait()
}

func (q *Queue) Enqueue(ctx context.Context, job Job) {
	if q.isDone(ctx, job.IdempotencyKey) {
		return
	}

	select {
	case q.jobs <- job:
		q.logState("info", job.IdempotencyKey, 1, "enqueued", "")
	default:
		q.jobs <- job
		q.logState("warn", job.IdempotencyKey, 1, "enqueued", "queue was full; enqueue blocked")
	}
}

func (q *Queue) worker(ctx context.Context) {
	defer q.wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case job, ok := <-q.jobs:
			if !ok {
				return
			}

			q.processJob(ctx, job)
		}
	}
}

func (q *Queue) processJob(ctx context.Context, job Job) {
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		q.logState("info", job.IdempotencyKey, attempt, "processing", "")

		err := q.callGateway(ctx, job)

		if err == nil {
			q.markDone(ctx, job.IdempotencyKey)
			q.logState("info", job.IdempotencyKey, attempt, "success", "")
			return
		}

		if attempt < maxAttempts {
			q.logState("warn", job.IdempotencyKey, attempt, "retry", err.Error())
			time.Sleep(backoffs[attempt-1])
			continue
		}

		q.logDeadLetter(job.IdempotencyKey, attempt, err.Error())
		return
	}
}

func (q *Queue) callGateway(ctx context.Context, job Job) error {
	b, err := json.Marshal(job)

	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/notify", q.gatewayUrl), bytes.NewReader(b))

	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := q.client.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	io.Copy(io.Discard, resp.Body)

	if resp.StatusCode == http.StatusOK {
		return nil
	}

	if resp.StatusCode == http.StatusServiceUnavailable {
		return errors.New("gateway returned 503")
	}

	return fmt.Errorf("gateway returned %d", resp.StatusCode)
}

func (q *Queue) isDone(ctx context.Context, idempotencyKey string) bool {
	if q.rdb == nil {
		return false
	}

	val, err := q.rdb.Get(ctx, fmt.Sprintf(notificationIdempotencyKeyFormat, idempotencyKey)).Result()

	if err == redis.Nil {
		return false
	}

	if err != nil {
		log.Printf("[jobqueue] idempotency read failed %s: %v", idempotencyKey, err.Error())
		return false
	}

	return val == "done"
}

func (q *Queue) markDone(ctx context.Context, idempotencyKey string) {
	if q.rdb == nil {
		return
	}

	err := q.rdb.Set(ctx, fmt.Sprintf(notificationIdempotencyKeyFormat, idempotencyKey), "done", idempotencyTTL).Err()

	if err != nil {
		log.Printf("[jobqueue] idempotency write failed %s: %v", idempotencyKey, err.Error())
	}
}

func (q *Queue) logState(level, jobID string, attempt int, statusText, errMsg string) {
	line := map[string]any{
		"time":    time.Now().UTC().Format(time.RFC3339),
		"level":   level,
		"job_id":  jobID,
		"attempt": attempt,
		"status":  statusText,
	}

	if errMsg != "" {
		line["error"] = errMsg
	}

	b, err := json.Marshal(line)

	if err != nil {
		log.Printf("[jobqueue] marshal log failed job_id=%s: %v", jobID, err)
		return
	}

	log.Println(string(b))
}

func (q *Queue) logDeadLetter(jobID string, attempt int, err string) {
	line := map[string]any{
		"time":    time.Now().UTC().Format(time.RFC3339),
		"level":   "error",
		"job_id":  jobID,
		"attempt": attempt,
		"status":  "dead_letter",
		"error":   err,
	}

	b, _ := json.Marshal(line)

	fmt.Fprintln(os.Stderr, string(b))
}
