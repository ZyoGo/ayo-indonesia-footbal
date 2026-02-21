package response

import "github.com/ZyoGo/ayo-indonesia-footbal/internal/club/domain"

type PlayerResponse struct {
	ID           string  `json:"id"`
	TeamID       string  `json:"team_id"`
	Name         string  `json:"name"`
	Height       float64 `json:"height"`
	Weight       float64 `json:"weight"`
	Position     string  `json:"position"`
	JerseyNumber int     `json:"jersey_number"`
}

func FromPlayer(player *domain.Player) PlayerResponse {
	return PlayerResponse{
		ID:           player.ID,
		TeamID:       player.TeamID,
		Name:         player.Name,
		Height:       player.Height,
		Weight:       player.Weight,
		Position:     player.Position.String(),
		JerseyNumber: player.JerseyNumber,
	}
}

func FromPlayers(players []domain.Player) []PlayerResponse {
	result := make([]PlayerResponse, len(players))
	for i, p := range players {
		result[i] = FromPlayer(&p)
	}
	return result
}
