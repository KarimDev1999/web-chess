package game

import (
	"context"
	"errors"
	"fmt"
	"time"

	"chess-backend/internal/app/interfaces"
	"chess-backend/internal/domain/chess"
	"chess-backend/internal/domain/events"
)

const DrawOfferTimeout = 30 * time.Second

type GameService struct {
	gameRepo       chess.GameRepository
	pub            interfaces.EventPublisher
	drawOfferStore interfaces.DrawOfferStore
}

func NewGameService(gameRepo chess.GameRepository, pub interfaces.EventPublisher, drawOfferStore interfaces.DrawOfferStore) *GameService {
	return &GameService{
		gameRepo:       gameRepo,
		pub:            pub,
		drawOfferStore: drawOfferStore,
	}
}

func (s *GameService) CreateGame(ctx context.Context, cmd CreateGameCommand) (*chess.Game, error) {
	game := chess.NewGame(cmd.PlayerID, cmd.TimeControl, cmd.ColorPref)
	if err := s.gameRepo.Save(ctx, game); err != nil {
		return nil, fmt.Errorf("failed to save game: %w", err)
	}

	go func() {
		_ = s.pub.Publish(ctx, events.GameCreated{
			BaseEvent: events.BaseEvent{
				Type:      events.EventGameCreated,
				Timestamp: time.Now(),
			},
			GameID:        string(game.ID),
			WhitePlayerID: cmd.PlayerID,
		})
	}()

	return game, nil
}

func (s *GameService) JoinGame(ctx context.Context, cmd JoinGameCommand) (*chess.Game, error) {
	game, err := s.gameRepo.FindByID(ctx, chess.GameID(cmd.GameID))
	if err != nil {
		return nil, fmt.Errorf("failed to find game: %w", err)
	}
	if game == nil {
		return nil, chess.ErrGameNotFound
	}

	if err := game.Join(cmd.PlayerID); err != nil {
		return nil, err
	}

	if err := s.gameRepo.Update(ctx, game); err != nil {
		if errors.Is(err, chess.ErrConcurrentUpdate) {
			return nil, chess.ErrConcurrentUpdate
		}
		return nil, fmt.Errorf("failed to update game: %w", err)
	}

	go func() {
		white, _ := s.playerIDs(game)
		_ = s.pub.Publish(ctx, events.GameJoined{
			BaseEvent: events.BaseEvent{
				Type:      events.EventGameJoined,
				Timestamp: time.Now(),
			},
			GameID:        cmd.GameID,
			WhitePlayerID: white,
			BlackPlayerID: cmd.PlayerID,
		})
	}()

	return game, nil
}

func (s *GameService) MakeMove(ctx context.Context, cmd MakeMoveCommand) (*chess.Game, error) {
	game, err := s.gameRepo.FindByID(ctx, chess.GameID(cmd.GameID))
	if err != nil {
		return nil, fmt.Errorf("failed to find game: %w", err)
	}
	if game == nil {
		return nil, chess.ErrGameNotFound
	}

	from, err := chess.ParseAlgebraic(cmd.From)
	if err != nil {
		return nil, fmt.Errorf("invalid from position: %w", err)
	}
	to, err := chess.ParseAlgebraic(cmd.To)
	if err != nil {
		return nil, fmt.Errorf("invalid to position: %w", err)
	}

	if err := game.MakeMove(cmd.PlayerID, from, to); err != nil {
		return nil, err
	}

	if err := s.gameRepo.Update(ctx, game); err != nil {
		if errors.Is(err, chess.ErrConcurrentUpdate) {
			return nil, chess.ErrConcurrentUpdate
		}
		return nil, fmt.Errorf("failed to update game: %w", err)
	}

	s.clearDrawOffer(ctx, cmd.GameID)

	var resultStr, endReasonStr *string
	if game.Result != nil {
		s := string(*game.Result)
		resultStr = &s
	}
	if game.EndReason != nil {
		s := string(*game.EndReason)
		endReasonStr = &s
	}

	go func() {
		white, black := s.playerIDs(game)
		_ = s.pub.Publish(ctx, events.MoveMade{
			BaseEvent: events.BaseEvent{
				Type:      events.EventMoveMade,
				Timestamp: time.Now(),
			},
			GameID:        cmd.GameID,
			WhitePlayerID: white,
			BlackPlayerID: black,
			PlayerID:      cmd.PlayerID,
			From:          cmd.From,
			To:            cmd.To,
			Result:        resultStr,
			EndReason:     endReasonStr,
		})
	}()

	return game, nil
}

