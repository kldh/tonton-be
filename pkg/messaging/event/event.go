package event

import "context"

type Event interface {
	EventName() string
}

type HandlerFunc func(ctx context.Context, e Event) error

type Bus struct {
	b *internal.Bus
}

func NewBus() *Bus {
	return &Bus{b: internal.NewBus()}
}

func (b *Bus) Publish(ctx context.Context, e Event) {
	b.b.Publish(ctx, internal.NewMessage(e.EventName(), e))
}

func (b *Bus) Subscribe(eventName string, h HandlerFunc) {
	b.b.Subscribe(eventName, func(ctx context.Context, msg *internal.Message) error {
		return h(ctx, msg.Data().(Event))
	})
}

func (b *Bus) Stop() {
	b.b.Stop()
}
