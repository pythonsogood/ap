package api

import (
	"sync"
	"time"

	"github.com/SultanYakupov/Assignment2/internal/model"
	"github.com/SultanYakupov/Assignment2/internal/store"
)

type Server struct {
	repo               *store.Repository[string, model.Task]
	queue              chan string
	active_handlers_mu sync.RWMutex
	active_handlers    uint
	accepting          bool
}

func NewServer(repo *store.Repository[string, model.Task], queue chan string) *Server {
	return &Server{
		repo:            repo,
		queue:           queue,
		active_handlers: 0,
		accepting:       true,
	}
}

func (s *Server) SetAccepting(accepting bool) {
	s.accepting = accepting
}

func (s *Server) IsAccepting() bool {
	return s.accepting
}

func (s *Server) ActiveHandlers() uint {
	s.active_handlers_mu.RLock()
	defer s.active_handlers_mu.RUnlock()

	return s.active_handlers
}

func (s *Server) IncrementActiveHandler() {
	s.active_handlers_mu.Lock()
	defer s.active_handlers_mu.Unlock()

	s.active_handlers++
}

func (s *Server) DecrementActiveHandler() {
	s.active_handlers_mu.Lock()
	defer s.active_handlers_mu.Unlock()

	s.active_handlers--
}

func (s *Server) WaitForHandlersDone(timeout time.Duration) {
	deadline := time.Now().Add(timeout)

	for {
		if s.ActiveHandlers() == 0 {
			return
		}

		if time.Now().After(deadline) {
			return
		}

		time.Sleep(50 * time.Millisecond)
	}
}
