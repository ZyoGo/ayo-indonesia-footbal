package request

import "github.com/ZyoGo/ayo-indonesia-footbal/internal/club/domain"

type CreatePlayerRequest struct {
	TeamID       string  `json:"team_id" binding:"required"`
	Name         string  `json:"name" binding:"required"`
	Height       float64 `json:"height" binding:"required"`
	Weight       float64 `json:"weight" binding:"required"`
	Position     string  `json:"position" binding:"required"`
	JerseyNumber int     `json:"jersey_number" binding:"required"`
}

func (r CreatePlayerRequest) ToDomain() *domain.Player {
	pos, _ := domain.ParsePosition(r.Position)
	return &domain.Player{
		TeamID:       r.TeamID,
		Name:         r.Name,
		Height:       r.Height,
		Weight:       r.Weight,
		Position:     pos,
		JerseyNumber: r.JerseyNumber,
	}
}

type UpdatePlayerRequest struct {
	Name         string  `json:"name" binding:"required"`
	Height       float64 `json:"height" binding:"required"`
	Weight       float64 `json:"weight" binding:"required"`
	Position     string  `json:"position" binding:"required"`
	JerseyNumber int     `json:"jersey_number" binding:"required"`
}

func (r UpdatePlayerRequest) ToDomain() *domain.Player {
	pos, _ := domain.ParsePosition(r.Position)
	return &domain.Player{
		Name:         r.Name,
		Height:       r.Height,
		Weight:       r.Weight,
		Position:     pos,
		JerseyNumber: r.JerseyNumber,
	}
}
