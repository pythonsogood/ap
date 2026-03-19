package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/SultanYakupov/Assignment2/internal/model"
)

func (s *Server) GetAllTasks() []model.Task {
	var tasks []model.Task

	for _, task := range s.repo.GetAll() {
		tasks = append(tasks, task)
	}

	return tasks
}

func (s *Server) GetTasksStats() model.TasksStats {
	submitted := 0
	completed := 0
	in_progress := 0

	for _, task := range s.repo.GetAll() {
		switch task.Status {
		case model.StatusPending:
			in_progress++

		case model.StatusDone:
			completed++
		}

		submitted++
	}

	return model.TasksStats{
		Submitted:  submitted,
		Completed:  completed,
		InProgress: in_progress,
	}
}

func (s *Server) HandleTasks(w http.ResponseWriter, r *http.Request) {
	if !s.IsAccepting() {
		http.Error(w, "server is not accepting responses", http.StatusServiceUnavailable)
		return
	}

	s.IncrementActiveHandler()
	defer s.DecrementActiveHandler()

	switch r.Method {
	case http.MethodGet:
		tasks := s.GetAllTasks()

		// empty slice encodes as null
		if len(tasks) > 0 {
			err := json.NewEncoder(w).Encode(model.TasksToJsonArrayModel(tasks))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			fmt.Fprint(w, "[]")
		}

	case http.MethodPost:
		var request model.TasksPostRequestJson

		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		task := model.NewTask(request.Payload)

		s.repo.Set(task.Id, *task)

		select {
		case s.queue <- task.Id: // queue

		default:
			s.queue <- task.Id // if queue is full, try to enqueue again
		}

		w.WriteHeader(http.StatusCreated)

		err = json.NewEncoder(w).Encode(task.ToJsonModel())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (s *Server) HandleTasksById(w http.ResponseWriter, r *http.Request) {
	if !s.IsAccepting() {
		http.Error(w, "server is not accepting responses", http.StatusServiceUnavailable)
		return
	}

	s.IncrementActiveHandler()
	defer s.DecrementActiveHandler()

	switch r.Method {
	case http.MethodGet:
		task_id := r.PathValue("id")

		task, exists := s.repo.Get(task_id)
		if !exists {
			http.Error(w, "", http.StatusNotFound)
			return
		}

		err := json.NewEncoder(w).Encode(task)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (s *Server) HandleTasksStats(w http.ResponseWriter, r *http.Request) {
	if !s.IsAccepting() {
		http.Error(w, "server is not accepting responses", http.StatusServiceUnavailable)
		return
	}

	s.IncrementActiveHandler()
	defer s.DecrementActiveHandler()

	switch r.Method {
	case http.MethodGet:
		stats := s.GetTasksStats()

		err := json.NewEncoder(w).Encode(stats.ToJsonModel())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
