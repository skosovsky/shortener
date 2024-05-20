package store

import (
	"errors"
	"sync"

	"shortener/internal/model"
)

var ErrSiteNotFound = errors.New("site not found")

type MemoryStore struct {
	memory map[string]model.Site
	mu     sync.Mutex
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		memory: map[string]model.Site{},
		mu:     sync.Mutex{},
	}
}

func (m *MemoryStore) Add(site model.Site) {
	m.mu.Lock()
	m.memory[site.ID] = site
	m.mu.Unlock()
}

func (m *MemoryStore) Get(id string) (model.Site, error) {
	m.mu.Lock()
	site, ok := m.memory[id]
	m.mu.Unlock()

	if !ok {
		return model.Site{}, ErrSiteNotFound
	}

	return site, nil
}