func (s *GameService) GetWaitingGames(ctx context.Context) ([]*chess.Game, error) {
	games, err := s.gameRepo.FindWaitingGames(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch waiting games: %w", err)
	}
	return games, nil
}

func (s *GameService) GetGameByID(ctx context.Context, id string) (*chess.Game, error) {
	game, err := s.gameRepo.FindByID(ctx, chess.GameID(id))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch game: %w", err)
	}
	if game == nil {
		return nil, chess.ErrGameNotFound
	}

	s.checkAndEndTimeout(ctx, game)

	return game, nil
}

func (s *GameService) checkAndEndTimeout(ctx context.Context, game *chess.Game) {
	if !game.TimeControl.IsTimed() || game.Status != chess.StatusActive {
		return
	}

	var playerID *string
	if game.Turn == chess.White {
		playerID = game.WhitePlayerID
	} else {
		playerID = game.BlackPlayerID
	}
	if playerID == nil {
		return
	}

	if game.RemainingFor(game.Turn) > 0 {
		return
	}

	if err := game.EndOnTimeout(*playerID); err != nil {
		return
	}
	if err := s.gameRepo.Update(ctx, game); err != nil {
		return
	}
	s.clearDrawOffer(ctx, string(game.ID))

	s.publishGameTimedOutEvent(ctx, game, *playerID)
}

func (s *GameService) GetMyGames(ctx context.Context, playerID string) ([]*chess.Game, error) {
	games, err := s.gameRepo.FindByPlayerID(ctx, playerID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch games: %w", err)
	}
	return games, nil
}

func (s *GameService) GetMoveHistory(ctx context.Context, id string) ([]chess.Move, error) {
	game, err := s.gameRepo.FindByID(ctx, chess.GameID(id))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch game: %w", err)
	}
	if game == nil {
		return nil, chess.ErrGameNotFound
	}
	return game.Moves, nil
}

func (s *GameService) ResignGame(ctx context.Context, cmd ResignGameCommand) (*chess.Game, error) {
	game, err := s.gameRepo.FindByID(ctx, chess.GameID(cmd.GameID))
	if err != nil {
		return nil, fmt.Errorf("failed to find game: %w", err)
	}
	if game == nil {
		return nil, chess.ErrGameNotFound
	}

	if err := game.Resign(cmd.PlayerID); err != nil {
		return nil, err
	}

	s.clearDrawOffer(ctx, cmd.GameID)

	if err := s.gameRepo.Update(ctx, game); err != nil {
		if errors.Is(err, chess.ErrConcurrentUpdate) {
			return nil, chess.ErrConcurrentUpdate
		}
		return nil, fmt.Errorf("failed to update game: %w", err)
	}

	s.publishResignEvent(ctx, game, cmd.PlayerID)
	return game, nil
}

func (s *GameService) OfferDraw(ctx context.Context, cmd OfferDrawCommand) (*chess.Game, error) {
	game, err := s.gameRepo.FindByID(ctx, chess.GameID(cmd.GameID))
	if err != nil {
		return nil, fmt.Errorf("failed to find game: %w", err)
	}
	if game == nil {
		return nil, chess.ErrGameNotFound
	}

	existing, err := s.drawOfferStore.Get(ctx, cmd.GameID)
	if err != nil {
		return nil, fmt.Errorf("failed to check draw offer: %w", err)
	}
	if existing != nil && existing.OfferedBy == cmd.PlayerID {
		return nil, chess.ErrDrawAlreadyOffered
	}

	offer := &chess.DrawOffer{
		OfferedBy: cmd.PlayerID,
		OfferedAt: time.Now(),
	}
	if err := s.drawOfferStore.Set(ctx, cmd.GameID, offer, DrawOfferTimeout); err != nil {
		return nil, fmt.Errorf("failed to store draw offer: %w", err)
	}

	s.publishDrawOfferedEvent(ctx, game, cmd.PlayerID)
	return game, nil
}

