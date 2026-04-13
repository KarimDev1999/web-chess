# Web Chess

A real-time multiplayer chess application with WebSocket support, built with Go and React.

## Features

- Real-time two-player chess with WebSocket communication
- Full chess rules: castling, en-passant, pawn promotion, checkmate, stalemate
- Time controls: Bullet, Blitz, Rapid, Classical
- Draw offers with timeout, resignation
- JWT authentication
- Redis-backed ephemeral state for draw offers
- PostgreSQL persistence
- Current game position storage in FEN format

## Architecture

```
chess-backend/          # Go backend (Clean Architecture)
├── cmd/api/            # Application entry point
├── internal/
│   ├── app/            # Application layer (services, commands)
│   ├── domain/         # Domain layer (chess engine, entities)
│   ├── infrastructure/ # PostgreSQL, Redis
│   └── transport/      # HTTP, WebSocket, DTOs
└── pkg/                # Shared utilities (JWT, math)

chess-frontend/         # React + TypeScript frontend
└── src/
    ├── pages/          # GameView, Lobby, Auth
    ├── components/     # ChessBoard, PlayerProfile
    ├── hooks/          # WebSocket with auto-reconnect
    └── lib/            # Chess logic, types, config
```

## Tech Stack

| Layer | Technology |
|-------|-----------|
| Backend | Go 1.25, chi router, gorilla/websocket |
| Auth | JWT (HS256), bcrypt |
| Database | PostgreSQL (pgx), Redis (pub/sub + ephemeral state) |
| Frontend | React 19, TypeScript 5.8, Vite 6 |
| Docs | Swagger/OpenAPI |

## Development

### Backend

```bash
cd chess-backend
make fmt       # Format Go files
make vet       # Run go vet
make test      # Run tests with race detection
make build     # Build all packages
make lint      # fmt + vet + test
make migrate-up
make migrate-down
make migrate-create
```

### Frontend

```bash
cd chess-frontend
npm run format       # Format with Prettier
npm run format:check # Check formatting (CI)
npm run lint         # Run ESLint
```

## Quick Start

### Backend

```bash
cd chess-backend
docker compose up -d
go run cmd/api/main.go
```

The server starts on `:8080` with Swagger docs at `/swagger/`.

### Frontend

```bash
cd chess-frontend
npm install
npm run dev
```

## API

- `POST /api/register` — Register
- `POST /api/login` — Login
- `GET /api/games` — My games
- `POST /api/games` — Create game
- `GET /api/games/waiting` — Waiting games
- `GET /api/games/:id` — Game details
- `POST /api/games/:id/join` — Join game
- `POST /api/games/:id/move` — Make move
- `POST /api/games/:id/resign` — Resign
- `POST /api/games/:id/draw/offer` — Offer draw
- `POST /api/games/:id/draw/accept` — Accept draw
- `POST /api/games/:id/draw/decline` — Decline draw
- `GET /api/time-controls` — Time control presets
- `GET /swagger/` — OpenAPI docs

## WebSocket

Connect to `/ws?token=<jwt>` for real-time updates. Messages use a `type` field for dispatching:

- `game.created`, `game.joined`, `game.move_made` — Game state changes
- `game.draw_offered`, `game.draw_declined`, `game.draw_accepted` — Draw events
- `game.resigned`, `game.timed_out` — Game end events
- `presence` — Player presence updates

## License

MIT
