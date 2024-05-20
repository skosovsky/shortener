package store

import (
	"shortener/internal/model"
)

type DummyStore struct{}

func NewDummyStore() *DummyStore {
	return &DummyStore{}
}

func (m *DummyStore) Add(_ model.Site) {
}

func (m *DummyStore) Get(_ string) (model.Site, error) {
	return model.Site{}, nil //nolint:exhaustruct // empty
}
