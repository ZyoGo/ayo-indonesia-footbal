package postgres

import (
	"context"
	"github.com/ZyoGo/ayo-indonesia-footbal/internal/reporting/domain"
	"github.com/ZyoGo/ayo-indonesia-footbal/pkg/derrors"
	"github.com/jackc/pgx/v5/pgxpool"
)

type reportingRepository struct {
	db *pgxpool.Pool
}

func NewReportingRepository(db *pgxpool.Pool) domain.ReportingRepository {
	return &reportingRepository{db: db}
}

func (r *reportingRepository) GetStandings(ctx context.Context) ([]domain.TeamStanding, error) {
	rows, err := r.db.Query(ctx, queryStandings)
	if err != nil {
		return nil, derrors.WrapErrorf(err, derrors.ErrorCodeInternal, "failed to query standings")
	}
	defer rows.Close()

	var standings []domain.TeamStanding
	for rows.Next() {
		var s domain.TeamStanding
		if err := rows.Scan(
			&s.TeamID,
			&s.TeamName,
			&s.Played,
			&s.Won,
			&s.Drawn,
			&s.Lost,
			&s.GF,
			&s.GA,
			&s.GD,
			&s.Points,
		); err != nil {
			return nil, derrors.WrapErrorf(err, derrors.ErrorCodeInternal, "failed to scan standing row")
		}
		standings = append(standings, s)
	}

	return standings, nil
}

func (r *reportingRepository) GetTopScorers(ctx context.Context) ([]domain.TopScorer, error) {
	rows, err := r.db.Query(ctx, queryTopScorers)
	if err != nil {
		return nil, derrors.WrapErrorf(err, derrors.ErrorCodeInternal, "failed to query top scorers")
	}
	defer rows.Close()

	var scorers []domain.TopScorer
	for rows.Next() {
		var s domain.TopScorer
		if err := rows.Scan(
			&s.PlayerID,
			&s.PlayerName,
			&s.TeamName,
			&s.Goals,
		); err != nil {
			return nil, derrors.WrapErrorf(err, derrors.ErrorCodeInternal, "failed to scan top scorer row")
		}
		scorers = append(scorers, s)
	}

	return scorers, nil
}
