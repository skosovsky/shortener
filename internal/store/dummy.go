package store

import "shortener/internal/service"

type DummyStore struct{}

func NewDummyStore() *DummyStore {
	return &DummyStore{}
}

func (m *DummyStore) Add(_ service.Site) {
}

func (m *DummyStore) Get(_ string) (service.Site, error) {
	return service.Site{}, nil //nolint:exhaustruct // empty
}
