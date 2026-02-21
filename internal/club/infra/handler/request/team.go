package request

import "github.com/ZyoGo/ayo-indonesia-footbal/internal/club/domain"

type CreateTeamRequest struct {
	Name        string `json:"name" binding:"required"`
	LogoURL     string `json:"logo_url"`
	YearFounded int    `json:"year_founded" binding:"required"`
	Address     string `json:"address"`
	City        string `json:"city" binding:"required"`
}

func (r CreateTeamRequest) ToDomain() *domain.Team {
	return &domain.Team{
		Name:        r.Name,
		LogoURL:     r.LogoURL,
		YearFounded: r.YearFounded,
		Address:     r.Address,
		City:        r.City,
	}
}

type UpdateTeamRequest struct {
	Name        string `json:"name" binding:"required"`
	LogoURL     string `json:"logo_url"`
	YearFounded int    `json:"year_founded" binding:"required"`
	Address     string `json:"address"`
	City        string `json:"city" binding:"required"`
}

func (r UpdateTeamRequest) ToDomain() *domain.Team {
	return &domain.Team{
		Name:        r.Name,
		LogoURL:     r.LogoURL,
		YearFounded: r.YearFounded,
		Address:     r.Address,
		City:        r.City,
	}
}
