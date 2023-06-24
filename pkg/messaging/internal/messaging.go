package internal

import (
	"context"
	"fmt"
	"runtime/debug"
	"sync"
	"time"

	"github.com/ldhk/tonton-be/pkg/telemetry/logging"
	"github.com/panjf2000/ants/v2"
)

const poolSize = 20000

type Message struct {
	topic string
	data  interface{}
}

func NewMessage(topic string, data interface{}) *Message {
	return &Message{topic: topic, data: data}
}

func (m *Message) Topic() string {
	if m == nil {
		return ""
	}

	return m.topic
}

func (m *Message) Data() interface{} {
	if m == nil {
		return nil
	}

	return m.data
}

type HandlerFunc func(ctx context.Context, msg *Message) error

type LocalQueue struct {
	mu          sync.RWMutex
	subscribers map[string][]HandlerFunc
	pool        *ants.Pool
	wg          *sync.WaitGroup
	errHandler  func(ctx context.Context, err error)
}

func NewBus() *LocalQueue {
	p, err := ants.NewPool(poolSize, ants.Option(func(opts *ants.Options) {
		opts.ExpiryDuration = time.Second
		opts.Nonblocking = false
		opts.PanicHandler = nil
	}))
	if err != nil {
		panic(fmt.Errorf("bus: init pool failed: %v", err))
	}
	return &LocalQueue{
		subscribers: make(map[string][]HandlerFunc),
		pool:        p,
		wg:          new(sync.WaitGroup),
		errHandler: func(ctx context.Context, err error) {
			logging.FromContext(ctx).Errorf("messaging: subscriber return error: %v", err)
		},
	}
}

// Publish message to the bus, this message will be delivered to all subscribers of the message's topic.
func (b *LocalQueue) Publish(ctx context.Context, msg *Message) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	handlers := b.subscribers[msg.topic]
	b.wg.Add(len(handlers))
	for _, h := range handlers {
		func(ctx context.Context, msg *Message, h HandlerFunc) {
			ctx, _ = logging.WithField(ctx, "messaging.topic", msg.topic)

			errSubmit := b.pool.Submit(func() {
				defer b.wg.Done()

				var err error
				defer func() {
					if r := recover(); r != nil {
						err = fmt.Errorf("recover from panic: %v, stack trace: %s", r, debug.Stack())
					}
					if err != nil && b.errHandler != nil {
						b.errHandler(ctx, err)
					}
				}()

				err = h(ctx, msg)
			})
			if errSubmit != nil {
				logging.FromContext(ctx).Errorf("submit message to pool failed: %v", errSubmit)
				b.wg.Done()
			}
		}(ctx, msg, h)
	}
}

// Subscribe to a topic, when a message is published to that topic, it will be pass to the given HandlerFunc
func (b *LocalQueue) Subscribe(topic string, h HandlerFunc) {
	b.mu.Lock()
	b.subscribers[topic] = append(b.subscribers[topic], h)
	b.mu.Unlock()
}

// Stop wait for all running handlers to finish their jobs before return,
// Stop should be called after stop all other incoming channels like gRPC servers or Kafka consumers
func (b *LocalQueue) Stop() {
	b.wg.Wait()
	b.pool.Release()
}

// SetErrHandler a custom error handler when subscriber return error or panic.
// By default, error will be log.
func (b *LocalQueue) SetErrHandler(h func(ctx context.Context, err error)) {
	b.errHandler = h
}
