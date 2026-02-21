package response

import "github.com/ZyoGo/ayo-indonesia-footbal/internal/club/domain"

type TeamResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	LogoURL     string `json:"logo_url"`
	YearFounded int    `json:"year_founded"`
	Address     string `json:"address"`
	City        string `json:"city"`
}

func FromTeam(team *domain.Team) TeamResponse {
	return TeamResponse{
		ID:          team.ID,
		Name:        team.Name,
		LogoURL:     team.LogoURL,
		YearFounded: team.YearFounded,
		Address:     team.Address,
		City:        team.City,
	}
}

func FromTeams(teams []domain.Team) []TeamResponse {
	result := make([]TeamResponse, len(teams))
	for i, t := range teams {
		result[i] = FromTeam(&t)
	}
	return result
}
