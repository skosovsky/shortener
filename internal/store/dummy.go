package store

import (
	"shortener/internal/model"
)

type DummyStore struct{}

func NewDummyStore() (*DummyStore, error) {
	return &DummyStore{}, nil
}

func (m *DummyStore) Add(_ model.Site) bool {
	return true
}

func (m *DummyStore) Get(_ string) (model.Site, bool) {
	return model.Site{}, true //nolint:exhaustruct // empty
}

func (m *DummyStore) Update(_ string, _ model.Site) bool {
	return true
}

func (m *DummyStore) Delete(_ string) bool {
	return true
}
