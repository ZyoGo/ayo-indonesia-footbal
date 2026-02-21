package request

import (
	"time"

	"github.com/ZyoGo/ayo-indonesia-footbal/internal/match/domain"
)

type CreateMatchRequest struct {
	HomeTeamID string `json:"home_team_id" binding:"required"`
	AwayTeamID string `json:"away_team_id" binding:"required"`
	MatchDate  string `json:"match_date" binding:"required"` // YYYY-MM-DD
	MatchTime  string `json:"match_time" binding:"required"` // HH:MM
}

func (r CreateMatchRequest) ToDomain() *domain.Match {
	date, _ := time.Parse("2006-01-02", r.MatchDate)
	return &domain.Match{
		HomeTeamID: r.HomeTeamID,
		AwayTeamID: r.AwayTeamID,
		MatchDate:  date,
		MatchTime:  r.MatchTime,
	}
}

type ReportResultRequest struct {
	HomeScore int          `json:"home_score" binding:"min=0"`
	AwayScore int          `json:"away_score" binding:"min=0"`
	Goals     []GoalInput  `json:"goals" binding:"required"`
}

type GoalInput struct {
	PlayerID   string `json:"player_id" binding:"required"`
	TeamID     string `json:"team_id" binding:"required"`
	GoalMinute int    `json:"goal_minute" binding:"required"`
}

func (r ReportResultRequest) ToDomain() *domain.MatchResult {
	goals := make([]domain.Goal, len(r.Goals))
	for i, g := range r.Goals {
		goals[i] = domain.Goal{
			PlayerID:   g.PlayerID,
			TeamID:     g.TeamID,
			GoalMinute: g.GoalMinute,
		}
	}
	return &domain.MatchResult{
		HomeScore: r.HomeScore,
		AwayScore: r.AwayScore,
		Goals:     goals,
	}
}
