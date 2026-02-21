package domain

import (
	"strings"
	"time"

	"github.com/ZyoGo/ayo-indonesia-footbal/pkg/derrors"
	"github.com/ZyoGo/ayo-indonesia-footbal/pkg/ulid"
)

const (
	maxTeamNameLength = 255
	maxCityLength     = 100
)

type Team struct {
	ID          string
	Name        string
	LogoURL     string
	YearFounded int
	Address     string
	City        string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}

func NewTeam(name, logoURL string, yearFounded int, address, city string) (*Team, error) {
	name = strings.TrimSpace(name)
	city = strings.TrimSpace(city)
	logoURL = strings.TrimSpace(logoURL)
	address = strings.TrimSpace(address)

	if name == "" {
		return nil, derrors.NewErrorf(derrors.ErrorCodeBadRequest, "team name is required")
	}
	if len(name) > maxTeamNameLength {
		return nil, derrors.NewErrorf(derrors.ErrorCodeBadRequest, "team name must not exceed %d characters", maxTeamNameLength)
	}
	if city == "" {
		return nil, derrors.NewErrorf(derrors.ErrorCodeBadRequest, "team city is required")
	}
	if len(city) > maxCityLength {
		return nil, derrors.NewErrorf(derrors.ErrorCodeBadRequest, "team city must not exceed %d characters", maxCityLength)
	}
	if yearFounded <= 0 {
		return nil, derrors.NewErrorf(derrors.ErrorCodeBadRequest, "year founded must be a positive number")
	}
	if yearFounded > time.Now().Year() {
		return nil, derrors.NewErrorf(derrors.ErrorCodeBadRequest, "year founded cannot be in the future")
	}

	now := time.Now()
	return &Team{
		ID:          ulid.GenerateID(),
		Name:        name,
		LogoURL:     logoURL,
		YearFounded: yearFounded,
		Address:     address,
		City:        city,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

func (t *Team) Update(name, logoURL string, yearFounded int, address, city string) error {
	name = strings.TrimSpace(name)
	city = strings.TrimSpace(city)
	logoURL = strings.TrimSpace(logoURL)
	address = strings.TrimSpace(address)

	if name == "" {
		return derrors.NewErrorf(derrors.ErrorCodeBadRequest, "team name is required")
	}
	if len(name) > maxTeamNameLength {
		return derrors.NewErrorf(derrors.ErrorCodeBadRequest, "team name must not exceed %d characters", maxTeamNameLength)
	}
	if city == "" {
		return derrors.NewErrorf(derrors.ErrorCodeBadRequest, "team city is required")
	}
	if len(city) > maxCityLength {
		return derrors.NewErrorf(derrors.ErrorCodeBadRequest, "team city must not exceed %d characters", maxCityLength)
	}
	if yearFounded <= 0 {
		return derrors.NewErrorf(derrors.ErrorCodeBadRequest, "year founded must be a positive number")
	}
	if yearFounded > time.Now().Year() {
		return derrors.NewErrorf(derrors.ErrorCodeBadRequest, "year founded cannot be in the future")
	}

	t.Name = name
	t.LogoURL = logoURL
	t.YearFounded = yearFounded
	t.Address = address
	t.City = city
	t.UpdatedAt = time.Now()
	return nil
}
