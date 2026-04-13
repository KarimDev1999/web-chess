package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"chess-backend/internal/domain/chess"

	"github.com/redis/go-redis/v9"
)

type DrawOfferStore struct {
	client *redis.Client
}

const drawOfferKeyPattern = "games:%s:draw_offer"

func NewDrawOfferStore(client *redis.Client) *DrawOfferStore {
	return &DrawOfferStore{client: client}
}

func (s *DrawOfferStore) Set(ctx context.Context, gameID string, offer *chess.DrawOffer, ttl time.Duration) error {
	key := s.keyFor(gameID)
	data, err := json.Marshal(offer)
	if err != nil {
		return fmt.Errorf("marshal draw offer: %w", err)
	}
	if err := s.client.Set(ctx, key, data, ttl).Err(); err != nil {
		return fmt.Errorf("set draw offer: %w", err)
	}
	return nil
}

func (s *DrawOfferStore) Get(ctx context.Context, gameID string) (*chess.DrawOffer, error) {
	key := s.keyFor(gameID)
	data, err := s.client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get draw offer: %w", err)
	}
	var offer chess.DrawOffer
	if err := json.Unmarshal([]byte(data), &offer); err != nil {
		return nil, fmt.Errorf("unmarshal draw offer: %w", err)
	}
	return &offer, nil
}

func (s *DrawOfferStore) Delete(ctx context.Context, gameID string) error {
	key := s.keyFor(gameID)
	if err := s.client.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("delete draw offer: %w", err)
	}
	return nil
}

func (s *DrawOfferStore) keyFor(gameID string) string {
	return fmt.Sprintf(drawOfferKeyPattern, gameID)
}
