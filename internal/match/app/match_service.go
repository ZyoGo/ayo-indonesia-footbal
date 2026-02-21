package app

import (
	"context"

	"github.com/ZyoGo/ayo-indonesia-footbal/internal/match/domain"
	"github.com/ZyoGo/ayo-indonesia-footbal/pkg/derrors"
)

type MatchService struct {
	matchRepo  domain.MatchRepository
	resultRepo domain.MatchResultRepository
	reportRepo domain.ReportRepository
}

func NewMatchService(
	matchRepo domain.MatchRepository,
	resultRepo domain.MatchResultRepository,
	reportRepo domain.ReportRepository,
) MatchServicePort {
	return &MatchService{
		matchRepo:  matchRepo,
		resultRepo: resultRepo,
		reportRepo: reportRepo,
	}
}

func (s *MatchService) CreateMatch(ctx context.Context, match *domain.Match) (string, error) {
	newMatch, err := domain.NewMatch(match.HomeTeamID, match.AwayTeamID, match.MatchDate, match.MatchTime, match.Stadium)
	if err != nil {
		return "", err
	}

	if err := s.matchRepo.Create(ctx, newMatch); err != nil {
		return "", err
	}

	return newMatch.ID, nil
}

func (s *MatchService) GetMatchByID(ctx context.Context, id string) (*domain.Match, error) {
	match, err := s.matchRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return match, nil
}

func (s *MatchService) GetAllMatches(ctx context.Context) ([]domain.Match, error) {
	matches, err := s.matchRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	return matches, nil
}

func (s *MatchService) ReportResult(ctx context.Context, matchID string, result *domain.MatchResult) (string, error) {
	// Verify match exists
	m, err := s.matchRepo.FindByID(ctx, matchID)
	if err != nil {
		return "", err
	}

	// Check if result already reported
	exists, err := s.resultRepo.ExistsByMatchID(ctx, matchID)
	if err != nil {
		return "", err
	}
	if exists {
		return "", derrors.WrapErrorf(domain.ErrResultAlreadyExists, derrors.ErrorCodeDuplicate, "%s", domain.ErrResultAlreadyExists.Error())
	}

	// Construct valid result via domain factory
	newResult, err := domain.NewMatchResult(matchID, m.HomeTeamID, m.AwayTeamID, result.HomeScore, result.AwayScore, result.Goals)
	if err != nil {
		return "", err
	}

	if err := s.resultRepo.Create(ctx, newResult); err != nil {
		return "", err
	}

	return newResult.ID, nil
}

func (s *MatchService) GetMatchReport(ctx context.Context, matchID string) (*domain.MatchReportView, error) {
	report, err := s.reportRepo.GetMatchReport(ctx, matchID)
	if err != nil {
		return nil, err
	}
	return report, nil
}

func (s *MatchService) GetAllMatchReports(ctx context.Context) ([]domain.MatchReportView, error) {
	reports, err := s.reportRepo.GetAllMatchReports(ctx)
	if err != nil {
		return nil, err
	}
	return reports, nil
}

func (s *MatchService) DeleteMatch(ctx context.Context, id string) error {
	return s.matchRepo.Delete(ctx, id)
}
