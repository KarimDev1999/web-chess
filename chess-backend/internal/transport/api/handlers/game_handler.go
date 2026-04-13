package handlers

import (
	"encoding/json"
	"net/http"

	"chess-backend/internal/app/game"
	"chess-backend/internal/constants/appconst"
	"chess-backend/internal/domain/chess"
	"chess-backend/internal/transport/api/middleware"
	"chess-backend/internal/transport/dto"

	"github.com/go-chi/chi/v5"
)

type GameHandler struct {
	gameService *game.GameService
	userLookup  dto.UserLookup
}

func NewGameHandler(gameService *game.GameService, userLookup dto.UserLookup) *GameHandler {
	return &GameHandler{gameService: gameService, userLookup: userLookup}
}

// GetTimeControlPresets godoc
// @Summary Get available time control presets
// @Description Return the list of standard time control presets
// @Tags game
// @Produce json
// @Success 200 {object} map[string][]dto.TimeControlPresetResponse
// @Router /time-controls [get]
func (h *GameHandler) GetTimeControlPresets(w http.ResponseWriter, r *http.Request) {
	presets := chess.StandardPresets()
	result := make(map[string][]dto.TimeControlPresetResponse, len(presets))
	for category, items := range presets {
		for _, p := range items {
			result[category] = append(result[category], dto.TimeControlPresetResponse{
				Label:     p.Label,
				Base:      p.TC.Base,
				Increment: p.TC.Increment,
			})
		}
	}
	w.Header().Set(appconst.HeaderContentType, appconst.MimeTypeJSON)
	json.NewEncoder(w).Encode(result)
}

// CreateGame godoc
// @Summary Create a new game
// @Description Create a chess game as white player
// @Tags game
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 201 {object} dto.GameResponse
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal server error"
// @Router /games [post]
func (h *GameHandler) CreateGame(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, appconst.MsgUnauthorized, http.StatusUnauthorized)
		return
	}

	var req dto.CreateGameRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {

		req = dto.CreateGameRequest{ColorPref: "white"}
	}

	colorPref := chess.ColorPreference(req.ColorPref)
	if colorPref != chess.PreferenceWhite && colorPref != chess.PreferenceBlack && colorPref != chess.PreferenceRandom {
		colorPref = chess.PreferenceWhite
	}

	cmd := game.CreateGameCommand{
		PlayerID: userID,
		TimeControl: chess.TimeControl{
			Base:      req.TimeBase,
			Increment: req.TimeIncrement,
		},
		ColorPref: colorPref,
	}
	g, err := h.gameService.CreateGame(r.Context(), cmd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set(appconst.HeaderContentType, appconst.MimeTypeJSON)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(dto.ToGameResponse(r.Context(), h.userLookup, g, h.gameService.GetPendingDrawOffer(r.Context(), string(g.ID))))
}

// JoinGame godoc
// @Summary Join a game as black player
// @Description Join a waiting game
// @Tags game
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Game ID"
// @Success 200 {object} dto.GameResponse
// @Failure 400 {string} string "Bad request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Game not found"
// @Router /games/{id}/join [post]
func (h *GameHandler) JoinGame(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, appconst.MsgUnauthorized, http.StatusUnauthorized)
		return
	}

	gameID := chi.URLParam(r, "id")
	cmd := game.JoinGameCommand{
		GameID:   gameID,
		PlayerID: userID,
	}
	g, err := h.gameService.JoinGame(r.Context(), cmd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set(appconst.HeaderContentType, appconst.MimeTypeJSON)
	json.NewEncoder(w).Encode(dto.ToGameResponse(r.Context(), h.userLookup, g, h.gameService.GetPendingDrawOffer(r.Context(), string(g.ID))))
}

