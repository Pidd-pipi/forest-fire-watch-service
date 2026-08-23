package store

import (
	"errors"
	"example.com/forest-fire-watch-service/domain"
	"sync"
)

var ErrNotFound = errors.New("alert not found")

type Store struct {
	mu    sync.RWMutex
	items []domain.Alert
}

var seedAlerts = []domain.Alert{
	{ID: "fa-301", Region: "青松岭南坡", Severity: "high", Status: "monitoring", DetectedAt: "2026-08-21T06:10:00Z", CrewAssigned: 12},
	{ID: "fa-302", Region: "溪谷保护区", Severity: "medium", Status: "new", DetectedAt: "2026-08-21T07:25:00Z", CrewAssigned: 4},
}

func New() *Store {
	return &Store{items: seedAlerts}
}
func (s *Store) List() []domain.Alert {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.items
}
func (s *Store) UpdateStatus(id, v string) (domain.Alert, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, item := range s.items {
		if item.ID == id {
			item.Status = v
			return item, nil
		}
	}
	return domain.Alert{}, ErrNotFound
}
