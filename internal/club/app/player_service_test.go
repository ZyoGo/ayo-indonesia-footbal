package app

import (
	"context"
	"errors"
	"testing"

	"github.com/ZyoGo/ayo-indonesia-footbal/internal/club/domain"
	mockDomain "github.com/ZyoGo/ayo-indonesia-footbal/internal/club/mock"
	"github.com/ZyoGo/ayo-indonesia-footbal/pkg/derrors"
	"go.uber.org/mock/gomock"
)

func setupPlayerService(t *testing.T) (*PlayerService, *mockDomain.MockPlayerRepository, *mockDomain.MockTeamRepository) {
	t.Helper()
	ctrl := gomock.NewController(t)
	mockPlayerRepo := mockDomain.NewMockPlayerRepository(ctrl)
	mockTeamRepo := mockDomain.NewMockTeamRepository(ctrl)
	svc := &PlayerService{
		playerRepo: mockPlayerRepo,
		teamRepo:   mockTeamRepo,
	}
	return svc, mockPlayerRepo, mockTeamRepo
}

// assertPlayerErrorCode verifies the error is a *derrors.Error with the expected code.
func assertPlayerErrorCode(t *testing.T, err error, expectedCode derrors.ErrorCode) {
	t.Helper()
	var dErr *derrors.Error
	if !errors.As(err, &dErr) {
		t.Fatalf("expected *derrors.Error, got %T: %v", err, err)
	}
	if dErr.Code() != expectedCode {
		t.Fatalf("expected error code %d, got %d", expectedCode, dErr.Code())
	}
}

// ---------------------------------------------------------------------------
// Create
// ---------------------------------------------------------------------------

func TestPlayerService_Create_Success(t *testing.T) {
	// Given
	svc, mockPlayerRepo, mockTeamRepo := setupPlayerService(t)
	ctx := context.Background()
	input := &domain.Player{
		TeamID:       "team-1",
		Name:         "Beckham",
		Height:       180.0,
		Weight:       75.0,
		Position:     domain.PositionCM,
		JerseyNumber: 10,
	}

	mockTeamRepo.EXPECT().FindByID(ctx, "team-1").Return(&domain.Team{ID: "team-1"}, nil)
	mockPlayerRepo.EXPECT().IsJerseyNumberTaken(ctx, "team-1", 10, "").Return(false, nil)
	mockPlayerRepo.EXPECT().Create(ctx, gomock.Any()).Return(nil)

	// When
	id, err := svc.Create(ctx, input)

	// Then
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if id == "" {
		t.Fatal("expected non-empty ID, got empty string")
	}
}

func TestPlayerService_Create_TeamNotFound(t *testing.T) {
	// Given
	svc, _, mockTeamRepo := setupPlayerService(t)
	ctx := context.Background()
	input := &domain.Player{
		TeamID:       "nonexistent",
		Name:         "Beckham",
		Height:       180.0,
		Weight:       75.0,
		Position:     domain.PositionCM,
		JerseyNumber: 10,
	}

	mockTeamRepo.EXPECT().FindByID(ctx, "nonexistent").Return(nil, domain.ErrTeamNotFound)

	// When
	id, err := svc.Create(ctx, input)

	// Then
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, domain.ErrTeamNotFound) {
		t.Fatalf("expected ErrTeamNotFound, got: %v", err)
	}
	if id != "" {
		t.Fatalf("expected empty ID on error, got %q", id)
	}
}

func TestPlayerService_Create_JerseyCheckRepoError(t *testing.T) {
	// Given
	svc, mockPlayerRepo, mockTeamRepo := setupPlayerService(t)
	ctx := context.Background()
	input := &domain.Player{
		TeamID:       "team-1",
		Name:         "Beckham",
		Height:       180.0,
		Weight:       75.0,
		Position:     domain.PositionCM,
		JerseyNumber: 10,
	}

	mockTeamRepo.EXPECT().FindByID(ctx, "team-1").Return(&domain.Team{ID: "team-1"}, nil)
	mockPlayerRepo.EXPECT().IsJerseyNumberTaken(ctx, "team-1", 10, "").Return(false, derrors.WrapErrorf(errors.New("db error"), derrors.ErrorCodeInternal, "failed to check jersey"))

	// When
	id, err := svc.Create(ctx, input)

	// Then
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	assertPlayerErrorCode(t, err, derrors.ErrorCodeInternal)
	if id != "" {
		t.Fatalf("expected empty ID on error, got %q", id)
	}
}

