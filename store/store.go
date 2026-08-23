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

// New creates an independent store. The seed slice is copied so that writes to
// one store (or to the list returned by List) never bleed into another store or
// back into the shared package-level seed.
func New() *Store {
	items := make([]domain.Alert, len(seedAlerts))
	copy(items, seedAlerts)
	return &Store{items: items}
}

// List returns a defensive copy of the alerts. Callers may mutate the returned
// slice and its elements without affecting the store's internal state or any
// other caller that previously called List.
func (s *Store) List() []domain.Alert {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]domain.Alert, len(s.items))
	copy(out, s.items)
	return out
}

// UpdateStatus sets the status of the alert with the given id. It mutates the
// store in place by index so the change is persisted and reflected by List and
// by the returned alert.
func (s *Store) UpdateStatus(id, v string) (domain.Alert, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := range s.items {
		if s.items[i].ID == id {
			s.items[i].Status = v
			return s.items[i], nil
		}
	}
	return domain.Alert{}, ErrNotFound
}
