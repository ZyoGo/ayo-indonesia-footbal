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

func setupTeamService(t *testing.T) (*TeamService, *mockDomain.MockTeamRepository) {
	t.Helper()
	ctrl := gomock.NewController(t)
	mockRepo := mockDomain.NewMockTeamRepository(ctrl)
	svc := &TeamService{teamRepo: mockRepo}
	return svc, mockRepo
}

// assertErrorCode verifies the error is a *derrors.Error with the expected code.
func assertErrorCode(t *testing.T, err error, expectedCode derrors.ErrorCode) {
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

func TestTeamService_Create_Success(t *testing.T) {
	// Given
	svc, mockRepo := setupTeamService(t)
	ctx := context.Background()
	input := &domain.Team{
		Name:        "Persib Bandung",
		LogoURL:     "https://example.com/logo.png",
		YearFounded: 1933,
		Address:     "Jl. Ahmad Yani",
		City:        "Bandung",
	}

	mockRepo.EXPECT().Create(ctx, gomock.Any()).Return(nil)

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

func TestTeamService_Create_ValidationError_EmptyName(t *testing.T) {
	// Given
	svc, _ := setupTeamService(t)
	ctx := context.Background()
	input := &domain.Team{
		Name:        "",
		YearFounded: 1933,
		City:        "Bandung",
	}

	// When
	id, err := svc.Create(ctx, input)

	// Then
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	assertErrorCode(t, err, derrors.ErrorCodeBadRequest)
	if id != "" {
		t.Fatalf("expected empty ID on error, got %q", id)
	}
}

func TestTeamService_Create_ValidationError_EmptyCity(t *testing.T) {
	// Given
	svc, _ := setupTeamService(t)
	ctx := context.Background()
	input := &domain.Team{
		Name:        "Persib",
		YearFounded: 1933,
		City:        "",
	}

	// When
	id, err := svc.Create(ctx, input)

	// Then
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	assertErrorCode(t, err, derrors.ErrorCodeBadRequest)
	if id != "" {
		t.Fatalf("expected empty ID on error, got %q", id)
	}
}

func TestTeamService_Create_ValidationError_InvalidYear(t *testing.T) {
	// Given
	svc, _ := setupTeamService(t)
	ctx := context.Background()
	input := &domain.Team{
		Name:        "Persib",
		YearFounded: 0,
		City:        "Bandung",
	}

	// When
	id, err := svc.Create(ctx, input)

	// Then
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	assertErrorCode(t, err, derrors.ErrorCodeBadRequest)
	if id != "" {
		t.Fatalf("expected empty ID on error, got %q", id)
	}
}

func TestTeamService_Create_RepoError(t *testing.T) {
	// Given
	svc, mockRepo := setupTeamService(t)
	ctx := context.Background()
	input := &domain.Team{
		Name:        "Persib",
		LogoURL:     "https://example.com/logo.png",
		YearFounded: 1933,
		Address:     "Jl. Ahmad Yani",
		City:        "Bandung",
	}

	mockRepo.EXPECT().Create(ctx, gomock.Any()).Return(derrors.WrapErrorf(errors.New("db error"), derrors.ErrorCodeInternal, "failed to create team"))

	// When
	id, err := svc.Create(ctx, input)

	// Then
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	assertErrorCode(t, err, derrors.ErrorCodeInternal)
	if id != "" {
		t.Fatalf("expected empty ID on error, got %q", id)
	}
}

// ---------------------------------------------------------------------------
// GetByID
// ---------------------------------------------------------------------------

func TestTeamService_GetByID_Success(t *testing.T) {
	// Given
	svc, mockRepo := setupTeamService(t)
	ctx := context.Background()
	expected := &domain.Team{ID: "team-1", Name: "Persib", City: "Bandung"}

	mockRepo.EXPECT().FindByID(ctx, "team-1").Return(expected, nil)

	// When
	team, err := svc.GetByID(ctx, "team-1")

	// Then
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if team == nil {
		t.Fatal("expected team, got nil")
	}
	if team.ID != "team-1" {
		t.Fatalf("expected team.ID = %q, got %q", "team-1", team.ID)
	}
	if team.Name != "Persib" {
		t.Fatalf("expected team.Name = %q, got %q", "Persib", team.Name)
	}
}

func TestTeamService_GetByID_NotFound(t *testing.T) {
	// Given
	svc, mockRepo := setupTeamService(t)
	ctx := context.Background()

	mockRepo.EXPECT().FindByID(ctx, "nonexistent").Return(nil, domain.ErrTeamNotFound)

	// When
	team, err := svc.GetByID(ctx, "nonexistent")

	// Then
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, domain.ErrTeamNotFound) {
		t.Fatalf("expected ErrTeamNotFound, got: %v", err)
	}
	if team != nil {
		t.Fatalf("expected nil team on error, got %+v", team)
	}
}

// ---------------------------------------------------------------------------
// GetAll
// ---------------------------------------------------------------------------

func TestTeamService_GetAll_Success(t *testing.T) {
	// Given
	svc, mockRepo := setupTeamService(t)
	ctx := context.Background()
	expected := []domain.Team{
		{ID: "team-1", Name: "Persib"},
		{ID: "team-2", Name: "Persija"},
	}

	mockRepo.EXPECT().FindAll(ctx).Return(expected, nil)

	// When
	teams, err := svc.GetAll(ctx)

	// Then
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if len(teams) != 2 {
		t.Fatalf("expected 2 teams, got %d", len(teams))
	}
	if teams[0].ID != "team-1" {
		t.Fatalf("expected first team ID = %q, got %q", "team-1", teams[0].ID)
	}
}

func TestTeamService_GetAll_EmptyList(t *testing.T) {
	// Given
	svc, mockRepo := setupTeamService(t)
	ctx := context.Background()

	mockRepo.EXPECT().FindAll(ctx).Return([]domain.Team{}, nil)

	// When
	teams, err := svc.GetAll(ctx)

	// Then
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if len(teams) != 0 {
		t.Fatalf("expected 0 teams, got %d", len(teams))
	}
}

func TestTeamService_GetAll_RepoError(t *testing.T) {
	// Given
	svc, mockRepo := setupTeamService(t)
	ctx := context.Background()

	mockRepo.EXPECT().FindAll(ctx).Return(nil, derrors.WrapErrorf(errors.New("db error"), derrors.ErrorCodeInternal, "failed to fetch teams"))

	// When
	teams, err := svc.GetAll(ctx)

	// Then
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	assertErrorCode(t, err, derrors.ErrorCodeInternal)
	if teams != nil {
		t.Fatalf("expected nil teams on error, got %d items", len(teams))
	}
}

// ---------------------------------------------------------------------------
// Update
// ---------------------------------------------------------------------------

func TestTeamService_Update_Success(t *testing.T) {
	// Given
	svc, mockRepo := setupTeamService(t)
	ctx := context.Background()
	existing := &domain.Team{ID: "team-1", Name: "Old Name", YearFounded: 1933, City: "Bandung"}
	update := &domain.Team{
		Name:        "New Name",
		LogoURL:     "https://example.com/new-logo.png",
		YearFounded: 1933,
		Address:     "Jl. Baru",
		City:        "Bandung",
	}

	mockRepo.EXPECT().FindByID(ctx, "team-1").Return(existing, nil)
	mockRepo.EXPECT().Update(ctx, gomock.Any()).Return(nil)

	// When
	err := svc.Update(ctx, "team-1", update)

	// Then
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestTeamService_Update_NotFound(t *testing.T) {
	// Given
	svc, mockRepo := setupTeamService(t)
	ctx := context.Background()
	update := &domain.Team{Name: "New Name", YearFounded: 1933, City: "Bandung"}

	mockRepo.EXPECT().FindByID(ctx, "nonexistent").Return(nil, domain.ErrTeamNotFound)

	// When
	err := svc.Update(ctx, "nonexistent", update)

	// Then
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, domain.ErrTeamNotFound) {
		t.Fatalf("expected ErrTeamNotFound, got: %v", err)
	}
}

func TestTeamService_Update_ValidationError_EmptyName(t *testing.T) {
	// Given
	svc, mockRepo := setupTeamService(t)
	ctx := context.Background()
	existing := &domain.Team{ID: "team-1", Name: "Old Name", YearFounded: 1933, City: "Bandung"}
	update := &domain.Team{Name: "", YearFounded: 1933, City: "Bandung"}

	mockRepo.EXPECT().FindByID(ctx, "team-1").Return(existing, nil)

	// When
	err := svc.Update(ctx, "team-1", update)

	// Then
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	assertErrorCode(t, err, derrors.ErrorCodeBadRequest)
}

func TestTeamService_Update_RepoError(t *testing.T) {
	// Given
	svc, mockRepo := setupTeamService(t)
	ctx := context.Background()
	existing := &domain.Team{ID: "team-1", Name: "Old Name", YearFounded: 1933, City: "Bandung"}
	update := &domain.Team{
		Name:        "New Name",
		LogoURL:     "https://example.com/logo.png",
		YearFounded: 1933,
		Address:     "Jl. Baru",
		City:        "Bandung",
	}

	mockRepo.EXPECT().FindByID(ctx, "team-1").Return(existing, nil)
	mockRepo.EXPECT().Update(ctx, gomock.Any()).Return(derrors.WrapErrorf(errors.New("db error"), derrors.ErrorCodeInternal, "failed to update team"))

	// When
	err := svc.Update(ctx, "team-1", update)

	// Then
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	assertErrorCode(t, err, derrors.ErrorCodeInternal)
}

// ---------------------------------------------------------------------------
// Delete
// ---------------------------------------------------------------------------

func TestTeamService_Delete_Success(t *testing.T) {
	// Given
	svc, mockRepo := setupTeamService(t)
	ctx := context.Background()

	mockRepo.EXPECT().FindByID(ctx, "team-1").Return(&domain.Team{ID: "team-1"}, nil)
	mockRepo.EXPECT().SoftDelete(ctx, "team-1").Return(nil)

	// When
	err := svc.Delete(ctx, "team-1")

	// Then
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestTeamService_Delete_NotFound(t *testing.T) {
	// Given
	svc, mockRepo := setupTeamService(t)
	ctx := context.Background()

	mockRepo.EXPECT().FindByID(ctx, "nonexistent").Return(nil, domain.ErrTeamNotFound)

	// When
	err := svc.Delete(ctx, "nonexistent")

	// Then
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, domain.ErrTeamNotFound) {
		t.Fatalf("expected ErrTeamNotFound, got: %v", err)
	}
}

func TestTeamService_Delete_RepoError(t *testing.T) {
	// Given
	svc, mockRepo := setupTeamService(t)
	ctx := context.Background()

	mockRepo.EXPECT().FindByID(ctx, "team-1").Return(&domain.Team{ID: "team-1"}, nil)
	mockRepo.EXPECT().SoftDelete(ctx, "team-1").Return(derrors.WrapErrorf(errors.New("db error"), derrors.ErrorCodeInternal, "failed to delete team"))

	// When
	err := svc.Delete(ctx, "team-1")

	// Then
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	assertErrorCode(t, err, derrors.ErrorCodeInternal)
}