func TestPlayerService_Create_JerseyNumberTaken(t *testing.T) {
	// Given
	svc, mockPlayerRepo, mockTeamRepo := setupPlayerService(t)
	ctx := context.Background()
	input := &domain.Player{
		TeamID:       "team-1",
		Name:         "Beckham",
		Height:       180.0,
		Weight:       75.0,
		Position:     domain.PositionCM,
		JerseyNumber: 10,
	}

	mockTeamRepo.EXPECT().FindByID(ctx, "team-1").Return(&domain.Team{ID: "team-1"}, nil)
	mockPlayerRepo.EXPECT().IsJerseyNumberTaken(ctx, "team-1", 10, "").Return(true, nil)

	// When
	id, err := svc.Create(ctx, input)

	// Then
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, domain.ErrJerseyNumberTaken) {
		t.Fatalf("expected ErrJerseyNumberTaken, got: %v", err)
	}
	assertPlayerErrorCode(t, err, derrors.ErrorCodeDuplicate)
	if id != "" {
		t.Fatalf("expected empty ID on error, got %q", id)
	}
}

func TestPlayerService_Create_ValidationError_EmptyName(t *testing.T) {
	// Given
	svc, mockPlayerRepo, mockTeamRepo := setupPlayerService(t)
	ctx := context.Background()
	input := &domain.Player{
		TeamID:       "team-1",
		Name:         "",
		Height:       180.0,
		Weight:       75.0,
		Position:     domain.PositionCM,
		JerseyNumber: 10,
	}

	mockTeamRepo.EXPECT().FindByID(ctx, "team-1").Return(&domain.Team{ID: "team-1"}, nil)
	mockPlayerRepo.EXPECT().IsJerseyNumberTaken(ctx, "team-1", 10, "").Return(false, nil)

	// When
	id, err := svc.Create(ctx, input)

	// Then
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	assertPlayerErrorCode(t, err, derrors.ErrorCodeBadRequest)
	if id != "" {
		t.Fatalf("expected empty ID on error, got %q", id)
	}
}

func TestPlayerService_Create_ValidationError_InvalidPosition(t *testing.T) {
	// Given
	svc, mockPlayerRepo, mockTeamRepo := setupPlayerService(t)
	ctx := context.Background()
	input := &domain.Player{
		TeamID:       "team-1",
		Name:         "Beckham",
		Height:       180.0,
		Weight:       75.0,
		Position:     domain.Position(99),
		JerseyNumber: 10,
	}

	mockTeamRepo.EXPECT().FindByID(ctx, "team-1").Return(&domain.Team{ID: "team-1"}, nil)
	mockPlayerRepo.EXPECT().IsJerseyNumberTaken(ctx, "team-1", 10, "").Return(false, nil)

	// When
	id, err := svc.Create(ctx, input)

	// Then
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	assertPlayerErrorCode(t, err, derrors.ErrorCodeBadRequest)
	if id != "" {
		t.Fatalf("expected empty ID on error, got %q", id)
	}
}

func TestPlayerService_Create_ValidationError_JerseyNumberZero(t *testing.T) {
	// Given
	svc, mockPlayerRepo, mockTeamRepo := setupPlayerService(t)
	ctx := context.Background()
	input := &domain.Player{
		TeamID:       "team-1",
		Name:         "Beckham",
		Height:       180.0,
		Weight:       75.0,
		Position:     domain.PositionCM,
		JerseyNumber: 0,
	}

	mockTeamRepo.EXPECT().FindByID(ctx, "team-1").Return(&domain.Team{ID: "team-1"}, nil)
	mockPlayerRepo.EXPECT().IsJerseyNumberTaken(ctx, "team-1", 0, "").Return(false, nil)

	// When
	id, err := svc.Create(ctx, input)

	// Then
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	assertPlayerErrorCode(t, err, derrors.ErrorCodeBadRequest)
	if id != "" {
		t.Fatalf("expected empty ID on error, got %q", id)
	}
}

func TestPlayerService_Create_ValidationError_JerseyNumberExceeds99(t *testing.T) {
	// Given
	svc, mockPlayerRepo, mockTeamRepo := setupPlayerService(t)
	ctx := context.Background()
	input := &domain.Player{
		TeamID:       "team-1",
		Name:         "Beckham",
		Height:       180.0,
		Weight:       75.0,
		Position:     domain.PositionCM,
		JerseyNumber: 100,
	}

	mockTeamRepo.EXPECT().FindByID(ctx, "team-1").Return(&domain.Team{ID: "team-1"}, nil)
	mockPlayerRepo.EXPECT().IsJerseyNumberTaken(ctx, "team-1", 100, "").Return(false, nil)

	// When
	id, err := svc.Create(ctx, input)

	// Then
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	assertPlayerErrorCode(t, err, derrors.ErrorCodeBadRequest)
	if id != "" {
		t.Fatalf("expected empty ID on error, got %q", id)
	}
}

