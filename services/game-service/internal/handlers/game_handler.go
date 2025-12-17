package handlers

import (
	"encoding/json"
	"net/http"

	"portfolio-game-service/internal/domain"
	"portfolio-game-service/internal/services"

	"github.com/sirupsen/logrus"
)

type GameHandler struct {
	gameService *services.GameService
	logger      *logrus.Logger
}

type StartGameResponse struct {
	GameID string `json:"game_id"`
	Status string `json:"status"`
}

type MoveRequest struct {
	GameID    string `json:"game_id"`
	Action    string `json:"action"`
	Direction string `json:"direction,omitempty"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func NewGameHandler(gameService *services.GameService, logger *logrus.Logger) *GameHandler {
	return &GameHandler{
		gameService: gameService,
		logger:      logger,
	}
}

func (h *GameHandler) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
}

func (h *GameHandler) StartGame(w http.ResponseWriter, r *http.Request) {
	game, err := h.gameService.StartGame()
	if err != nil {
		h.writeError(w, "Failed to start game", http.StatusInternalServerError)
		return
	}

	response := StartGameResponse{
		GameID: game.ID,
		Status: string(game.Status),
	}

	h.writeJSON(w, response, http.StatusCreated)
}

func (h *GameHandler) MakeMove(w http.ResponseWriter, r *http.Request) {
	var req MoveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	game, err := h.gameService.MakeMove(req.GameID, req.Action, req.Direction)
	if err != nil {
		switch err {
		case domain.ErrGameNotFound:
			h.writeError(w, "Game not found", http.StatusNotFound)
		case domain.ErrGameOver:
			h.writeError(w, "Game is over", http.StatusBadRequest)
		case domain.ErrInvalidMove:
			h.writeError(w, "Invalid move", http.StatusBadRequest)
		default:
			h.writeError(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	h.writeJSON(w, game, http.StatusOK)
}

func (h *GameHandler) GetStatus(w http.ResponseWriter, r *http.Request) {
	gameID := r.URL.Query().Get("game_id")
	if gameID == "" {
		h.writeError(w, "Missing game_id parameter", http.StatusBadRequest)
		return
	}

	game, err := h.gameService.GetGameStatus(gameID)
	if err != nil {
		if err == domain.ErrGameNotFound {
			h.writeError(w, "Game not found", http.StatusNotFound)
		} else {
			h.writeError(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	h.writeJSON(w, game, http.StatusOK)
}

func (h *GameHandler) writeJSON(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.logger.WithError(err).Error("Failed to encode JSON response")
	}
}

func (h *GameHandler) writeError(w http.ResponseWriter, message string, status int) {
	h.logger.WithFields(logrus.Fields{
		"error": message,
		"status": status,
	}).Error("Request error")
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}