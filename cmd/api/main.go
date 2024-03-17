package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
)

type application struct {
	logger *slog.Logger
}

func main() {
	handler := slog.NewJSONHandler(os.Stdout, nil)
	logger := slog.New(handler)
	instanceLogger := logger.With(
		slog.Group(
			"application_instance",
			slog.String("instance_id", uuid.New().String()),
		),
	)
	slog.SetDefault(instanceLogger)

	slog.Info("starting application")

	app := application{
		logger: logger,
	}

	slog.Info("creating server")
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", 4000),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	logger.Info("starting server", "addr", srv.Addr)
	err := srv.ListenAndServe()
	logger.Error(err.Error())

	os.Exit(1)
}
