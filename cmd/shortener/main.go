package main

import (
	"context"
	"fmt"
	"os"

	"github.com/joho/godotenv"

	_ "modernc.org/sqlite"

	"shortener/config"
	"shortener/internal/app"
	"shortener/internal/service"
	"shortener/internal/store"
	log "shortener/pkg/logger"
)

func main() {
	loggerInit()
	loadEnv()

	appInfo, err := config.NewAppInfo()
	if err != nil {
		log.Fatal("appInfo", log.ErrAttr(err))
	}
	log.Info("appInfo", log.AnyAttr("app", fmt.Sprint(appInfo)))

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal("cfg", log.ErrAttr(err))
	}
	log.Info("config", log.AnyAttr("cfg", fmt.Sprint(cfg)))

	_, err = config.GetDomains() // for validate
	if err != nil {
		log.Fatal("domains", log.ErrAttr(err))
	}

	db, err := store.NewMemoryStore() // add defer db.CloseDBStore() - only for sqlite3
	if err != nil {
		log.Fatal("store", log.ErrAttr(err))
	}

	shortener := service.NewSiteService(db)

	ctx := context.WithValue(context.Background(), app.KeyServiceCtx{}, shortener)

	if err = app.RunServer(ctx, cfg); err != nil {
		log.Fatal("run", log.ErrAttr(err))
	}
}

func loggerInit() {
	log.NewLogger(
		log.WithLevel("DEBUG"),
		log.WithAddSource(false),
		log.WithIsJSON(true),
		log.WithMiddleware(false),
		log.WithSetDefault(true))
}

func loadEnv() {
	if os.Getenv("APP_MODE") == "test" || os.Getenv("APP_MODE") == "production" {
		return
	}

	if err := godotenv.Load(".env"); err != nil {
		workDir, errGetWD := os.Getwd()
		if errGetWD != nil {
			log.Error("Error getting work dir", log.ErrAttr(errGetWD))
		}
		log.Error("Error loading .env file", log.ErrAttr(err), log.StringAttr("work dir", workDir))
	}
}
