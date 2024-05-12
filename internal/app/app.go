package app

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"

	"shortener/config"
	log "shortener/pkg/logger"
)

const (
	ReadTimeout  = 60 * time.Second
	WriteTimeout = 60 * time.Second
	IdleTimeout  = 60 * time.Second
)

func Handler() http.Handler {
	router := http.NewServeMux()
	router.HandleFunc(http.MethodPost+" /", AddSite)
	router.HandleFunc(http.MethodGet+" /{id}", GetSite)

	return router
}

func RunServer(ctx context.Context, cfg config.Config) error {
	hostPort := cfg.Server.Host + ":" + strconv.Itoa(cfg.Server.Port)
	server := http.Server{
		Addr:                         hostPort,
		Handler:                      Handler(),
		DisableGeneralOptionsHandler: false,
		TLSConfig:                    nil,
		ReadTimeout:                  ReadTimeout,
		ReadHeaderTimeout:            0,
		WriteTimeout:                 WriteTimeout,
		IdleTimeout:                  IdleTimeout,
		MaxHeaderBytes:               0,
		TLSNextProto:                 nil,
		ConnState:                    nil,
		ErrorLog:                     nil,
		BaseContext:                  nil,
		ConnContext: func(_ context.Context, _ net.Conn) context.Context {
			return ctx
		},
	}

	if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("could not start server: %w", err)
	}

	defer func(server *http.Server) { //nolint:contextcheck // only close
		err := server.Close()
		if err != nil {
			log.Error("could not close server", log.ErrAttr(err))
		}
	}(&server)

	return nil
}