func TestPlayerService_Create_RepoError(t *testing.T) {
	// Given
	svc, mockPlayerRepo, mockTeamRepo := setupPlayerService(t)
	ctx := context.Background()
	input := &domain.Player{
		TeamID:       "team-1",
		Name:         "Beckham",
		Height:       180.0,
		Weight:       75.0,
		Position:     domain.PositionCM,
		JerseyNumber: 10,
	}

	mockTeamRepo.EXPECT().FindByID(ctx, "team-1").Return(&domain.Team{ID: "team-1"}, nil)
	mockPlayerRepo.EXPECT().IsJerseyNumberTaken(ctx, "team-1", 10, "").Return(false, nil)
	mockPlayerRepo.EXPECT().Create(ctx, gomock.Any()).Return(derrors.WrapErrorf(errors.New("db error"), derrors.ErrorCodeInternal, "failed to create player"))

	// When
	id, err := svc.Create(ctx, input)

	// Then
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	assertPlayerErrorCode(t, err, derrors.ErrorCodeInternal)
	if id != "" {
		t.Fatalf("expected empty ID on error, got %q", id)
	}
}

// ---------------------------------------------------------------------------
// GetByID
// ---------------------------------------------------------------------------

func TestPlayerService_GetByID_Success(t *testing.T) {
	// Given
	svc, mockPlayerRepo, _ := setupPlayerService(t)
	ctx := context.Background()
	expected := &domain.Player{ID: "player-1", TeamID: "team-1", Name: "Beckham"}

	mockPlayerRepo.EXPECT().FindByID(ctx, "player-1").Return(expected, nil)

	// When
	player, err := svc.GetByID(ctx, "player-1")

	// Then
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if player == nil {
		t.Fatal("expected player, got nil")
	}
	if player.ID != "player-1" {
		t.Fatalf("expected player.ID = %q, got %q", "player-1", player.ID)
	}
	if player.Name != "Beckham" {
		t.Fatalf("expected player.Name = %q, got %q", "Beckham", player.Name)
	}
}

func TestPlayerService_GetByID_NotFound(t *testing.T) {
	// Given
	svc, mockPlayerRepo, _ := setupPlayerService(t)
	ctx := context.Background()

	mockPlayerRepo.EXPECT().FindByID(ctx, "nonexistent").Return(nil, domain.ErrPlayerNotFound)

	// When
	player, err := svc.GetByID(ctx, "nonexistent")

	// Then
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, domain.ErrPlayerNotFound) {
		t.Fatalf("expected ErrPlayerNotFound, got: %v", err)
	}
	if player != nil {
		t.Fatalf("expected nil player on error, got %+v", player)
	}
}

// ---------------------------------------------------------------------------
// GetByTeamID
// ---------------------------------------------------------------------------

func TestPlayerService_GetByTeamID_Success(t *testing.T) {
	// Given
	svc, mockPlayerRepo, mockTeamRepo := setupPlayerService(t)
	ctx := context.Background()
	expected := []domain.Player{
		{ID: "player-1", TeamID: "team-1", Name: "Beckham"},
		{ID: "player-2", TeamID: "team-1", Name: "Zidane"},
	}

	mockTeamRepo.EXPECT().FindByID(ctx, "team-1").Return(&domain.Team{ID: "team-1"}, nil)
	mockPlayerRepo.EXPECT().FindByTeamID(ctx, "team-1").Return(expected, nil)

	// When
	players, err := svc.GetByTeamID(ctx, "team-1")

	// Then
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if len(players) != 2 {
		t.Fatalf("expected 2 players, got %d", len(players))
	}
	if players[0].Name != "Beckham" {
		t.Fatalf("expected first player = %q, got %q", "Beckham", players[0].Name)
	}
}

func TestPlayerService_GetByTeamID_TeamNotFound(t *testing.T) {
	// Given
	svc, _, mockTeamRepo := setupPlayerService(t)
	ctx := context.Background()

	mockTeamRepo.EXPECT().FindByID(ctx, "nonexistent").Return(nil, domain.ErrTeamNotFound)

	// When
	players, err := svc.GetByTeamID(ctx, "nonexistent")

	// Then
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, domain.ErrTeamNotFound) {
		t.Fatalf("expected ErrTeamNotFound, got: %v", err)
	}
	if players != nil {
		t.Fatalf("expected nil players on error, got %d items", len(players))
	}
}

