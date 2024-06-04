package config

import (
	"os"

	"github.com/joho/godotenv"

	"shortener/internal/log"
)

func LoadEnv() {
	if os.Getenv("APP_MODE") == testMode || os.Getenv("APP_MODE") == prodMode {
		return
	}

	if err := godotenv.Load(".env"); err != nil {
		workDir, errGetWD := os.Getwd()

		if errGetWD != nil {
			log.Error("Error getting work dir",
				log.ErrAttr(errGetWD))
		}

		log.Error("Error loading .env file",
			log.StringAttr("work dir", workDir),
			log.ErrAttr(err),
		)

		setEnvDefault()
	}
}

func setEnvDefault() {
	var cfg Config
	cfg.App.Mode = testMode

	err := os.Setenv("APP_MODE", cfg.App.Mode)
	if err != nil {
		log.Error("Error setting APP_MODE",
			log.ErrAttr(err))
	}

	log.Info("Environment variables set default",
		log.StringAttr("app mode", cfg.App.Mode))
}
