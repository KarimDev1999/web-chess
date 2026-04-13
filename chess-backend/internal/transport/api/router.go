package transport

import (
	"net/http"

	"chess-backend/internal/constants/appconst"
	"chess-backend/internal/transport/api/handlers"
	"chess-backend/internal/transport/api/middleware"
	"chess-backend/pkg/jwt"

	"github.com/go-chi/chi/v5"
	chimdlw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func NewRouter(
	authHandler *handlers.AuthHandler,
	gameHandler *handlers.GameHandler,
	wsHandler *handlers.WebSocketHandler,
	tokenSvc *jwt.TokenService,
) http.Handler {
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Use(middleware.WSAuthMiddleware(tokenSvc))
		r.Get(appconst.RouteWS, wsHandler.ServeHTTP)
	})

	r.Group(func(r chi.Router) {
		r.Use(chimdlw.Logger)
		r.Use(chimdlw.Recoverer)
		r.Use(cors.Handler(cors.Options{
			AllowedOrigins:   []string{"http://localhost:3000", "http://localhost:5173"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: true,
			MaxAge:           300,
		}))

		r.Post(appconst.RouteRegister, authHandler.Register)
		r.Post(appconst.RouteLogin, authHandler.Login)
		r.Get(appconst.RouteTimeControls, gameHandler.GetTimeControlPresets)

		r.Group(func(r chi.Router) {
			r.Use(middleware.AuthMiddleware(tokenSvc))
			r.Get(appconst.RouteGames, gameHandler.ListMyGames)
			r.Post(appconst.RouteGames, gameHandler.CreateGame)
			r.Get(appconst.RouteGamesWaiting, gameHandler.ListWaitingGames)
			r.Get(appconst.RouteGamesByID, gameHandler.GetGame)
			r.Get(appconst.RouteGamesMoves, gameHandler.GetMoveHistory)
			r.Post(appconst.RouteGamesJoin, gameHandler.JoinGame)
			r.Post(appconst.RouteGamesMove, gameHandler.MakeMove)
			r.Post(appconst.RouteGamesResign, gameHandler.ResignGame)
			r.Post(appconst.RouteGamesOfferDraw, gameHandler.OfferDraw)
			r.Post(appconst.RouteGamesAcceptDraw, gameHandler.AcceptDraw)
			r.Post(appconst.RouteGamesDeclineDraw, gameHandler.DeclineDraw)
		})

		r.Get("/swagger/doc.json", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, appconst.SwaggerDocPath)
		})
		r.Get("/swagger/*", httpSwagger.Handler(
			httpSwagger.URL("/swagger/doc.json"),
		))
	})

	return r
}
