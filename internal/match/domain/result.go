package domain

import (
	"strings"
	"time"


	"github.com/ZyoGo/ayo-indonesia-footbal/pkg/derrors"
	"github.com/ZyoGo/ayo-indonesia-footbal/pkg/ulid"
)

const (
	maxGoalMinute = 150 // including extra time and penalties
)

type MatchResult struct {
	ID        string
	MatchID   string
	HomeScore int
	AwayScore int
	Goals     []Goal
	DeletedAt *time.Time
}

type Goal struct {
	ID         string
	ResultID   string
	PlayerID   string
	PlayerName string
	TeamID     string
	GoalMinute int
	DeletedAt  *time.Time
}

func NewMatchResult(matchID string, homeTeamID, awayTeamID string, homeScore, awayScore int, goals []Goal) (*MatchResult, error) {
	matchID = strings.TrimSpace(matchID)
	homeTeamID = strings.TrimSpace(homeTeamID)
	awayTeamID = strings.TrimSpace(awayTeamID)

	if matchID == "" {
		return nil, derrors.NewErrorf(derrors.ErrorCodeBadRequest, "match ID is required")
	}
	if homeTeamID == "" || awayTeamID == "" {
		return nil, derrors.NewErrorf(derrors.ErrorCodeBadRequest, "home and away team IDs are required")
	}
	if homeScore < 0 || awayScore < 0 {
		return nil, derrors.NewErrorf(derrors.ErrorCodeBadRequest, "score cannot be negative")
	}

	// Validate that total goals match the reported score
	totalGoals := len(goals)
	expectedTotal := homeScore + awayScore
	if totalGoals != expectedTotal {
		return nil, derrors.NewErrorf(derrors.ErrorCodeBadRequest, "total goals (%d) does not match reported score (%d)", totalGoals, expectedTotal)
	}

	homeGoalCount := 0
	awayGoalCount := 0

	resultID := ulid.GenerateID()
	for i := range goals {
		goals[i].ID = ulid.GenerateID()
		goals[i].ResultID = resultID
		goals[i].PlayerID = strings.TrimSpace(goals[i].PlayerID)
		goals[i].TeamID = strings.TrimSpace(goals[i].TeamID)

		if goals[i].PlayerID == "" {
			return nil, derrors.NewErrorf(derrors.ErrorCodeBadRequest, "player ID is required for each goal")
		}
		if goals[i].TeamID == "" {
			return nil, derrors.NewErrorf(derrors.ErrorCodeBadRequest, "team ID is required for each goal")
		}
		if goals[i].GoalMinute <= 0 || goals[i].GoalMinute > maxGoalMinute {
			return nil, derrors.NewErrorf(derrors.ErrorCodeBadRequest, "goal minute must be between 1 and %d", maxGoalMinute)
		}

		if goals[i].TeamID == homeTeamID {
			homeGoalCount++
		} else if goals[i].TeamID == awayTeamID {
			awayGoalCount++
		} else {
			return nil, derrors.NewErrorf(derrors.ErrorCodeBadRequest, "goal team ID %s does not belong to match participants", goals[i].TeamID)
		}
	}

	if homeGoalCount != homeScore {
		return nil, derrors.NewErrorf(derrors.ErrorCodeBadRequest, "home goals (%d) does not match home score (%d)", homeGoalCount, homeScore)
	}
	if awayGoalCount != awayScore {
		return nil, derrors.NewErrorf(derrors.ErrorCodeBadRequest, "away goals (%d) does not match away score (%d)", awayGoalCount, awayScore)
	}

	return &MatchResult{
		ID:        resultID,
		MatchID:   matchID,
		HomeScore: homeScore,
		AwayScore: awayScore,
		Goals:     goals,
	}, nil
}

// Status returns the match outcome based on the score.
func (r *MatchResult) Status() string {
	if r.HomeScore > r.AwayScore {
		return "Home Win"
	}
	if r.AwayScore > r.HomeScore {
		return "Away Win"
	}
	return "Draw"
}
