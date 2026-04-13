package redis

import (
	"context"
	"encoding/json"
	"log"
	"sync"

	"chess-backend/internal/domain/events"

	"github.com/redis/go-redis/v9"
)

const redisChannelPattern = "*"

type EventBus struct {
	client   *redis.Client
	handlers map[string][]func(ctx context.Context, event events.DomainEvent) error
	mu       sync.RWMutex
}

func NewEventBus(client *redis.Client) *EventBus {
	return &EventBus{
		client:   client,
		handlers: make(map[string][]func(ctx context.Context, event events.DomainEvent) error),
	}
}

func (b *EventBus) Publish(ctx context.Context, event events.DomainEvent) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return b.client.Publish(ctx, event.EventName(), data).Err()
}

func (b *EventBus) Subscribe(eventName string, handler func(ctx context.Context, event events.DomainEvent) error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.handlers[eventName] = append(b.handlers[eventName], handler)
}

func (b *EventBus) StartListening(ctx context.Context) {
	pubsub := b.client.PSubscribe(ctx, redisChannelPattern)
	defer pubsub.Close()

	ch := pubsub.Channel()
	for msg := range ch {
		event, err := events.UnmarshalEvent(msg.Channel, []byte(msg.Payload))
		if err != nil {
			log.Printf("eventbus: failed to unmarshal event on channel %s: %v", msg.Channel, err)
			continue
		}
		b.mu.RLock()
		handlers := b.handlers[event.EventName()]
		b.mu.RUnlock()
		for _, h := range handlers {
			if err := h(ctx, event); err != nil {
				log.Printf("eventbus: handler error for event %s: %v", event.EventName(), err)
			}
		}
	}
}