func (s *GameService) AcceptDraw(ctx context.Context, cmd AcceptDrawCommand) (*chess.Game, error) {
	game, err := s.gameRepo.FindByID(ctx, chess.GameID(cmd.GameID))
	if err != nil {
		return nil, fmt.Errorf("failed to find game: %w", err)
	}
	if game == nil {
		return nil, chess.ErrGameNotFound
	}

	offer, err := s.drawOfferStore.Get(ctx, cmd.GameID)
	if err != nil {
		return nil, fmt.Errorf("failed to get draw offer: %w", err)
	}
	if offer == nil {
		return nil, chess.ErrNoDrawOffer
	}
	if offer.OfferedBy == cmd.PlayerID {
		return nil, chess.ErrCannotAcceptOwnOffer
	}
	if !game.IsPlayerInGame(cmd.PlayerID) {
		return nil, chess.ErrPlayerNotInGame
	}
	if game.Status != chess.StatusActive {
		return nil, chess.ErrGameNotActive
	}

	if err := s.drawOfferStore.Delete(ctx, cmd.GameID); err != nil {
		return nil, fmt.Errorf("failed to remove draw offer: %w", err)
	}

	draw := chess.ResultDraw
	game.Result = &draw
	reason := chess.ReasonDrawAgree
	game.EndReason = &reason
	game.Status = chess.StatusFinished
	game.UpdatedAt = time.Now()

	if err := s.gameRepo.Update(ctx, game); err != nil {
		if errors.Is(err, chess.ErrConcurrentUpdate) {
			return nil, chess.ErrConcurrentUpdate
		}
		return nil, fmt.Errorf("failed to update game: %w", err)
	}

	s.publishDrawAcceptedEvent(ctx, game, cmd.PlayerID)
	return game, nil
}

func (s *GameService) DeclineDraw(ctx context.Context, cmd DeclineDrawCommand) (*chess.Game, error) {
	game, err := s.gameRepo.FindByID(ctx, chess.GameID(cmd.GameID))
	if err != nil {
		return nil, fmt.Errorf("failed to find game: %w", err)
	}
	if game == nil {
		return nil, chess.ErrGameNotFound
	}

	if err := s.drawOfferStore.Delete(ctx, cmd.GameID); err != nil {
		return nil, fmt.Errorf("failed to remove draw offer: %w", err)
	}

	if !game.IsPlayerInGame(cmd.PlayerID) {
		return nil, chess.ErrPlayerNotInGame
	}

	game.UpdatedAt = time.Now()
	if err := s.gameRepo.Update(ctx, game); err != nil {
		if errors.Is(err, chess.ErrConcurrentUpdate) {
			return nil, chess.ErrConcurrentUpdate
		}
		return nil, fmt.Errorf("failed to update game: %w", err)
	}

	s.publishDrawDeclinedEvent(ctx, game, cmd.PlayerID)
	return game, nil
}

func (s *GameService) clearDrawOffer(ctx context.Context, gameID string) {
	_ = s.drawOfferStore.Delete(ctx, gameID)
}

func (s *GameService) GetPendingDrawOffer(ctx context.Context, gameID string) *chess.DrawOffer {
	offer, err := s.drawOfferStore.Get(ctx, gameID)
	if err != nil {
		return nil
	}
	return offer
}

func (s *GameService) playerIDs(game *chess.Game) (white, black string) {
	if game.WhitePlayerID != nil {
		white = *game.WhitePlayerID
	}
	if game.BlackPlayerID != nil {
		black = *game.BlackPlayerID
	}
	return
}

