package response

import "github.com/ZyoGo/ayo-indonesia-footbal/internal/match/domain"

type MatchResponse struct {
	ID         string `json:"id"`
	HomeTeamID string `json:"home_team_id"`
	AwayTeamID string `json:"away_team_id"`
	MatchDate  string `json:"match_date"`
	MatchTime  string `json:"match_time"`
}

func FromMatch(match *domain.Match) MatchResponse {
	return MatchResponse{
		ID:         match.ID,
		HomeTeamID: match.HomeTeamID,
		AwayTeamID: match.AwayTeamID,
		MatchDate:  match.MatchDate.Format("2006-01-02"),
		MatchTime:  match.MatchTime,
	}
}

func FromMatches(matches []domain.Match) []MatchResponse {
	result := make([]MatchResponse, len(matches))
	for i, m := range matches {
		result[i] = FromMatch(&m)
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
