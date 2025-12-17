package services

import (
	"crypto/rand"
	"encoding/hex"
	"sync"

	"portfolio-game-service/internal/domain"
	"portfolio-game-service/pkg/metrics"

	"github.com/sirupsen/logrus"
)

type GameService struct {
	games  map[string]*domain.Game
	mutex  sync.RWMutex
	logger *logrus.Logger
}

func NewGameService(logger *logrus.Logger) *GameService {
	return &GameService{
		games:  make(map[string]*domain.Game),
		logger: logger,
	}
}

func (s *GameService) StartGame() (*domain.Game, error) {
	gameID := s.generateGameID()
	
	game := domain.NewGame(gameID)
	
	s.mutex.Lock()
	s.games[gameID] = game
	s.mutex.Unlock()
	
	metrics.GamesStarted.Inc()
	s.logger.WithFields(logrus.Fields{
		"game_id": gameID,
		"score": game.Score,
	}).Info("New game started")
	
	return game, nil
}

func (s *GameService) MakeMove(gameID, action, direction string) (*domain.Game, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	
	game, exists := s.games[gameID]
	if !exists {
		return nil, domain.ErrGameNotFound
	}
	
	var err error
	switch action {
	case "move":
		err = game.MovePlayer(direction)
		metrics.GuessesTotal.Inc()
	case "shoot":
		err = game.Shoot()
		metrics.GuessesTotal.Inc()
	case "update":
		// Just update game state, no player action
	default:
		return nil, domain.ErrInvalidMove
	}
	
	if err != nil {
		metrics.InvalidGuesses.Inc()
		return nil, err
	}
	
	// Always update game state
	game.Update()
	
	if game.Status == domain.StatusWon {
		metrics.GamesWon.Inc()
		s.logger.WithField("game_id", gameID).Info("Game won")
	} else if game.Status == domain.StatusLost {
		metrics.GamesLost.Inc()
		s.logger.WithField("game_id", gameID).Info("Game lost")
	}
	
	s.logger.WithFields(logrus.Fields{
		"game_id": gameID,
		"action": action,
		"score": game.Score,
	}).Debug("Move processed")
	
	return game, nil
}

func (s *GameService) GetGameStatus(gameID string) (*domain.Game, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	
	game, exists := s.games[gameID]
	if !exists {
		return nil, domain.ErrGameNotFound
	}
	
	return game, nil
}

func (s *GameService) generateGameID() string {
	bytes := make([]byte, 8)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