func TestPlayerService_GetByTeamID_RepoError(t *testing.T) {
	// Given
	svc, mockPlayerRepo, mockTeamRepo := setupPlayerService(t)
	ctx := context.Background()

	mockTeamRepo.EXPECT().FindByID(ctx, "team-1").Return(&domain.Team{ID: "team-1"}, nil)
	mockPlayerRepo.EXPECT().FindByTeamID(ctx, "team-1").Return(nil, derrors.WrapErrorf(errors.New("db error"), derrors.ErrorCodeInternal, "failed to fetch players"))

	// When
	players, err := svc.GetByTeamID(ctx, "team-1")

	// Then
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	assertPlayerErrorCode(t, err, derrors.ErrorCodeInternal)
	if players != nil {
		t.Fatalf("expected nil players on error, got %d items", len(players))
	}
}

// ---------------------------------------------------------------------------
// Update
// ---------------------------------------------------------------------------

func TestPlayerService_Update_Success(t *testing.T) {
	// Given
	svc, mockPlayerRepo, _ := setupPlayerService(t)
	ctx := context.Background()
	existing := &domain.Player{
		ID:           "player-1",
		TeamID:       "team-1",
		Name:         "Old Name",
		Height:       175.0,
		Weight:       70.0,
		Position:     domain.PositionST,
		JerseyNumber: 9,
	}
	update := &domain.Player{
		Name:         "New Name",
		Height:       180.0,
		Weight:       75.0,
		Position:     domain.PositionCM,
		JerseyNumber: 10,
	}

	mockPlayerRepo.EXPECT().FindByID(ctx, "player-1").Return(existing, nil)
	mockPlayerRepo.EXPECT().IsJerseyNumberTaken(ctx, "team-1", 10, "player-1").Return(false, nil)
	mockPlayerRepo.EXPECT().Update(ctx, gomock.Any()).Return(nil)

	// When
	err := svc.Update(ctx, "player-1", update)

	// Then
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestPlayerService_Update_NotFound(t *testing.T) {
	// Given
	svc, mockPlayerRepo, _ := setupPlayerService(t)
	ctx := context.Background()
	update := &domain.Player{
		Name:         "New Name",
		Height:       180.0,
		Weight:       75.0,
		Position:     domain.PositionCM,
		JerseyNumber: 10,
	}

	mockPlayerRepo.EXPECT().FindByID(ctx, "nonexistent").Return(nil, domain.ErrPlayerNotFound)

	// When
	err := svc.Update(ctx, "nonexistent", update)

	// Then
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, domain.ErrPlayerNotFound) {
		t.Fatalf("expected ErrPlayerNotFound, got: %v", err)
	}
}

func TestPlayerService_Update_JerseyCheckRepoError(t *testing.T) {
	// Given
	svc, mockPlayerRepo, _ := setupPlayerService(t)
	ctx := context.Background()
	existing := &domain.Player{
		ID:           "player-1",
		TeamID:       "team-1",
		Name:         "Old Name",
		JerseyNumber: 9,
	}
	update := &domain.Player{
		Name:         "New Name",
		Height:       180.0,
		Weight:       75.0,
		Position:     domain.PositionCM,
		JerseyNumber: 10,
	}

	mockPlayerRepo.EXPECT().FindByID(ctx, "player-1").Return(existing, nil)
	mockPlayerRepo.EXPECT().IsJerseyNumberTaken(ctx, "team-1", 10, "player-1").Return(false, derrors.WrapErrorf(errors.New("db error"), derrors.ErrorCodeInternal, "failed to check jersey"))

	// When
	err := svc.Update(ctx, "player-1", update)

	// Then
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	assertPlayerErrorCode(t, err, derrors.ErrorCodeInternal)
}

func TestPlayerService_Update_JerseyNumberTaken(t *testing.T) {
	// Given
	svc, mockPlayerRepo, _ := setupPlayerService(t)
	ctx := context.Background()
	existing := &domain.Player{
		ID:           "player-1",
		TeamID:       "team-1",
		Name:         "Old Name",
		JerseyNumber: 9,
	}
	update := &domain.Player{
		Name:         "New Name",
		Height:       180.0,
		Weight:       75.0,
		Position:     domain.PositionCM,
		JerseyNumber: 10,
	}

	mockPlayerRepo.EXPECT().FindByID(ctx, "player-1").Return(existing, nil)
	mockPlayerRepo.EXPECT().IsJerseyNumberTaken(ctx, "team-1", 10, "player-1").Return(true, nil)

	// When
	err := svc.Update(ctx, "player-1", update)

	// Then
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, domain.ErrJerseyNumberTaken) {
		t.Fatalf("expected ErrJerseyNumberTaken, got: %v", err)
	}
	assertPlayerErrorCode(t, err, derrors.ErrorCodeDuplicate)
}

