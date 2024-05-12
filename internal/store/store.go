package store

import (
	"shortener/internal/model"
)

//go:generate mockgen -source store.go -destination=mock.go -package=store
type Store interface {
	Add(site model.Site) bool
	Get(id string) (model.Site, bool)
	Update(id string, site model.Site) bool
	Delete(id string) bool
}
