package app

import (
	"context"

	"github.com/ZyoGo/ayo-indonesia-footbal/internal/club/domain"
)

type TeamService struct {
	teamRepo domain.TeamRepository
}

func NewTeamService(teamRepo domain.TeamRepository) TeamServicePort {
	return &TeamService{teamRepo: teamRepo}
}

func (s *TeamService) Create(ctx context.Context, team *domain.Team) (string, error) {
	newTeam, err := domain.NewTeam(team.Name, team.LogoURL, team.YearFounded, team.Address, team.City)
	if err != nil {
		return "", err
	}

	if err := s.teamRepo.Create(ctx, newTeam); err != nil {
		return "", err
	}

	return newTeam.ID, nil
}

func (s *TeamService) GetByID(ctx context.Context, id string) (*domain.Team, error) {
	team, err := s.teamRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return team, nil
}

func (s *TeamService) GetAll(ctx context.Context) ([]domain.Team, error) {
	teams, err := s.teamRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	return teams, nil
}

func (s *TeamService) Update(ctx context.Context, id string, team *domain.Team) error {
	existing, err := s.teamRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if err := existing.Update(team.Name, team.LogoURL, team.YearFounded, team.Address, team.City); err != nil {
		return err
	}

	if err := s.teamRepo.Update(ctx, existing); err != nil {
		return err
	}

	return nil
}

func (s *TeamService) Delete(ctx context.Context, id string) error {
	_, err := s.teamRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if err := s.teamRepo.SoftDelete(ctx, id); err != nil {
		return err
	}

	return nil
}
