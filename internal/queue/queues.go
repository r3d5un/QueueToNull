package queue

import "log/slog"

type Queues struct {
	HelloWorldQueue  HelloWorldQueue
	ExampleWorkQueue ExampleWorkQueue
}

func NewQueues(pool *ChannelPool) (*Queues, error) {
	helloWorldQueue, err := NewHelloWorldQueue(pool)
	if err != nil {
		slog.Error("unable to create a new queue", "error", err)
		return nil, err
	}

	exampleWorkQueue, err := NewExampleWorkQueue(pool)
	if err != nil {
		slog.Error("unable to create new queue", "error", err)
		return nil, err
	}

	qs := Queues{
		HelloWorldQueue:  *helloWorldQueue,
		ExampleWorkQueue: *exampleWorkQueue,
	}

	return &qs, nil
}
