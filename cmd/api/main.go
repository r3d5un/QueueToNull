package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/google/uuid"
	"github.com/r3d5un/QueueToNull/internal/queue"
	amqp "github.com/rabbitmq/amqp091-go"
)

type application struct {
	logger *slog.Logger
	queues *queue.Queues
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
	slog.Info("connected to RabbitMQ")

	slog.Info("creating RabbitMQ channel pool")
	mqPool, err := queue.NewChannelPool(conn, 10)
	if err != nil {
		logger.Error("unable to create RabbitMQ channel pool", "error", err)
		os.Exit(1)
	}
	defer mqPool.Shutdown()
	slog.Info("RabbitMQ channel pool created")

	slog.Info("creating queues")
	queues, err := queue.NewQueues(mqPool)
	if err != nil {
		logger.Error("unable to create queues", "error", err)
		os.Exit(1)
	}
	slog.Info("queues created")

	slog.Info("setting up background processes")
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	// Creating two consumers for demonstration purposes
	go queue.ConsumeExampleWorkQueue(queues.ExampleWorkQueue, done)
	go queue.ConsumeExampleWorkQueue(queues.ExampleWorkQueue, done)

	app := application{
		logger: logger,
		queues: queues,
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
