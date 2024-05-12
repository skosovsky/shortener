package service

import (
	"errors"
	"fmt"
	"net/url"

	"shortener/internal/model"
	"shortener/internal/store"
	log "shortener/pkg/logger"
)

var (
	ErrSiteNotAdded = errors.New("site not added")
	ErrSiteNotFound = errors.New("site not found")
)

type Shortener struct {
	store store.Store
}

func NewSiteService(store store.Store) Shortener {
	return Shortener{store: store}
}

func (s Shortener) Add(link string) (model.Site, error) {
	_, err := url.Parse(link)
	if err != nil {
		return model.Site{}, fmt.Errorf("invalid link URL: %w", err)
	}

	site := SiteGenerate(link)

	if ok := s.store.Add(site); !ok {
		return model.Site{}, ErrSiteNotAdded
	}

	log.Info("site added", log.StringAttr("site", fmt.Sprint(site)))

	return site, nil
}

func (s Shortener) Get(id string) (model.Site, error) {
	site, ok := s.store.Get(id)
	if !ok {
		return model.Site{}, ErrSiteNotFound
	}

	log.Info("site returned", log.StringAttr("site", fmt.Sprint(site)))

	return site, nil
}
