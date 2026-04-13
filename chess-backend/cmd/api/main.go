package main

import (
	"log"

	"chess-backend/internal/app"
	"chess-backend/internal/infrastructure/config"
)

// @title Chess Backend API
// @version 1.0
// @description Real-time chess game backend with WebSocket.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.example.com/support
// @contact.email support@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api
// @schemes http

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	cfg := config.Load()
	application, err := app.NewApp(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize app: %v", err)
	}
	if err := application.Run(); err != nil {
		log.Fatalf("App run error: %v", err)
	}
}
