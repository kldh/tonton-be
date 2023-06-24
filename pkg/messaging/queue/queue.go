package queue

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/ldhk/tonton-be/pkg/messaging/internal"
)

var errMultipleListeners = errors.New("queue must have only 1 listener")

type Job interface {
	QueueName() string
}

type HandlerFunc func(ctx context.Context, j Job) error

type Queue struct {
	b *internal.LocalQueue

	mu          sync.Mutex
	hasListener map[string]struct{}
}

func New() *Queue {
	return &Queue{
		b:           internal.NewBus(),
		hasListener: make(map[string]struct{}),
	}
}

func (q *Queue) Push(ctx context.Context, j Job) {
	q.b.Publish(ctx, internal.NewMessage(j.QueueName(), j))
}

// Listen to a queue. Each queue can only have 1 listener
func (q *Queue) Listen(queueName string, h HandlerFunc) {
	q.mu.Lock()
	defer q.mu.Unlock()

	_, ok := q.hasListener[queueName]
	if ok {
		panic(fmt.Errorf("%w: queue_name = %s", errMultipleListeners, queueName))
	}
	q.hasListener[queueName] = struct{}{}

	q.b.Subscribe(queueName, func(ctx context.Context, msg *internal.Message) error {
		j := msg.Data().(Job)
		return h(ctx, j)
	})
}

func (q *Queue) Stop() {
	q.b.Stop()
}