// ListWaitingGames godoc
// @Summary List waiting games
// @Description Get all games waiting for a player
// @Tags game
// @Produce json
// @Security BearerAuth
// @Success 200 {array} dto.GameResponse
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal server error"
// @Router /games/waiting [get]
func (h *GameHandler) ListWaitingGames(w http.ResponseWriter, r *http.Request) {
	games, err := h.gameService.GetWaitingGames(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set(appconst.HeaderContentType, appconst.MimeTypeJSON)
	json.NewEncoder(w).Encode(dto.ToGameResponses(r.Context(), h.userLookup, games))
}

// MakeMove godoc
// @Summary Make a move in a game
// @Description Submit a move in algebraic notation
// @Tags game
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Game ID"
// @Param request body dto.MoveRequest true "Move details"
// @Success 200 {object} dto.GameResponse
// @Failure 400 {string} string "Invalid move"
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Game not found"
// @Router /games/{id}/move [post]
func (h *GameHandler) MakeMove(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, appconst.MsgUnauthorized, http.StatusUnauthorized)
		return
	}

	gameID := chi.URLParam(r, "id")
	var req dto.MoveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, appconst.MsgInvalidRequest, http.StatusBadRequest)
		return
	}

	cmd := game.MakeMoveCommand{
		GameID:   gameID,
		PlayerID: userID,
		From:     req.From,
		To:       req.To,
	}
	g, err := h.gameService.MakeMove(r.Context(), cmd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set(appconst.HeaderContentType, appconst.MimeTypeJSON)
	json.NewEncoder(w).Encode(dto.ToGameResponse(r.Context(), h.userLookup, g, h.gameService.GetPendingDrawOffer(r.Context(), string(g.ID))))
}

// GetGame godoc
// @Summary Get a game by ID
// @Description Retrieve game details including board state and moves
// @Tags game
// @Produce json
// @Security BearerAuth
// @Param id path string true "Game ID"
// @Success 200 {object} dto.GameResponse
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Game not found"
// @Router /games/{id} [get]
func (h *GameHandler) GetGame(w http.ResponseWriter, r *http.Request) {
	if _, ok := middleware.GetUserID(r.Context()); !ok {
		http.Error(w, appconst.MsgUnauthorized, http.StatusUnauthorized)
		return
	}

	gameID := chi.URLParam(r, "id")
	g, err := h.gameService.GetGameByID(r.Context(), gameID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set(appconst.HeaderContentType, appconst.MimeTypeJSON)
	json.NewEncoder(w).Encode(dto.ToGameResponse(r.Context(), h.userLookup, g, h.gameService.GetPendingDrawOffer(r.Context(), string(g.ID))))
}

// ListMyGames godoc
// @Summary List current user's games
// @Description Get all games the authenticated user participates in
// @Tags game
// @Produce json
// @Security BearerAuth
// @Success 200 {array} dto.GameResponse
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal server error"
// @Router /games [get]
func (h *GameHandler) ListMyGames(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, appconst.MsgUnauthorized, http.StatusUnauthorized)
		return
	}

	games, err := h.gameService.GetMyGames(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set(appconst.HeaderContentType, appconst.MimeTypeJSON)
	json.NewEncoder(w).Encode(dto.ToGameResponses(r.Context(), h.userLookup, games))
}

