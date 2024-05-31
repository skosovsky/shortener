package store

import (
	"errors"
	"sync"

	"shortener/internal/service"
)

var ErrSiteNotFound = errors.New("site not found")

type MemoryStore struct {
	memory map[string]service.Site
	mu     sync.Mutex
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		memory: map[string]service.Site{},
		mu:     sync.Mutex{},
	}
}

func (m *MemoryStore) Add(site service.Site) error {
	m.mu.Lock()
	m.memory[site.ID] = site
	m.mu.Unlock()

	return nil
}

func (m *MemoryStore) Get(id string) (service.Site, error) {
	m.mu.Lock()
	site, ok := m.memory[id]
	m.mu.Unlock()

	if !ok {
		return service.Site{}, ErrSiteNotFound
	}

	return site, nil
}

func (*MemoryStore) Close() {}
