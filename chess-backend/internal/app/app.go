package app

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"

	"chess-backend/internal/app/auth"
	"chess-backend/internal/app/game"
	"chess-backend/internal/constants/appconst"
	"chess-backend/internal/domain/events"
	"chess-backend/internal/infrastructure/config"
	"chess-backend/internal/infrastructure/persistence/postgres/db"
	pgrepo "chess-backend/internal/infrastructure/persistence/postgres/repos"
	redispkg "chess-backend/internal/infrastructure/persistence/redis"
	"chess-backend/internal/transport/api"
	"chess-backend/internal/transport/api/handlers"
	"chess-backend/internal/transport/ws"
	"chess-backend/pkg/jwt"
)

type App struct {
	config     *config.Config
	httpServer *http.Server
	hub        *ws.Hub
	eventBus   *redispkg.EventBus
}

func NewApp(cfg *config.Config) (*App, error) {

	pgPool, err := db.NewPostgresPool(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}

	redisClient := redispkg.NewRedisClient(cfg)
	eventBus := redispkg.NewEventBus(redisClient)

	userRepo := pgrepo.NewUserRepository(pgPool)
	gameRepo := pgrepo.NewGameRepository(pgPool)

	drawOfferStore := redispkg.NewDrawOfferStore(redisClient)

	tokenSvc := jwt.NewTokenService(cfg.JWTSecret)
	authSvc := auth.NewAuthService(userRepo, tokenSvc)
	gameSvc := game.NewGameService(gameRepo, eventBus, drawOfferStore)

	hub := ws.NewHub()

	authHandler := handlers.NewAuthHandler(authSvc)
	gameHandler := handlers.NewGameHandler(gameSvc, userRepo.FindByID)
	wsHandler := handlers.NewWebSocketHandler(hub, gameHandler)

	router := transport.NewRouter(authHandler, gameHandler, wsHandler, tokenSvc)

	eventBus.Subscribe(events.EventGameCreated, func(ctx context.Context, event events.DomainEvent) error {
		e, ok := event.(events.GameCreated)
		if !ok {
			return fmt.Errorf("unexpected event type for %s", events.EventGameCreated)
		}
		data, err := json.Marshal(event)
		if err != nil {
			return err
		}
		hub.NotifyPlayers([]string{e.PlayerID}, data)
		return nil
	})
	eventBus.Subscribe(events.EventGameJoined, func(ctx context.Context, event events.DomainEvent) error {
		e, ok := event.(events.GameJoined)
		if !ok {
			return fmt.Errorf("unexpected event type for %s", events.EventGameJoined)
		}
		data, err := json.Marshal(event)
		if err != nil {
			return err
		}
		hub.NotifyPlayers([]string{e.WhitePlayerID, e.BlackPlayerID}, data)
		return nil
	})
	eventBus.Subscribe(events.EventMoveMade, func(ctx context.Context, event events.DomainEvent) error {
		e, ok := event.(events.MoveMade)
		if !ok {
			return fmt.Errorf("unexpected event type for %s", events.EventMoveMade)
		}
		data, err := json.Marshal(event)
		if err != nil {
			return err
		}
		hub.NotifyPlayers([]string{e.WhitePlayerID, e.BlackPlayerID}, data)
		return nil
	})
	eventBus.Subscribe(events.EventDrawOffered, func(ctx context.Context, event events.DomainEvent) error {
		e, ok := event.(events.DrawOffered)
		if !ok {
			return fmt.Errorf("unexpected event type for %s", events.EventDrawOffered)
		}
		data, err := json.Marshal(event)
		if err != nil {
			return err
		}
		hub.NotifyPlayers([]string{e.WhitePlayerID, e.BlackPlayerID}, data)
		return nil
	})
	eventBus.Subscribe(events.EventDrawDeclined, func(ctx context.Context, event events.DomainEvent) error {
		e, ok := event.(events.DrawDeclined)
		if !ok {
			return fmt.Errorf("unexpected event type for %s", events.EventDrawDeclined)
		}
		data, err := json.Marshal(event)
		if err != nil {
			return err
		}
		hub.NotifyPlayers([]string{e.WhitePlayerID, e.BlackPlayerID}, data)
		return nil
	})
	eventBus.Subscribe(events.EventGameResigned, func(ctx context.Context, event events.DomainEvent) error {
		e, ok := event.(events.GameResigned)
		if !ok {
			return fmt.Errorf("unexpected event type for %s", events.EventGameResigned)
		}
		data, err := json.Marshal(event)
		if err != nil {
			return err
		}
		hub.NotifyPlayers([]string{e.WhitePlayerID, e.BlackPlayerID}, data)
		return nil
	})
	eventBus.Subscribe(events.EventDrawAccepted, func(ctx context.Context, event events.DomainEvent) error {
		e, ok := event.(events.DrawAccepted)
		if !ok {
			return fmt.Errorf("unexpected event type for %s", events.EventDrawAccepted)
		}
		data, err := json.Marshal(event)
		if err != nil {
			return err
		}
		hub.NotifyPlayers([]string{e.WhitePlayerID, e.BlackPlayerID}, data)
		return nil
	})
	eventBus.Subscribe(events.EventGameTimedOut, func(ctx context.Context, event events.DomainEvent) error {
		e, ok := event.(events.GameTimedOut)
		if !ok {
			return fmt.Errorf("unexpected event type for %s", events.EventGameTimedOut)
		}
		data, err := json.Marshal(event)
		if err != nil {
			return err
		}
		hub.NotifyPlayers([]string{e.WhitePlayerID, e.BlackPlayerID}, data)
		return nil
	})

	httpServer := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	return &App{
		config:     cfg,
		httpServer: httpServer,
		hub:        hub,
		eventBus:   eventBus,
	}, nil
}

func (a *App) Run() error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go a.hub.Run()
	go a.eventBus.StartListening(ctx)

	go func() {
		log.Printf("HTTP server listening on %s", a.httpServer.Addr)
		if err := a.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("HTTP server error: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("Shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), appconst.HTTPShutdownTimeout)
	defer cancel()

	if err := a.httpServer.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("HTTP server shutdown error: %w", err)
	}

	log.Println("Server gracefully stopped")
	return nil
}
