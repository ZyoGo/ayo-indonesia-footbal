package domain

import "context"

type TeamStanding struct {
	TeamID   string
	TeamName string
	Played   int
	Won      int
	Drawn    int
	Lost     int
	GF       int // Goals For
	GA       int // Goals Against
	GD       int // Goal Difference
	Points   int
}

type TopScorer struct {
	PlayerID   string
	PlayerName string
	TeamName   string
	Goals      int
}

type ReportingRepository interface {
	GetStandings(ctx context.Context) ([]TeamStanding, error)
	GetTopScorers(ctx context.Context) ([]TopScorer, error)
}
