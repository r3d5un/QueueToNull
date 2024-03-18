package queue

import (
	"context"
	"log/slog"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type ExampleWorkQueue struct {
	Queue *amqp.Queue
	Pool  *ChannelPool
}

func NewExampleWorkQueue(pool *ChannelPool) (*ExampleWorkQueue, error) {
	ch, err := pool.GetChannel()
	if err != nil {
		slog.Error("unable to get channel", "error", err)
		return nil, err
	}
	defer pool.ReturnChannel(ch)

	q, err := ch.QueueDeclare(
		"task_queue", // name
		true,         // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		slog.Error("unable to declare RabbitMQ queue", "error", err)
		return nil, err
	}

	return &ExampleWorkQueue{Queue: &q, Pool: pool}, nil
}

func (q ExampleWorkQueue) PublishMessage(msg string) error {
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
		"",
		q.Queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(msg),
		})
	if err != nil {
		slog.Error("unable to publish message", "message", msg, "error", err)
		return err
	}

	return nil
}

// GetMessages returns a Go chan that consumes messages from the RabbitMQ
// queue.
func (q ExampleWorkQueue) GetMessages(autoAck bool) (<-chan amqp.Delivery, error) {
	ch, err := q.Pool.GetChannel()
	if err != nil {
		slog.Error("unable to get channel", "error", err)
		return nil, err
	}
	defer q.Pool.ReturnChannel(ch)

	// The following snippet tells RabbitMQ not to give more than one message
	// to a worker at the time, until said worker has processed and acknowledged
	// the previous message.
	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		slog.Error("unable to set fair dispatch", "error", err)
		return nil, err
	}

	msgs, err := ch.Consume(
		q.Queue.Name, // queue
		"",           // consumer
		autoAck,      // auto-ack
		false,        // exclusive
		false,        // no-local
		false,        // no/wait
		nil,          // args
	)
	if err != nil {
		slog.Error("unable to consume messages", "error", err)
		return nil, err
	}

	return msgs, nil
}
