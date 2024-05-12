package store //nolint:dupl // fake

import (
	"sync"

	"shortener/internal/model"
)

type FakeStore struct {
	memory map[string]model.Site
	mu     sync.RWMutex
}

func NewFakeStore() (*FakeStore, error) {
	return &FakeStore{
		memory: map[string]model.Site{},
		mu:     sync.RWMutex{},
	}, nil
}

func (m *FakeStore) Add(site model.Site) bool {
	m.mu.Lock()
	m.memory[site.ID] = site
	m.mu.Unlock()

	return true
}

func (m *FakeStore) Get(id string) (model.Site, bool) {
	m.mu.RLock()
	site, ok := m.memory[id]
	m.mu.RUnlock()

	return site, ok
}

func (m *FakeStore) Update(id string, site model.Site) bool {
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

func (m *FakeStore) Delete(id string) bool {
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
