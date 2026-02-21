package response

import "github.com/ZyoGo/ayo-indonesia-footbal/internal/match/domain"

type MatchSummaryResponse struct {
	ID           string `json:"id"`
	HomeTeamID   string `json:"home_team_id"`
	HomeTeamName string `json:"home_team_name"`
	AwayTeamID   string `json:"away_team_id"`
	AwayTeamName string `json:"away_team_name"`
	MatchDate    string `json:"match_date"`
	MatchTime    string `json:"match_time"`
}

type MatchDetailResponse struct {
	ID           string `json:"id"`
	HomeTeamID   string `json:"home_team_id"`
	HomeTeamName string `json:"home_team_name"`
	AwayTeamID   string `json:"away_team_id"`
	AwayTeamName string `json:"away_team_name"`
	MatchDate    string `json:"match_date"`
	MatchTime    string `json:"match_time"`
	Stadium      string `json:"stadium"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

func FromMatch(match *domain.Match) MatchDetailResponse {
	return MatchDetailResponse{
		ID:           match.ID,
		HomeTeamID:   match.HomeTeamID,
		HomeTeamName: match.HomeTeamName,
		AwayTeamID:   match.AwayTeamID,
		AwayTeamName: match.AwayTeamName,
		MatchDate:    match.MatchDate.Format("2006-01-02"),
		MatchTime:    match.MatchTime,
		Stadium:      match.Stadium,
		CreatedAt:    match.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:    match.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func FromMatchSummary(match *domain.Match) MatchSummaryResponse {
	return MatchSummaryResponse{
		ID:           match.ID,
		HomeTeamID:   match.HomeTeamID,
		HomeTeamName: match.HomeTeamName,
		AwayTeamID:   match.AwayTeamID,
		AwayTeamName: match.AwayTeamName,
		MatchDate:    match.MatchDate.Format("2006-01-02"),
		MatchTime:    match.MatchTime,
	}
}

func FromMatches(matches []domain.Match) []MatchSummaryResponse {
	result := make([]MatchSummaryResponse, len(matches))
	for i, m := range matches {
		result[i] = FromMatchSummary(&m)
	}
	return result
}

type MatchReportResponse struct {
	MatchID        string `json:"match_id"`
	MatchDate      string `json:"match_date"`
	MatchTime      string `json:"match_time"`
	HomeTeamName   string `json:"home_team_name"`
	AwayTeamName   string `json:"away_team_name"`
	HomeScore      int    `json:"home_score"`
	AwayScore      int    `json:"away_score"`
	MatchStatus    string `json:"match_status"`
	TopScorer      string `json:"top_scorer"`
	TopScorerGoals int    `json:"top_scorer_goals"`
	HomeTeamWins   int    `json:"home_team_wins"`
	AwayTeamWins   int    `json:"away_team_wins"`
}

func FromMatchReport(report *domain.MatchReportView) MatchReportResponse {
	return MatchReportResponse{
		MatchID:        report.MatchID,
		MatchDate:      report.MatchDate,
		MatchTime:      report.MatchTime,
		HomeTeamName:   report.HomeTeamName,
		AwayTeamName:   report.AwayTeamName,
		HomeScore:      report.HomeScore,
		AwayScore:      report.AwayScore,
		MatchStatus:    report.MatchStatus,
		TopScorer:      report.TopScorer,
		TopScorerGoals: report.TopScorerGoals,
		HomeTeamWins:   report.HomeTeamWins,
		AwayTeamWins:   report.AwayTeamWins,
	}
}

func FromMatchReports(reports []domain.MatchReportView) []MatchReportResponse {
	result := make([]MatchReportResponse, len(reports))
	for i, r := range reports {
		result[i] = FromMatchReport(&r)
	}
	return result
}
