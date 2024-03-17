package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/r3d5un/QueueToNull/internal/queue"
	amqp "github.com/rabbitmq/amqp091-go"
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

	slog.Info("connecting to RabbitMQ")
	conn, err := amqp.Dial("amqp://devUser:password@localhost:5672/my_vhost")
	if err != nil {
		logger.Error("unable to dial RabbitMQ", "error", err)
		os.Exit(1)
	}
	defer conn.Close()

	mqPool, err := queue.NewChannelPool(conn, 10)
	if err != nil {
		logger.Error("unable to create RabbitMQ channel pool", "error", err)
		os.Exit(1)
	}

	c, err := mqPool.GetChannel()
	if err != nil {
		logger.Error("unable to get channel from pool", "error", err)
		os.Exit(1)
	}
	defer mqPool.ReturnChannel(c)

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
	err = srv.ListenAndServe()
	logger.Error(err.Error())

	os.Exit(1)
}
