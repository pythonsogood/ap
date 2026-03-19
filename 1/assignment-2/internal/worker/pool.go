package worker

import (
	"context"
	"math/rand"
	"sync"
	"time"

	"github.com/SultanYakupov/Assignment2/internal/model"
	"github.com/SultanYakupov/Assignment2/internal/store"
)

func processTask(workerID int, id string, repo *store.Repository[string, model.Task]) {
	t, exists := repo.Get(id)
	if !exists {
		return
	}

	// simulate work
	time.Sleep(time.Duration(500+rand.Intn(2000)) * time.Millisecond)

	t, exists = repo.Get(id)
	if !exists {
		return
	}

	t.Status = model.StatusDone
	repo.Set(id, t)
}

func StartWorkerPool(ctx context.Context, workers int, queue <-chan string, repo *store.Repository[string, model.Task]) *sync.WaitGroup {
	var wg sync.WaitGroup
	wg.Add(workers)

	for i := range workers {
		go func(workerID int) {
			defer wg.Done()

			for {
				select {
				case <-ctx.Done():
					return

				case id, ok := <-queue:
					if !ok {
						return
					}

					processTask(workerID, id, repo)
				}
			}
		}(i + 1)
	}

	return &wg
}
