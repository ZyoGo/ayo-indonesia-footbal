package response

import "github.com/ZyoGo/ayo-indonesia-footbal/internal/reporting/domain"

type StandingResponse struct {
	TeamID   string `json:"team_id"`
	TeamName string `json:"team_name"`
	Played   int    `json:"played"`
	Won      int    `json:"won"`
	Drawn    int    `json:"drawn"`
	Lost     int    `json:"lost"`
	GF       int    `json:"gf"`
	GA       int    `json:"ga"`
	GD       int    `json:"gd"`
	Points   int    `json:"points"`
}

type TopScorerResponse struct {
	PlayerID   string `json:"player_id"`
	PlayerName string `json:"player_name"`
	TeamName   string `json:"team_name"`
	Goals      int    `json:"goals"`
}

func FromStandingDomain(d domain.TeamStanding) StandingResponse {
	return StandingResponse{
		TeamID:   d.TeamID,
		TeamName: d.TeamName,
		Played:   d.Played,
		Won:      d.Won,
		Drawn:    d.Drawn,
		Lost:     d.Lost,
		GF:       d.GF,
		GA:       d.GA,
		GD:       d.GD,
		Points:   d.Points,
	}
}

func FromTopScorerDomain(d domain.TopScorer) TopScorerResponse {
	return TopScorerResponse{
		PlayerID:   d.PlayerID,
		PlayerName: d.PlayerName,
		TeamName:   d.TeamName,
		Goals:      d.Goals,
	}
}
