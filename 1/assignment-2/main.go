package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/SultanYakupov/Assignment2/internal/api"
	"github.com/SultanYakupov/Assignment2/internal/model"
	"github.com/SultanYakupov/Assignment2/internal/store"
	"github.com/SultanYakupov/Assignment2/internal/worker"
)

const ADDR = ":8000"
const WORKERS = 4
const QUEUE_SIZE = 16
const MONITOR_INTERVAL = 5
const SHUTDOWN_HANDLERS_TIMEOUT = 5
const SHUTDOWN_SERVER_TIMEOUT = 5

func monitorLog(server_api *api.Server) {
	stats := server_api.GetTasksStats()

	log.Printf("[Monitor] Submitted: %d, Completed: %d, InProgress: %d\n", stats.Submitted, stats.Completed, stats.InProgress)
}

func main() {
	queue := make(chan string, QUEUE_SIZE)
	repo := store.NewRepository[string, model.Task]()

	worker_ctx, worker_cancel := context.WithCancel(context.Background())
	worker_wg := worker.StartWorkerPool(worker_ctx, WORKERS, queue, repo)

	server_api := api.NewServer(repo, queue)

	server_mux := http.NewServeMux()
	server_mux.HandleFunc("/tasks", server_api.HandleTasks)
	server_mux.HandleFunc("/tasks/{id}", server_api.HandleTasksById)
	server_mux.HandleFunc("/stats", server_api.HandleTasksStats)

	server := &http.Server{
		Addr:    ADDR,
		Handler: server_mux,
	}

	server_errors := make(chan error)
	go func() {
		log.Printf("Server running at %s", ADDR)
		server_errors <- server.ListenAndServe()
	}()

	monitor_stop := make(chan struct{})
	var monitor_wg sync.WaitGroup
	monitor_wg.Go(func() {
		ticker := time.NewTicker(MONITOR_INTERVAL * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				monitorLog(server_api)

			case <-monitor_stop:
				return
			}
		}
	})

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-server_errors:
		log.Printf("Server error: %v\n", err)

	case signal := <-quit:
		log.Printf("Shutdown signal: %v\n", signal)

		server_api.SetAccepting(false)
		server_api.WaitForHandlersDone(SHUTDOWN_HANDLERS_TIMEOUT * time.Second)

		ctx_shutdown, shutdown_cancel := context.WithTimeout(context.Background(), SHUTDOWN_SERVER_TIMEOUT*time.Second)
		defer shutdown_cancel()

		if err := server.Shutdown(ctx_shutdown); err != nil {
			log.Printf("Server shutdown error: %v\n", err)
		}

		close(queue)

		close(monitor_stop)
		monitor_wg.Wait()

		worker_cancel()
		worker_wg.Wait()

		log.Println("Shutdown done")
	}
}
