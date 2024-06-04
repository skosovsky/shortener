package shortener

import (
	"context"
	"fmt"

	"shortener/config"
	"shortener/internal/service"
	"shortener/internal/store"
)

func Run(cfg config.Config) error {
	var db service.Store
	var err error

	if cfg.Store.FileStoragePath == "" {
		db = store.NewMemoryStore()
	} else {
		db, err = store.NewFileStore(cfg.Store.FileStoragePath)
		if err != nil {
			return fmt.Errorf("create file store: %w", err)
		}

		defer db.Close()
	}

	generator := service.NewIDGenerator()

	shortener := service.NewService(db, cfg, generator)

	handler := NewHandler(shortener)

	if err = RunServer(context.Background(), handler, cfg); err != nil {
		return fmt.Errorf("run server: %w", err)
	}

	return nil
}
