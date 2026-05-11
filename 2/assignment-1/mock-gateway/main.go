package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

type notifyRequest struct {
	IdempotencyKey string `json:"idempotency_key"`
	Channel        string `json:"channel"`
	Recipient      string `json:"recipient"`
	Message        string `json:"message"`
}

type gateway struct {
	mu   sync.Mutex
	seen map[string]struct{}
	rng  *rand.Rand
}

func newGateway() *gateway {
	return &gateway{
		seen: make(map[string]struct{}),
		rng:  rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (g *gateway) notifyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req notifyRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid json"}`, http.StatusBadRequest)
		return
	}

	reqLog := map[string]any{
		"time":            time.Now().UTC().Format(time.RFC3339),
		"path":            r.URL.Path,
		"idempotency_key": req.IdempotencyKey,
		"channel":         req.Channel,
		"recipient":       req.Recipient,
		"message":         req.Message,
	}

	if b, err := json.Marshal(reqLog); err == nil {
		log.Println(string(b))
	}

	if g.rng.Intn(100) < 20 {
		w.Header().Set("Content-Type", "application/json")

		w.WriteHeader(http.StatusServiceUnavailable)

		_, _ = w.Write([]byte(`{"status":"temporary_failure"}`))

		return
	}

	g.mu.Lock()

	_, exists := g.seen[req.IdempotencyKey]

	if !exists {
		g.seen[req.IdempotencyKey] = struct{}{}
	}

	g.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if exists {
		_, _ = w.Write([]byte(`{"status":"duplicate"}`))
		return
	}

	_, _ = w.Write([]byte(`{"status":"accepted"}`))
}

func main() {
	log.SetFlags(0)

	g := newGateway()

	mux := http.NewServeMux()

	mux.HandleFunc("/notify", g.notifyHandler)

	if err := http.ListenAndServe(":8080", mux); err != nil {
		panic(err.Error())
	}
}
