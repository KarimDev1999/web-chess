package interfaces

import (
	"context"
	"time"

	"chess-backend/internal/domain/chess"
	"chess-backend/internal/domain/events"
)

type EventPublisher interface {
	Publish(ctx context.Context, event events.DomainEvent) error
}

type DrawOfferStore interface {
	Set(ctx context.Context, gameID string, offer *chess.DrawOffer, ttl time.Duration) error
	Get(ctx context.Context, gameID string) (*chess.DrawOffer, error)
	Delete(ctx context.Context, gameID string) error
}
