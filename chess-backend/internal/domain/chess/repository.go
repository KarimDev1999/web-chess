package chess

import "context"

type GameRepository interface {
	Save(ctx context.Context, game *Game) error
	FindByID(ctx context.Context, id GameID) (*Game, error)
	FindWaitingGames(ctx context.Context) ([]*Game, error)
	FindByPlayerID(ctx context.Context, playerID string) ([]*Game, error)
	Update(ctx context.Context, game *Game) error
}
