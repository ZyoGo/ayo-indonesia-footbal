package domain

import "context"

// MatchRepository defines the port for match persistence.
type MatchRepository interface {
	Create(ctx context.Context, match *Match) error
	FindByID(ctx context.Context, id string) (*Match, error)
	FindAll(ctx context.Context) ([]Match, error)
	Update(ctx context.Context, match *Match) error
	Delete(ctx context.Context, id string) error
}

// MatchResultRepository defines the port for match result persistence.
type MatchResultRepository interface {
	Create(ctx context.Context, result *MatchResult) error
	FindByMatchID(ctx context.Context, matchID string) (*MatchResult, error)
	ExistsByMatchID(ctx context.Context, matchID string) (bool, error)
}

// MatchReportView defines the read model for match reports.
type MatchReportView struct {
	MatchID        string
	MatchDate      string
	MatchTime      string
	HomeTeamName   string
	AwayTeamName   string
	HomeScore      int
	AwayScore      int
	MatchStatus    string // "Tim Home Menang" / "Tim Away Menang" / "Draw"
	TopScorer      string // Player name with most goals in this match
	TopScorerGoals int
	HomeTeamWins   int // Accumulated total home team wins
	AwayTeamWins   int // Accumulated total away team wins
}

// ReportRepository defines the port for report queries.
type ReportRepository interface {
	GetMatchReport(ctx context.Context, matchID string) (*MatchReportView, error)
	GetAllMatchReports(ctx context.Context) ([]MatchReportView, error)
}
