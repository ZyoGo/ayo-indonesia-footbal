package domain

import (
	"regexp"
	"strings"
	"time"

	"github.com/ZyoGo/ayo-indonesia-footbal/pkg/derrors"
	"github.com/ZyoGo/ayo-indonesia-footbal/pkg/ulid"
)

var matchTimeRegex = regexp.MustCompile(`^([01]\d|2[0-3]):([0-5]\d)$`)

type Match struct {
	ID           string
	HomeTeamID   string
	AwayTeamID   string
	MatchDate    time.Time
	MatchTime    string // HH:MM format
	Stadium      string
	HomeTeamName string // Populated on read
	AwayTeamName string // Populated on read
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time
}

func NewMatch(homeTeamID, awayTeamID string, matchDate time.Time, matchTime string, stadium string) (*Match, error) {
	homeTeamID = strings.TrimSpace(homeTeamID)
	awayTeamID = strings.TrimSpace(awayTeamID)
	matchTime = strings.TrimSpace(matchTime)
	stadium = strings.TrimSpace(stadium)

	if homeTeamID == "" {
		return nil, derrors.NewErrorf(derrors.ErrorCodeBadRequest, "home team ID is required")
	}
	if awayTeamID == "" {
		return nil, derrors.NewErrorf(derrors.ErrorCodeBadRequest, "away team ID is required")
	}
	if homeTeamID == awayTeamID {
		return nil, derrors.NewErrorf(derrors.ErrorCodeBadRequest, "home team and away team cannot be the same")
	}
	if matchTime == "" {
		return nil, derrors.NewErrorf(derrors.ErrorCodeBadRequest, "match time is required")
	}
	if !matchTimeRegex.MatchString(matchTime) {
		return nil, derrors.NewErrorf(derrors.ErrorCodeBadRequest, "match time must be in HH:MM format (00:00 - 23:59)")
	}
	if stadium == "" {
		return nil, derrors.NewErrorf(derrors.ErrorCodeBadRequest, "stadium is required")
	}
	now := time.Now()
	// Allow matches on the same day (today) or in the future
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	if matchDate.Before(today) {
		return nil, derrors.NewErrorf(derrors.ErrorCodeBadRequest, "match date cannot be in the past")
	}

	return &Match{
		ID:         ulid.GenerateID(),
		HomeTeamID: homeTeamID,
		AwayTeamID: awayTeamID,
		MatchDate:  matchDate,
		MatchTime:  matchTime,
		Stadium:    stadium,
		CreatedAt:  now,
		UpdatedAt:  now,
	}, nil
}
