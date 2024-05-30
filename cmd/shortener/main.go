package main

import (
	"shortener/config"
	"shortener/internal/log"
	"shortener/internal/shortener"
)

func main() {
	log.Prepare()

	config.LoadEnv()

	err := config.LogAppInfo()
	if err != nil {
		log.Fatal("appInfo",
			log.ErrAttr(err))
	}

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal("cfg", log.ErrAttr(err))
	}

	log.Info("config",
		log.StringAttr("address", string(cfg.Shortener.Address)),
		log.StringAttr("filepath", cfg.Store.FileStoragePath))

	err = shortener.Run(cfg)
	if err != nil {
		log.Fatal("shortener run error",
			log.ErrAttr(err))
	}
}
