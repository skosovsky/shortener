package service

import (
	"database/sql"
	"errors"
	"fmt"
	"net/url"

	"shortener/config"
	"shortener/internal/log"
)

type Site struct {
	ID        string `json:"id"`
	Link      string `json:"link"`
	ShortLink string `json:"shortLink"`
}

var (
	ErrSiteNotAdded = errors.New("site not added")
	ErrSiteNotFound = errors.New("site not found")
)

type Store interface {
	Add(Site) error
	Get(string) (Site, error)
	Close()
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

	err = s.store.Add(site)
	if err != nil {
		return Site{}, fmt.Errorf("site not add, %w - %w", err, ErrSiteNotAdded)
	}

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

func (s Shortener) Ping() error {
	db, err := sql.Open("pgx", s.config.Store.DatabaseDSN)
	if err != nil {
		return fmt.Errorf("could not connect to database: %w", err)
	}

	defer func(db *sql.DB) {
		err = db.Close()
		if err != nil {
			log.Error("could not close database connection",
				log.ErrAttr(err))
		}
	}(db)

	err = db.Ping()
	if err != nil {
		return fmt.Errorf("could not ping database: %w", err)
	}

	return nil
}
