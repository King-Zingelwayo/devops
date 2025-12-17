package domain

import (
	"errors"
	"fmt"
	"time"
)

var (
	ErrGameNotFound = errors.New("game not found")
	ErrGameOver     = errors.New("game is over")
	ErrInvalidMove  = errors.New("invalid move")
)

type GameStatus string

const (
	StatusActive GameStatus = "active"
	StatusWon    GameStatus = "won"
	StatusLost   GameStatus = "lost"
)

type Position struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type GameObject struct {
	ID       string   `json:"id"`
	Position Position `json:"position"`
	Active   bool     `json:"active"`
}

type Game struct {
	ID        string       `json:"id"`
	Score     int          `json:"score"`
	Level     int          `json:"level"`
	Player    GameObject   `json:"player"`
	Enemies   []GameObject `json:"enemies"`
	Bullets   []GameObject `json:"bullets"`
	Status    GameStatus   `json:"status"`
	CreatedAt time.Time    `json:"created_at"`
}

func NewGame(id string) *Game {
	return &Game{
		ID:     id,
		Score:  0,
		Level:  1,
		Player: GameObject{ID: "player", Position: Position{X: 385, Y: 550}, Active: true},
		Enemies: []GameObject{
			{ID: "enemy1", Position: Position{X: 100, Y: 50}, Active: true},
			{ID: "enemy2", Position: Position{X: 200, Y: 50}, Active: true},
			{ID: "enemy3", Position: Position{X: 300, Y: 50}, Active: true},
			{ID: "enemy4", Position: Position{X: 400, Y: 50}, Active: true},
			{ID: "enemy5", Position: Position{X: 500, Y: 50}, Active: true},
			{ID: "enemy6", Position: Position{X: 150, Y: 100}, Active: true},
			{ID: "enemy7", Position: Position{X: 250, Y: 100}, Active: true},
			{ID: "enemy8", Position: Position{X: 350, Y: 100}, Active: true},
			{ID: "enemy9", Position: Position{X: 450, Y: 100}, Active: true},
		},
		Bullets:   make([]GameObject, 0),
		Status:    StatusActive,
		CreatedAt: time.Now(),
	}
}

func (g *Game) MovePlayer(direction string) error {
	if g.Status != StatusActive {
		return ErrGameOver
	}

	switch direction {
	case "left":
		if g.Player.Position.X > 0 {
			g.Player.Position.X -= 20
		}
	case "right":
		if g.Player.Position.X < 780 {
			g.Player.Position.X += 20
		}
	default:
		return ErrInvalidMove
	}

	return nil
}

func (g *Game) Shoot() error {
	if g.Status != StatusActive {
		return ErrGameOver
	}

	bullet := GameObject{
		ID:       "bullet",
		Position: Position{X: g.Player.Position.X, Y: g.Player.Position.Y - 10},
		Active:   true,
	}
	g.Bullets = append(g.Bullets, bullet)

	return nil
}

func (g *Game) Update() {
	if g.Status != StatusActive {
		return
	}

	// Move bullets up
	activeBullets := make([]GameObject, 0)
	for i := range g.Bullets {
		if g.Bullets[i].Active {
			g.Bullets[i].Position.Y -= 15
			if g.Bullets[i].Position.Y >= 0 {
				activeBullets = append(activeBullets, g.Bullets[i])
			}
		}
	}
	g.Bullets = activeBullets

	// Move enemies down slowly
	for i := range g.Enemies {
		if g.Enemies[i].Active {
			g.Enemies[i].Position.Y += 1
			// Check if enemies reached bottom
			if g.Enemies[i].Position.Y > 550 {
				g.Status = StatusLost
				return
			}
		}
	}

	// Check collisions
	for i := range g.Bullets {
		for j := range g.Enemies {
			if !g.Enemies[j].Active {
				continue
			}
			if g.checkCollision(g.Bullets[i], g.Enemies[j]) {
				g.Bullets[i].Active = false
				g.Enemies[j].Active = false
				g.Score += 10
				break
			}
		}
	}

	// Remove inactive bullets
	activeBullets = make([]GameObject, 0)
	for _, bullet := range g.Bullets {
		if bullet.Active {
			activeBullets = append(activeBullets, bullet)
		}
	}
	g.Bullets = activeBullets

	// Check win condition
	allEnemiesDestroyed := true
	for _, enemy := range g.Enemies {
		if enemy.Active {
			allEnemiesDestroyed = false
			break
		}
	}
	if allEnemiesDestroyed {
		g.nextLevel()
	}
}

func (g *Game) checkCollision(obj1, obj2 GameObject) bool {
	return abs(obj1.Position.X-obj2.Position.X) < 30 && abs(obj1.Position.Y-obj2.Position.Y) < 30
}

func (g *Game) nextLevel() {
	g.Level++
	g.Score += 50 // Level bonus
	
	// Reset player position
	g.Player.Position = Position{X: 385, Y: 550}
	
	// Create new enemies with increased difficulty
	enemyCount := 9 + g.Level // More enemies each level
	g.Enemies = make([]GameObject, 0)
	
	for i := 0; i < enemyCount; i++ {
		row := i / 5
		col := i % 5
		g.Enemies = append(g.Enemies, GameObject{
			ID:       fmt.Sprintf("enemy%d", i+1),
			Position: Position{X: 100 + col*100, Y: 50 + row*50},
			Active:   true,
		})
	}
	
	// Clear bullets
	g.Bullets = make([]GameObject, 0)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}