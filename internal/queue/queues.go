package queue

import "log/slog"

type Queues struct {
	HelloWorldQueue HelloWorldQueue
}

func NewQueues(pool *ChannelPool) (*Queues, error) {
	helloWorldQueue, err := NewHelloWorldQueue(pool)
	if err != nil {
		slog.Error("unable to create a new queue", "error", err)
		return nil, err
	}

	qs := Queues{
		HelloWorldQueue: *helloWorldQueue,
	}

	return &qs, nil
}
