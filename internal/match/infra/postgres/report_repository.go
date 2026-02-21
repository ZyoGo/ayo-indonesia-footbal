package postgres

import (
	"context"
	"errors"

	"github.com/ZyoGo/ayo-indonesia-footbal/internal/match/domain"
	"github.com/ZyoGo/ayo-indonesia-footbal/pkg/derrors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type reportRepository struct {
	db *pgxpool.Pool
}

func NewReportRepository(db *pgxpool.Pool) domain.ReportRepository {
	return &reportRepository{db: db}
}

func (r *reportRepository) GetMatchReport(ctx context.Context, matchID string) (*domain.MatchReportView, error) {
	var report domain.MatchReportView
	err := r.db.QueryRow(ctx, queryMatchReport, matchID).Scan(
		&report.MatchID,
		&report.MatchDate,
		&report.MatchTime,
		&report.HomeTeamName,
		&report.AwayTeamName,
		&report.HomeScore,
		&report.AwayScore,
		&report.MatchStatus,
		&report.TopScorer,
		&report.TopScorerGoals,
		&report.HomeTeamWins,
		&report.AwayTeamWins,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, derrors.WrapErrorf(domain.ErrMatchResultNotFound, derrors.ErrorCodeNotFound, "%s", domain.ErrMatchResultNotFound.Error())
		}
		return nil, derrors.WrapErrorf(err, derrors.ErrorCodeInternal, "failed to get match report")
	}
	return &report, nil
}

func (r *reportRepository) GetAllMatchReports(ctx context.Context) ([]domain.MatchReportView, error) {
	rows, err := r.db.Query(ctx, queryAllMatchReports)
	if err != nil {
		return nil, derrors.WrapErrorf(err, derrors.ErrorCodeInternal, "failed to query match reports")
	}
	defer rows.Close()

	var reports []domain.MatchReportView
	for rows.Next() {
		var report domain.MatchReportView
		if err := rows.Scan(
			&report.MatchID,
			&report.MatchDate,
			&report.MatchTime,
			&report.HomeTeamName,
			&report.AwayTeamName,
			&report.HomeScore,
			&report.AwayScore,
			&report.MatchStatus,
			&report.TopScorer,
			&report.TopScorerGoals,
			&report.HomeTeamWins,
			&report.AwayTeamWins,
		); err != nil {
			return nil, derrors.WrapErrorf(err, derrors.ErrorCodeInternal, "failed to scan match report row")
		}
		reports = append(reports, report)
	}

	return reports, nil
}
