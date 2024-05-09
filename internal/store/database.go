package store

import "shortener/internal/model"

type Database interface {
	Add(site model.Site) bool
	Get(id string) (model.Site, bool)
	Update(id string, site model.Site) bool
	Delete(id string) bool
}
