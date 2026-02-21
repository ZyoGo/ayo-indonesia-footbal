package app

import (
	"context"
	"github.com/ZyoGo/ayo-indonesia-footbal/internal/reporting/domain"
)

type ReportingService struct {
	repo domain.ReportingRepository
}

func NewReportingService(repo domain.ReportingRepository) ReportingServicePort {
	return &ReportingService{repo: repo}
}

func (s *ReportingService) GetStandings(ctx context.Context) ([]domain.TeamStanding, error) {
	return s.repo.GetStandings(ctx)
}

func (s *ReportingService) GetTopScorers(ctx context.Context) ([]domain.TopScorer, error) {
	return s.repo.GetTopScorers(ctx)
}
