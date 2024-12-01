package cmd

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"effective-mobile-song-library/config"
	httphandl "effective-mobile-song-library/internal/delivery/http"
	"effective-mobile-song-library/pkg/logger"
)

type httpserver struct {
	handler *httphandl.Handler
	config  *config.Config
}

func NewServer(handler *httphandl.Handler, cfg *config.Config) httpserver {
	return httpserver{
		handler: handler,
		config:  cfg,
	}
}

func (s httpserver) Start() error {
	// Declare a HTTP server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", s.config.Port),
		Handler:      s.handler.Routes(),
		ErrorLog:     log.New(logger.New(os.Stdout, logger.LevelInfo), "", 0),
		IdleTimeout:  time.Minute,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	shutdownError := make(chan error)

	// Starting a background goroutine. graceful shutdown
	go func() {
		quit := make(chan os.Signal, 1)

		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		signal := <-quit

		logger.PrintInfo("shutting down server", map[string]any{
			"signal": signal.String(),
		})

		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		shutdownError <- srv.Shutdown(ctx)
	}()

	logger.PrintInfo("starting server", map[string]any{
		"addr": srv.Addr,
	})

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdownError
	if err != nil {
		return err
	}

	logger.PrintInfo("stopped server", map[string]any{
		"addr": srv.Addr,
	})

	return nil
}