func TestPlayerService_Update_ValidationError_EmptyName(t *testing.T) {
	// Given
	svc, mockPlayerRepo, _ := setupPlayerService(t)
	ctx := context.Background()
	existing := &domain.Player{
		ID:           "player-1",
		TeamID:       "team-1",
		Name:         "Old Name",
		JerseyNumber: 9,
	}
	update := &domain.Player{
		Name:         "",
		Height:       180.0,
		Weight:       75.0,
		Position:     domain.PositionCM,
		JerseyNumber: 10,
	}

	mockPlayerRepo.EXPECT().FindByID(ctx, "player-1").Return(existing, nil)
	mockPlayerRepo.EXPECT().IsJerseyNumberTaken(ctx, "team-1", 10, "player-1").Return(false, nil)

	// When
	err := svc.Update(ctx, "player-1", update)

	// Then
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	assertPlayerErrorCode(t, err, derrors.ErrorCodeBadRequest)
}

func TestPlayerService_Update_RepoError(t *testing.T) {
	// Given
	svc, mockPlayerRepo, _ := setupPlayerService(t)
	ctx := context.Background()
	existing := &domain.Player{
		ID:           "player-1",
		TeamID:       "team-1",
		Name:         "Old Name",
		Height:       175.0,
		Weight:       70.0,
		Position:     domain.PositionST,
		JerseyNumber: 9,
	}
	update := &domain.Player{
		Name:         "New Name",
		Height:       180.0,
		Weight:       75.0,
		Position:     domain.PositionCM,
		JerseyNumber: 10,
	}

	mockPlayerRepo.EXPECT().FindByID(ctx, "player-1").Return(existing, nil)
	mockPlayerRepo.EXPECT().IsJerseyNumberTaken(ctx, "team-1", 10, "player-1").Return(false, nil)
	mockPlayerRepo.EXPECT().Update(ctx, gomock.Any()).Return(derrors.WrapErrorf(errors.New("db error"), derrors.ErrorCodeInternal, "failed to update player"))

	// When
	err := svc.Update(ctx, "player-1", update)

	// Then
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	assertPlayerErrorCode(t, err, derrors.ErrorCodeInternal)
}

// ---------------------------------------------------------------------------
// Delete
// ---------------------------------------------------------------------------

func TestPlayerService_Delete_Success(t *testing.T) {
	// Given
	svc, mockPlayerRepo, _ := setupPlayerService(t)
	ctx := context.Background()

	mockPlayerRepo.EXPECT().FindByID(ctx, "player-1").Return(&domain.Player{ID: "player-1"}, nil)
	mockPlayerRepo.EXPECT().SoftDelete(ctx, "player-1").Return(nil)

	// When
	err := svc.Delete(ctx, "player-1")

	// Then
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestPlayerService_Delete_NotFound(t *testing.T) {
	// Given
	svc, mockPlayerRepo, _ := setupPlayerService(t)
	ctx := context.Background()

	mockPlayerRepo.EXPECT().FindByID(ctx, "nonexistent").Return(nil, domain.ErrPlayerNotFound)

	// When
	err := svc.Delete(ctx, "nonexistent")

	// Then
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, domain.ErrPlayerNotFound) {
		t.Fatalf("expected ErrPlayerNotFound, got: %v", err)
	}
}

func TestPlayerService_Delete_RepoError(t *testing.T) {
	// Given
	svc, mockPlayerRepo, _ := setupPlayerService(t)
	ctx := context.Background()

	mockPlayerRepo.EXPECT().FindByID(ctx, "player-1").Return(&domain.Player{ID: "player-1"}, nil)
	mockPlayerRepo.EXPECT().SoftDelete(ctx, "player-1").Return(derrors.WrapErrorf(errors.New("db error"), derrors.ErrorCodeInternal, "failed to delete player"))

	// When
	err := svc.Delete(ctx, "player-1")

	// Then
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	assertPlayerErrorCode(t, err, derrors.ErrorCodeInternal)
}
