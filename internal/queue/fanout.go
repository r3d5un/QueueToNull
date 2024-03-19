package queue

import (
	"context"
	"log/slog"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type FanoutExample struct {
	FirstQueue  *amqp.Queue
	SecondQueue *amqp.Queue
	Pool        *ChannelPool
}

func NewFanoutExample(pool *ChannelPool) (*FanoutExample, error) {
	ch, err := pool.GetChannel()
	if err != nil {
		slog.Error("unable to get channel", "error", err)
		return nil, err
	}
	defer pool.ReturnChannel(ch)

	err = ch.ExchangeDeclare(
		"fanout_example", // name
		"fanout",         // type
		true,             // durable
		false,            // auto-deleted
		false,            // internal
		false,            // no-wait
		nil,              // arguments
	)
	if err != nil {
		slog.Error("unable to create exchange", "error", err)
		return nil, err
	}

	firstQ, err := ch.QueueDeclare(
		"first_queue", // name
		false,         // durable
		false,         // delete when unused
		true,          // exclusive
		false,         // no-wait
		nil,           // arguments
	)
	if err != nil {
		slog.Error("unable to create the first queue", "error", err)
		return nil, err
	}

	secondQ, err := ch.QueueDeclare(
		"second_queue", // name
		false,          // durable
		false,          // delete when unused
		true,           // exclusive
		false,          // no-wait
		nil,            // arguments
	)
	if err != nil {
		slog.Error("unable to create the second queue", "error", err)
		return nil, err
	}

	err = ch.QueueBind(
		firstQ.Name,      // name
		"",               // routing key
		"fanout_example", // exchange
		false,
		nil,
	)
	if err != nil {
		slog.Error("unable to bind the first queue", "error", err)
		return nil, err
	}

	err = ch.QueueBind(
		secondQ.Name,     // name
		"",               // routing key
		"fanout_example", // exchange
		false,
		nil,
	)
	if err != nil {
		slog.Error("unable to bind the second queue", "error", err)
		return nil, err
	}

	return &FanoutExample{FirstQueue: &firstQ, SecondQueue: &secondQ, Pool: pool}, nil
}

// Publishes a plain text string to an exachge that fansout to several queues.
func (q *FanoutExample) PublishMessage(msg string) error {
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
		"fanout_example",
		"",
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
// queue called "first_queue".
func (q FanoutExample) GetFirstQueueMessages(autoAck bool) (<-chan amqp.Delivery, error) {
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
		q.FirstQueue.Name, // queue
		"",                // consumer
		autoAck,           // auto-ack
		false,             // exclusive
		false,             // no-local
		false,             // no/wait
		nil,               // args
	)
	if err != nil {
		slog.Error("unable to consume messages", "error", err)
		return nil, err
	}

	return msgs, nil
}

// GetMessages returns a Go chan that consumes messages from the RabbitMQ
// queue called "first_queue".
func (q FanoutExample) GetSecondQueueMessages(autoAck bool) (<-chan amqp.Delivery, error) {
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
		q.SecondQueue.Name, // queue
		"",                 // consumer
		autoAck,            // auto-ack
		false,              // exclusive
		false,              // no-local
		false,              // no/wait
		nil,                // args
	)
	if err != nil {
		slog.Error("unable to consume messages", "error", err)
		return nil, err
	}

	return msgs, nil
}
