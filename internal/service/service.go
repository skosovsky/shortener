package service

import (
	"errors"
	"fmt"
	"net/url"

	"shortener/config"
	"shortener/internal/log"
)

type Site struct {
	ID        string
	Link      string
	ShortLink string
}

var (
	ErrSiteNotAdded = errors.New("site not added")
	ErrSiteNotFound = errors.New("site not found")
)

type Store interface {
	Add(Site)
	Get(string) (Site, error)
}

type Generator interface {
	Generate(domain string, link string) Site
}

type Shortener struct {
	store     Store
	config    config.Config
	generator Generator
}

func NewService(store Store, config config.Config, generator Generator) Shortener {
	return Shortener{
		store:     store,
		config:    config,
		generator: generator,
	}
}

func (s Shortener) Add(link string) (Site, error) {
	_, err := url.Parse(link)
	if err != nil {
		return Site{}, fmt.Errorf("invalid link URL: %w, %w", err, ErrSiteNotAdded)
	}

	site := s.generator.Generate(s.config.Shortener.Domain, link)

	s.store.Add(site)

	log.Debug("site added",
		log.StringAttr("site", fmt.Sprint(site)))

	return site, nil
}

func (s Shortener) Get(id string) (Site, error) {
	site, err := s.store.Get(id)
	if err != nil {
		return Site{}, ErrSiteNotFound
	}

	log.Debug("site returned",
		log.StringAttr("site", fmt.Sprint(site)))

	return site, nil
}
