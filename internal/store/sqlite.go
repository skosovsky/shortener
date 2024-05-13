package store

import (
	"database/sql"
	"fmt"

	"shortener/internal/model"
	log "shortener/pkg/logger"
)

type SQLStore struct {
	*sql.DB
}

func NewDBStore() (*SQLStore, error) {
	db, err := sql.Open("sqlite", "./data/sites.db")
	if err != nil {
		return nil, fmt.Errorf("opening sqlite db: %w", err)
	}

	return &SQLStore{db}, nil
}

func (s SQLStore) CloseDBStore() {
	err := s.Close()
	if err != nil {
		log.Error("closing sqlite db", log.ErrAttr(err))
	}
}

func (s SQLStore) Add(site model.Site) bool {
	if _, ok := s.Get(site.ID); ok {
		return false
	}

	_, err := s.Exec("INSERT INTO sites (id, link, short_link) VALUES (:id, :link, :shortLink)",
		sql.Named("id", site.Link),
		sql.Named("link", site.Link),
		sql.Named("shortLink", site.ShortLink))
	if err != nil {
		log.Error("adding site", log.ErrAttr(err))

		return false
	}

	return true
}

func (s SQLStore) Get(id string) (model.Site, bool) {
	var site model.Site

	row := s.QueryRow("SELECT id, link, short_link FROM sites WHERE id = :id",
		sql.Named("id", id))

	err := row.Scan(&site.ID, &site.Link, &site.ShortLink)
	if err != nil {
		log.Error("querying site", log.ErrAttr(err))

		return model.Site{}, false //nolint:exhaustruct // empty
	}

	return site, true
}

func (s SQLStore) Update(id string, site model.Site) bool {
	_, err := s.Exec("UPDATE sites SET id = :id, link = :link, short_link = :shortLink WHERE id = :id",
		sql.Named("id", id),
		sql.Named("link", site.Link),
		sql.Named("shortLink", site.ShortLink))
	if err != nil {
		log.Error("updating site", log.ErrAttr(err))

		return false
	}

	return true
}

func (s SQLStore) Delete(id string) bool {
	_, err := s.Exec("DELETE FROM sites WHERE id = :id",
		sql.Named("id", id))
	if err != nil {
		log.Error("deleting site", log.ErrAttr(err))

		return false
	}

	return true
}
