package shortener

import (
	"context"

	"shortener/config"
	log "shortener/internal/logger"
	"shortener/internal/service"
	"shortener/internal/store"
)

func Run(cfg config.Config) {
	db := store.NewMemoryStore()

	generator := service.NewIDGenerator()

	shortener := service.NewService(db, cfg, generator)

	handler := NewHandler(shortener)

	if err := RunServer(context.Background(), handler, cfg); err != nil {
		log.Fatal("run", log.ErrAttr(err))
	}
}