// GetMoveHistory godoc
// @Summary Get move history for a game
// @Description Retrieve the complete move history of a game
// @Tags game
// @Produce json
// @Security BearerAuth
// @Param id path string true "Game ID"
// @Success 200 {array} dto.MoveResponse
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Game not found"
// @Router /games/{id}/moves [get]
func (h *GameHandler) GetMoveHistory(w http.ResponseWriter, r *http.Request) {
	if _, ok := middleware.GetUserID(r.Context()); !ok {
		http.Error(w, appconst.MsgUnauthorized, http.StatusUnauthorized)
		return
	}

	gameID := chi.URLParam(r, "id")
	moves, err := h.gameService.GetMoveHistory(r.Context(), gameID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set(appconst.HeaderContentType, appconst.MimeTypeJSON)
	resp := make([]dto.MoveResponse, len(moves))
	for i, m := range moves {
		resp[i] = dto.ToMoveResponse(m)
	}
	json.NewEncoder(w).Encode(resp)
}

// ResignGame godoc
// @Summary Resign from a game
// @Description Resign and concede the game to your opponent
// @Tags game
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Game ID"
// @Success 200 {object} dto.GameResponse
// @Failure 400 {string} string "Bad request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Game not found"
// @Router /games/{id}/resign [post]
func (h *GameHandler) ResignGame(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, appconst.MsgUnauthorized, http.StatusUnauthorized)
		return
	}

	gameID := chi.URLParam(r, "id")
	cmd := game.ResignGameCommand{
		GameID:   gameID,
		PlayerID: userID,
	}
	g, err := h.gameService.ResignGame(r.Context(), cmd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set(appconst.HeaderContentType, appconst.MimeTypeJSON)
	json.NewEncoder(w).Encode(dto.ToGameResponse(r.Context(), h.userLookup, g, h.gameService.GetPendingDrawOffer(r.Context(), string(g.ID))))
}

// OfferDraw godoc
// @Summary Offer a draw
// @Description Offer a draw to your opponent
// @Tags game
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Game ID"
// @Success 200 {object} dto.GameResponse
// @Failure 400 {string} string "Bad request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Game not found"
// @Router /games/{id}/draw/offer [post]
func (h *GameHandler) OfferDraw(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, appconst.MsgUnauthorized, http.StatusUnauthorized)
		return
	}

	gameID := chi.URLParam(r, "id")
	cmd := game.OfferDrawCommand{
		GameID:   gameID,
		PlayerID: userID,
	}
	g, err := h.gameService.OfferDraw(r.Context(), cmd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set(appconst.HeaderContentType, appconst.MimeTypeJSON)
	json.NewEncoder(w).Encode(dto.ToGameResponse(r.Context(), h.userLookup, g, h.gameService.GetPendingDrawOffer(r.Context(), string(g.ID))))
}

// AcceptDraw godoc
// @Summary Accept a draw offer
// @Description Accept the opponent's draw offer
// @Tags game
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Game ID"
// @Success 200 {object} dto.GameResponse
// @Failure 400 {string} string "Bad request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Game not found"
// @Router /games/{id}/draw/accept [post]
func (h *GameHandler) AcceptDraw(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, appconst.MsgUnauthorized, http.StatusUnauthorized)
		return
	}

	gameID := chi.URLParam(r, "id")
	cmd := game.AcceptDrawCommand{
		GameID:   gameID,
		PlayerID: userID,
	}
	g, err := h.gameService.AcceptDraw(r.Context(), cmd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set(appconst.HeaderContentType, appconst.MimeTypeJSON)
	json.NewEncoder(w).Encode(dto.ToGameResponse(r.Context(), h.userLookup, g, h.gameService.GetPendingDrawOffer(r.Context(), string(g.ID))))
}

// DeclineDraw godoc
// @Summary Decline a draw offer
// @Description Decline the opponent's draw offer
// @Tags game
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Game ID"
// @Success 200 {object} dto.GameResponse
// @Failure 400 {string} string "Bad request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Game not found"
// @Router /games/{id}/draw/decline [post]
func (h *GameHandler) DeclineDraw(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, appconst.MsgUnauthorized, http.StatusUnauthorized)
		return
	}

	gameID := chi.URLParam(r, "id")
	cmd := game.DeclineDrawCommand{
		GameID:   gameID,
		PlayerID: userID,
	}
	g, err := h.gameService.DeclineDraw(r.Context(), cmd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set(appconst.HeaderContentType, appconst.MimeTypeJSON)
	json.NewEncoder(w).Encode(dto.ToGameResponse(r.Context(), h.userLookup, g, h.gameService.GetPendingDrawOffer(r.Context(), string(g.ID))))
}
