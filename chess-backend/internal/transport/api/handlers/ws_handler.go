package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"chess-backend/internal/app/game"
	"chess-backend/internal/constants/appconst"
	"chess-backend/internal/transport/api/middleware"
	"chess-backend/internal/transport/dto"
	"chess-backend/internal/transport/ws"
	"chess-backend/internal/transport/wsmsg"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type WebSocketHandler struct {
	hub     *ws.Hub
	handler *GameHandler
}

func NewWebSocketHandler(hub *ws.Hub, gameHandler *GameHandler) *WebSocketHandler {
	return &WebSocketHandler{hub: hub, handler: gameHandler}
}

func (h *WebSocketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, appconst.MsgUnauthorized, http.StatusUnauthorized)
		return
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, appconst.MsgWebSocketUpgradeFailed, http.StatusInternalServerError)
		return
	}
	client := &ws.Client{
		ID:     uuid.New().String(),
		UserID: userID,
		Conn:   conn,
		Send:   make(chan []byte, appconst.WSClientBufferSize),
		Hub:    h.hub,
	}
	h.hub.Register(client)

	go client.WritePump()
	go client.ReadPump(h)
}

func (h *WebSocketHandler) HandleClientMessage(client *ws.Client, raw []byte) []byte {
	var msg wsmsg.ClientMessage
	if err := json.Unmarshal(raw, &msg); err != nil {
		return nil
	}

	switch msg.Type {
	case wsmsg.ClientTypePing:
		resp, _ := json.Marshal(map[string]string{wsmsg.KeyType: wsmsg.ClientTypePong})
		return resp

	case wsmsg.ClientTypeJoinGame:
		gameID, ok := extractString(msg.Data, wsmsg.KeyGameID)
		if !ok {
			return nil
		}
		h.hub.JoinGame(client.UserID, gameID)
		return nil

	case wsmsg.ClientTypeLeaveGame:
		gameID, ok := extractString(msg.Data, wsmsg.KeyGameID)
		if !ok {
			return nil
		}
		h.hub.LeaveGame(client.UserID, gameID)
		return nil

	case wsmsg.ClientTypeResign:
		gameID, ok := extractString(msg.Data, wsmsg.KeyGameID)
		if !ok {
			return nil
		}
		return h.handleResign(client.UserID, gameID)

	case wsmsg.ClientTypeOfferDraw:
		gameID, ok := extractString(msg.Data, wsmsg.KeyGameID)
		if !ok {
			return nil
		}
		return h.handleOfferDraw(client.UserID, gameID)

	case wsmsg.ClientTypeAcceptDraw:
		gameID, ok := extractString(msg.Data, wsmsg.KeyGameID)
		if !ok {
			return nil
		}
		return h.handleAcceptDraw(client.UserID, gameID)

	case wsmsg.ClientTypeDeclineDraw:
		gameID, ok := extractString(msg.Data, wsmsg.KeyGameID)
		if !ok {
			return nil
		}
		return h.handleDeclineDraw(client.UserID, gameID)

	default:
		return nil
	}
}

func (h *WebSocketHandler) handleResign(userID, gameID string) []byte {
	if h.handler == nil {
		return nil
	}
	cmd := game.ResignGameCommand{GameID: gameID, PlayerID: userID}
	g, err := h.handler.gameService.ResignGame(context.Background(), cmd)
	if err != nil {
		return nil
	}

	drawOffer := h.handler.gameService.GetPendingDrawOffer(context.Background(), gameID)
	resp := wsmsg.NewGameMessage(wsmsg.TypeGameResigned, gameID,
		dto.ToGameResponse(context.Background(), h.handler.userLookup, g, drawOffer))
	b, _ := json.Marshal(resp)
	return b
}

func (h *WebSocketHandler) handleOfferDraw(userID, gameID string) []byte {
	if h.handler == nil {
		return nil
	}
	cmd := game.OfferDrawCommand{GameID: gameID, PlayerID: userID}
	g, err := h.handler.gameService.OfferDraw(context.Background(), cmd)
	if err != nil {
		return nil
	}

	drawOffer := h.handler.gameService.GetPendingDrawOffer(context.Background(), gameID)
	resp := wsmsg.NewGameMessage(wsmsg.TypeDrawOffered, gameID,
		dto.ToGameResponse(context.Background(), h.handler.userLookup, g, drawOffer))
	b, _ := json.Marshal(resp)
	return b
}

func (h *WebSocketHandler) handleAcceptDraw(userID, gameID string) []byte {
	if h.handler == nil {
		return nil
	}
	cmd := game.AcceptDrawCommand{GameID: gameID, PlayerID: userID}
	g, err := h.handler.gameService.AcceptDraw(context.Background(), cmd)
	if err != nil {
		return nil
	}
	drawOffer := h.handler.gameService.GetPendingDrawOffer(context.Background(), gameID)
	resp := wsmsg.NewGameMessage(wsmsg.TypeDrawAccepted, gameID,
		dto.ToGameResponse(context.Background(), h.handler.userLookup, g, drawOffer))
	b, _ := json.Marshal(resp)
	return b
}

func (h *WebSocketHandler) handleDeclineDraw(userID, gameID string) []byte {
	if h.handler == nil {
		return nil
	}
	cmd := game.DeclineDrawCommand{GameID: gameID, PlayerID: userID}
	g, err := h.handler.gameService.DeclineDraw(context.Background(), cmd)
	if err != nil {
		return nil
	}
	drawOffer := h.handler.gameService.GetPendingDrawOffer(context.Background(), gameID)
	resp := wsmsg.NewGameMessage(wsmsg.TypeDrawDeclined, gameID,
		dto.ToGameResponse(context.Background(), h.handler.userLookup, g, drawOffer))
	b, _ := json.Marshal(resp)
	return b
}

func extractString(data interface{}, key string) (string, bool) {
	m, ok := data.(map[string]interface{})
	if !ok {
		return "", false
	}
	v, ok := m[key]
	if !ok {
		return "", false
	}
	s, ok := v.(string)
	return s, ok
}
