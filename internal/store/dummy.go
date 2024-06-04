package store

import "shortener/internal/service"

type DummyStore struct{}

func NewDummyStore() *DummyStore {
	return &DummyStore{}
}

func (*DummyStore) Add(_ service.Site) error {
	return nil
}

func (*DummyStore) Get(_ string) (service.Site, error) {
	return service.Site{}, nil //nolint:exhaustruct // empty
}

func (*DummyStore) Close() {}
