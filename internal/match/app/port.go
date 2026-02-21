package app

import (
	"context"

	"github.com/ZyoGo/ayo-indonesia-footbal/internal/match/domain"
)

//go:generate mockgen -source=port.go -destination=mock/mock_service.go -package=mock

// MatchServicePort defines the contract for match business operations.
type MatchServicePort interface {
	CreateMatch(ctx context.Context, match *domain.Match) (string, error)
	GetMatchByID(ctx context.Context, id string) (*domain.Match, error)
	GetAllMatches(ctx context.Context) ([]domain.Match, error)
	ReportResult(ctx context.Context, matchID string, result *domain.MatchResult) (string, error)
	GetMatchReport(ctx context.Context, matchID string) (*domain.MatchReportView, error)
	GetAllMatchReports(ctx context.Context) ([]domain.MatchReportView, error)
	DeleteMatch(ctx context.Context, id string) error
}