func (s *GameService) publishResignEvent(ctx context.Context, game *chess.Game, playerID string) {
	white, black := s.playerIDs(game)
	var resultStr, endReasonStr *string
	if game.Result != nil {
		r := string(*game.Result)
		resultStr = &r
	}
	if game.EndReason != nil {
		e := string(*game.EndReason)
		endReasonStr = &e
	}
	go func() {
		_ = s.pub.Publish(ctx, events.GameResigned{
			BaseEvent: events.BaseEvent{
				Type:      events.EventGameResigned,
				Timestamp: time.Now(),
			},
			GameID:        string(game.ID),
			WhitePlayerID: white,
			BlackPlayerID: black,
			ResignedBy:    playerID,
			Result:        resultStr,
			EndReason:     endReasonStr,
		})
	}()
}

func (s *GameService) publishDrawOfferedEvent(ctx context.Context, game *chess.Game, playerID string) {
	white, black := s.playerIDs(game)
	go func() {
		_ = s.pub.Publish(ctx, events.DrawOffered{
			BaseEvent: events.BaseEvent{
				Type:      events.EventDrawOffered,
				Timestamp: time.Now(),
			},
			GameID:        string(game.ID),
			WhitePlayerID: white,
			BlackPlayerID: black,
			OfferedBy:     playerID,
		})
	}()
}

func (s *GameService) publishDrawDeclinedEvent(ctx context.Context, game *chess.Game, playerID string) {
	white, black := s.playerIDs(game)
	go func() {
		_ = s.pub.Publish(ctx, events.DrawDeclined{
			BaseEvent: events.BaseEvent{
				Type:      events.EventDrawDeclined,
				Timestamp: time.Now(),
			},
			GameID:        string(game.ID),
			WhitePlayerID: white,
			BlackPlayerID: black,
			DeclinedBy:    playerID,
		})
	}()
}

func (s *GameService) publishDrawOfferExpiredEvent(ctx context.Context, game *chess.Game) {
	white, black := s.playerIDs(game)
	go func() {
		_ = s.pub.Publish(ctx, events.DrawOfferExpired{
			BaseEvent: events.BaseEvent{
				Type:      events.EventDrawOfferExpired,
				Timestamp: time.Now(),
			},
			GameID:        string(game.ID),
			WhitePlayerID: white,
			BlackPlayerID: black,
		})
	}()
}

func (s *GameService) publishDrawAcceptedEvent(ctx context.Context, game *chess.Game, playerID string) {
	white, black := s.playerIDs(game)
	var resultStr, endReasonStr *string
	if game.Result != nil {
		r := string(*game.Result)
		resultStr = &r
	}
	if game.EndReason != nil {
		e := string(*game.EndReason)
		endReasonStr = &e
	}
	go func() {
		_ = s.pub.Publish(ctx, events.DrawAccepted{
			BaseEvent: events.BaseEvent{
				Type:      events.EventDrawAccepted,
				Timestamp: time.Now(),
			},
			GameID:        string(game.ID),
			WhitePlayerID: white,
			BlackPlayerID: black,
			AcceptedBy:    playerID,
			Result:        resultStr,
			EndReason:     endReasonStr,
		})
	}()
}

func (s *GameService) publishGameTimedOutEvent(ctx context.Context, game *chess.Game, timedOutPlayerID string) {
	white, black := s.playerIDs(game)
	var resultStr, endReasonStr *string
	if game.Result != nil {
		r := string(*game.Result)
		resultStr = &r
	}
	if game.EndReason != nil {
		e := string(*game.EndReason)
		endReasonStr = &e
	}
	go func() {
		_ = s.pub.Publish(ctx, events.GameTimedOut{
			BaseEvent: events.BaseEvent{
				Type:      events.EventGameTimedOut,
				Timestamp: time.Now(),
			},
			GameID:        string(game.ID),
			WhitePlayerID: white,
			BlackPlayerID: black,
			PlayerID:      timedOutPlayerID,
			Result:        resultStr,
			EndReason:     endReasonStr,
		})
	}()
}
