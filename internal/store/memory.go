package store //nolint:dupl // false positive

import (
	"sync"

	"shortener/internal/model"
)

type MemoryStore struct {
	memory map[string]model.Site
	mu     sync.RWMutex
}

func NewMemoryStore() (*MemoryStore, error) {
	return &MemoryStore{
		memory: map[string]model.Site{},
		mu:     sync.RWMutex{},
	}, nil
}

func (m *MemoryStore) Add(site model.Site) bool {
	m.mu.Lock()
	m.memory[site.ID] = site
	m.mu.Unlock()

	return true
}

func (m *MemoryStore) Get(id string) (model.Site, bool) {
	m.mu.RLock()
	site, ok := m.memory[id]
	m.mu.RUnlock()

	return site, ok
}

func (m *MemoryStore) Update(id string, site model.Site) bool {
	m.mu.RLock()
	_, ok := m.memory[id]
	m.mu.RUnlock()

	if !ok {
		return false
	}

	m.mu.Lock()
	m.memory[site.ID] = site
	m.mu.Unlock()

	return true
}

func (m *MemoryStore) Delete(id string) bool {
	m.mu.RLock()
	_, ok := m.memory[id]
	m.mu.RUnlock()

	if !ok {
		return false
	}

	m.mu.Lock()
	delete(m.memory, id)
	m.mu.Unlock()

	return true
}
