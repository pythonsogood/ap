package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/pythonsogood/ap-assignment1/notification/cmd/notification/config"
	"github.com/pythonsogood/ap-assignment1/notification/internal/jobqueue"
	"github.com/pythonsogood/ap-assignment1/notification/internal/subscriber"
	"github.com/redis/go-redis/v9"
)

type logWriter struct {
	time_format string
	file_path   string
}

func connectNATSWithRetry(connectionURL string, maxAttempts int, initialBackoff time.Duration) (*nats.Conn, error) {
	if maxAttempts <= 0 {
		maxAttempts = 1
	}

	if initialBackoff <= 0 {
		initialBackoff = time.Second
	}

	backoff := initialBackoff

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		nc, err := nats.Connect(connectionURL)

		if err == nil {
			return nc, nil
		}

		if attempt == maxAttempts {
			return nil, fmt.Errorf("failed to connect to nats after %d attempts: %w", maxAttempts, err)
		}

		log.Printf("nats unavailable (attempt %d/%d): %v; retrying in %s\n", attempt, maxAttempts, err, backoff)
		time.Sleep(backoff)
		backoff *= 2
	}

	return nil, errors.New("failed to connect to nats")
}

func (lw *logWriter) Write(bs []byte) (int, error) {
	text := string(bs)

	f, err := os.OpenFile(lw.file_path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)

	if err == nil {
		defer f.Close()

		f.WriteString(text)
	}

	return fmt.Print(text)
}

func main() {
	conf, err := config.NewDefaultConfig()

	if err != nil {
		panic(err.Error())
	}

	log.SetFlags(0)
	log.SetOutput(&logWriter{
		time_format: time.RFC3339,
		file_path:   conf.Log.File,
	})

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	opts, err := redis.ParseURL(conf.RedisUrl)

	if err != nil {
		panic(err.Error())
	}

	rdb := redis.NewClient(opts)

	queue := jobqueue.NewQueue(rdb, conf.JobQueue.GatewayUrl, int(conf.JobQueue.WorkerPoolSize), 0)
	queue.Start(ctx)
	defer queue.Stop()

	var event_subscriber subscriber.EventSubscriber
	var nc *nats.Conn

	switch conf.MessageBroker.Type {
	case config.MessageBrokerTypeNATS:
		nc, err = connectNATSWithRetry(conf.MessageBroker.Nats.ConnectionUrl, 5, time.Second)

		if err != nil {
			log.Fatalf("broker unavailable at startup: %s", err.Error())
		}

		event_subscriber = subscriber.NewNATSEventSubscriber(nc, queue)
	default:
		panic("Unsupported message broker type!")
	}

	for _, subj := range conf.MessageBroker.LoggedSubjects {
		if err := event_subscriber.SubscribeNotification(subj); err != nil {
			panic(err.Error())
		}
	}

	<-ctx.Done()
	log.Println("received shutdown signal")

	if nc != nil {
		if err := nc.Drain(); err != nil {
			log.Printf("failed to drain nats connection: %s\n", err.Error())
			nc.Close()
		}
	}

	log.Println("notification service stopped")
}
