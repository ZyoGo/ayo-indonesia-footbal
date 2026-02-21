package app

import (
	"context"
	"github.com/ZyoGo/ayo-indonesia-footbal/internal/reporting/domain"
)

type ReportingServicePort interface {
	GetStandings(ctx context.Context) ([]domain.TeamStanding, error)
	GetTopScorers(ctx context.Context) ([]domain.TopScorer, error)
}
