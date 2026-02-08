package store

import (
	"awesomeProject1/internal/models"
	"sort"
	"sync"
)

type Store struct {
	mu     sync.Mutex
	nextID int
	tasks  map[int]models.Task
}

func New() *Store {
	return &Store{
		nextID: 1,
		tasks:  make(map[int]models.Task),
	}
}

func (s *Store) Create(title string) models.Task {
	s.mu.Lock()
	defer s.mu.Unlock()

	t := models.Task{
		ID:    s.nextID,
		Title: title,
		Done:  false,
	}
	s.tasks[t.ID] = t
	s.nextID++
	return t
}

func (s *Store) Get(id int) (models.Task, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	t, ok := s.tasks[id]
	return t, ok
}

func (s *Store) List() []models.Task {
	s.mu.Lock()
	defer s.mu.Unlock()

	out := make([]models.Task, 0, len(s.tasks))
	for _, t := range s.tasks {
		out = append(out, t)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}

func (s *Store) UpdateDone(id int, done bool) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	t, ok := s.tasks[id]
	if !ok {
		return false
	}
	t.Done = done
	s.tasks[id] = t
	return true
}
