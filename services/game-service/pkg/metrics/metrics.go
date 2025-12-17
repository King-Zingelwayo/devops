package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// Low-cardinality metrics only - no user IDs or game IDs in labels
	GamesStarted = promauto.NewCounter(prometheus.CounterOpts{
		Name: "wordle_games_started_total",
		Help: "Total number of games started",
	})

	GamesWon = promauto.NewCounter(prometheus.CounterOpts{
		Name: "wordle_games_won_total",
		Help: "Total number of games won",
	})

	GamesLost = promauto.NewCounter(prometheus.CounterOpts{
		Name: "wordle_games_lost_total",
		Help: "Total number of games lost",
	})

	GuessesTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "wordle_guesses_total",
		Help: "Total number of guesses made",
	})

	InvalidGuesses = promauto.NewCounter(prometheus.CounterOpts{
		Name: "wordle_invalid_guesses_total",
		Help: "Total number of invalid guesses",
	})

	ActiveGames = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "wordle_active_games",
		Help: "Number of currently active games",
	})
)

func Init() {
	// Register custom metrics if needed
	// All metrics are auto-registered via promauto
}