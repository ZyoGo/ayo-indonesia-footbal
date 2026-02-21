package app

import (
	"context"

	"github.com/ZyoGo/ayo-indonesia-footbal/internal/club/domain"
	"github.com/ZyoGo/ayo-indonesia-footbal/pkg/derrors"
)

type PlayerService struct {
	playerRepo domain.PlayerRepository
	teamRepo   domain.TeamRepository
}

func NewPlayerService(playerRepo domain.PlayerRepository, teamRepo domain.TeamRepository) PlayerServicePort {
	return &PlayerService{
		playerRepo: playerRepo,
		teamRepo:   teamRepo,
	}
}

func (s *PlayerService) Create(ctx context.Context, player *domain.Player) (string, error) {
	// Verify team exists
	_, err := s.teamRepo.FindByID(ctx, player.TeamID)
	if err != nil {
		return "", err
	}

	// Check jersey number uniqueness
	taken, err := s.playerRepo.IsJerseyNumberTaken(ctx, player.TeamID, player.JerseyNumber, "")
	if err != nil {
		return "", err
	}
	if taken {
		return "", derrors.WrapErrorf(domain.ErrJerseyNumberTaken, derrors.ErrorCodeDuplicate, "jersey number %d is already taken", player.JerseyNumber)
	}

	// Construct valid entity via domain factory
	newPlayer, err := domain.NewPlayer(player.TeamID, player.Name, player.Height, player.Weight, player.Position, player.JerseyNumber)
	if err != nil {
		return "", err
	}

	if err := s.playerRepo.Create(ctx, newPlayer); err != nil {
		return "", err
	}

	return newPlayer.ID, nil
}

func (s *PlayerService) GetByID(ctx context.Context, id string) (*domain.Player, error) {
	player, err := s.playerRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return player, nil
}

func (s *PlayerService) GetByTeamID(ctx context.Context, teamID string) ([]domain.Player, error) {
	_, err := s.teamRepo.FindByID(ctx, teamID)
	if err != nil {
		return nil, err
	}

	players, err := s.playerRepo.FindByTeamID(ctx, teamID)
	if err != nil {
		return nil, err
	}
	return players, nil
}

func (s *PlayerService) Update(ctx context.Context, id string, player *domain.Player) error {
	existing, err := s.playerRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	// Check jersey number uniqueness (exclude current player)
	taken, err := s.playerRepo.IsJerseyNumberTaken(ctx, existing.TeamID, player.JerseyNumber, existing.ID)
	if err != nil {
		return err
	}
	if taken {
		return derrors.WrapErrorf(domain.ErrJerseyNumberTaken, derrors.ErrorCodeDuplicate, "jersey number %d is already taken", player.JerseyNumber)
	}

	if err := existing.Update(player.Name, player.Height, player.Weight, player.Position, player.JerseyNumber); err != nil {
		return err
	}

	if err := s.playerRepo.Update(ctx, existing); err != nil {
		return err
	}

	return nil
}

func (s *PlayerService) Delete(ctx context.Context, id string) error {
	_, err := s.playerRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if err := s.playerRepo.SoftDelete(ctx, id); err != nil {
		return err
	}

	return nil
}
