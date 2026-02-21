package postgres

import (
	"context"
	"errors"

	"github.com/ZyoGo/ayo-indonesia-footbal/internal/match/domain"
	"github.com/ZyoGo/ayo-indonesia-footbal/pkg/derrors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type matchResultRepository struct {
	db *pgxpool.Pool
}

func NewMatchResultRepository(db *pgxpool.Pool) domain.MatchResultRepository {
	return &matchResultRepository{db: db}
}

func (r *matchResultRepository) Create(ctx context.Context, result *domain.MatchResult) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return derrors.WrapErrorf(err, derrors.ErrorCodeInternal, "failed to begin transaction")
	}
	defer tx.Rollback(ctx)

	// Insert match result
	_, err = tx.Exec(ctx, queryInsertMatchResult,
		result.ID,
		result.MatchID,
		result.HomeScore,
		result.AwayScore,
	)
	if err != nil {
		return derrors.WrapErrorf(err, derrors.ErrorCodeInternal, "failed to insert match result")
	}

	// Insert each goal event
	for _, goal := range result.Goals {
		_, err = tx.Exec(ctx, queryInsertGoal,
			goal.ID,
			goal.ResultID,
			goal.PlayerID,
			goal.TeamID,
			goal.GoalMinute,
		)
		if err != nil {
			return derrors.WrapErrorf(err, derrors.ErrorCodeInternal, "failed to insert goal event")
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return derrors.WrapErrorf(err, derrors.ErrorCodeInternal, "failed to commit transaction")
	}

	return nil
}

func (r *matchResultRepository) FindByMatchID(ctx context.Context, matchID string) (*domain.MatchResult, error) {
	var result domain.MatchResult
	err := r.db.QueryRow(ctx, queryFindResultByMatchID, matchID).Scan(
		&result.ID,
		&result.MatchID,
		&result.HomeScore,
		&result.AwayScore,
		&result.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, derrors.WrapErrorf(domain.ErrMatchResultNotFound, derrors.ErrorCodeNotFound, "%s", domain.ErrMatchResultNotFound.Error())
		}
		return nil, derrors.WrapErrorf(err, derrors.ErrorCodeInternal, "failed to find match result")
	}

	// Fetch goals
	rows, err := r.db.Query(ctx, queryFindGoalsByResultID, result.ID)
	if err != nil {
		return nil, derrors.WrapErrorf(err, derrors.ErrorCodeInternal, "failed to query goals")
	}
	defer rows.Close()

	for rows.Next() {
		var goal domain.Goal
		if err := rows.Scan(
			&goal.ID,
			&goal.ResultID,
			&goal.PlayerID,
			&goal.PlayerName,
			&goal.TeamID,
			&goal.GoalMinute,
			&goal.DeletedAt,
		); err != nil {
			return nil, derrors.WrapErrorf(err, derrors.ErrorCodeInternal, "failed to scan goal row")
		}
		result.Goals = append(result.Goals, goal)
	}

	return &result, nil
}

func (r *matchResultRepository) ExistsByMatchID(ctx context.Context, matchID string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx, queryExistsResultByMatchID, matchID).Scan(&exists)
	if err != nil {
		return false, derrors.WrapErrorf(err, derrors.ErrorCodeInternal, "failed to check match result existence")
	}
	return exists, nil
}
