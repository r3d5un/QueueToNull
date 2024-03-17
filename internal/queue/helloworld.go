package queue

import (
	"context"
	"log/slog"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type HelloWorldQueue struct {
	Queue *amqp.Queue
	Pool  *ChannelPool
}

func NewHelloWorldQueue(pool *ChannelPool) (*HelloWorldQueue, error) {

	ch, err := pool.GetChannel()
	if err != nil {
		slog.Error("unable to get channel", "error", err)
		return nil, err
	}
	defer pool.ReturnChannel(ch)

	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		slog.Error("unable to declare RabbitMQ queue", "error", err)
		return nil, err
	}

	return &HelloWorldQueue{Queue: &q, Pool: pool}, nil
}

func (q HelloWorldQueue) PublishHelloWorld() error {
	ch, err := q.Pool.GetChannel()
	if err != nil {
		slog.Error("unable to get channel", "error", err)
		return err
	}
	defer q.Pool.ReturnChannel(ch)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = ch.PublishWithContext(
		ctx,
		"",           // exchange
		q.Queue.Name, // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte("Hello, World!"),
		})
	if err != nil {
		slog.Error("unable to publish message", "error", err)
		return err
	}

	return nil
}
