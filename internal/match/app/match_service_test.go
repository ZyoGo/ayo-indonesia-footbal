package app

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/ZyoGo/ayo-indonesia-footbal/internal/match/domain"
	mockDomain "github.com/ZyoGo/ayo-indonesia-footbal/internal/match/mock"
	"github.com/ZyoGo/ayo-indonesia-footbal/pkg/derrors"
	"go.uber.org/mock/gomock"
)

func setupMatchService(t *testing.T) (
	*MatchService,
	*mockDomain.MockMatchRepository,
	*mockDomain.MockMatchResultRepository,
	*mockDomain.MockReportRepository,
) {
	t.Helper()
	ctrl := gomock.NewController(t)
	mockMatchRepo := mockDomain.NewMockMatchRepository(ctrl)
	mockResultRepo := mockDomain.NewMockMatchResultRepository(ctrl)
	mockReportRepo := mockDomain.NewMockReportRepository(ctrl)
	svc := &MatchService{
		matchRepo:  mockMatchRepo,
		resultRepo: mockResultRepo,
		reportRepo: mockReportRepo,
	}
	return svc, mockMatchRepo, mockResultRepo, mockReportRepo
}

func assertMatchErrorCode(t *testing.T, err error, expectedCode derrors.ErrorCode) {
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
// CreateMatch
// ---------------------------------------------------------------------------

func TestMatchService_CreateMatch_Success(t *testing.T) {
	svc, mockMatchRepo, _, _ := setupMatchService(t)
	ctx := context.Background()
	input := &domain.Match{
		HomeTeamID: "team-1",
		AwayTeamID: "team-2",
		MatchDate:  time.Date(2026, 3, 15, 0, 0, 0, 0, time.UTC),
		MatchTime:  "19:30",
		Stadium:    "Gelora Bung Karno",
	}

	mockMatchRepo.EXPECT().Create(ctx, gomock.Any()).Return(nil)

	id, err := svc.CreateMatch(ctx, input)

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if id == "" {
		t.Fatal("expected non-empty ID, got empty string")
	}
}

func TestMatchService_CreateMatch_SameTeamError(t *testing.T) {
	svc, _, _, _ := setupMatchService(t)
	ctx := context.Background()
	input := &domain.Match{
		HomeTeamID: "team-1",
		AwayTeamID: "team-1",
		MatchDate:  time.Date(2026, 3, 15, 0, 0, 0, 0, time.UTC),
		MatchTime:  "19:30",
		Stadium:    "Gelora Bung Karno",
	}

	id, err := svc.CreateMatch(ctx, input)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	assertMatchErrorCode(t, err, derrors.ErrorCodeBadRequest)
	if id != "" {
		t.Fatalf("expected empty ID on error, got %q", id)
	}
}

func TestMatchService_CreateMatch_InvalidTime(t *testing.T) {
	svc, _, _, _ := setupMatchService(t)
	ctx := context.Background()
	input := &domain.Match{
		HomeTeamID: "team-1",
		AwayTeamID: "team-2",
		MatchDate:  time.Date(2026, 3, 15, 0, 0, 0, 0, time.UTC),
		MatchTime:  "25:00",
		Stadium:    "Gelora Bung Karno",
	}

	id, err := svc.CreateMatch(ctx, input)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	assertMatchErrorCode(t, err, derrors.ErrorCodeBadRequest)
	if id != "" {
		t.Fatalf("expected empty ID on error, got %q", id)
	}
}

func TestMatchService_CreateMatch_EmptyHomeTeam(t *testing.T) {
	svc, _, _, _ := setupMatchService(t)
	ctx := context.Background()
	input := &domain.Match{
		HomeTeamID: "",
		AwayTeamID: "team-2",
		MatchDate:  time.Date(2026, 3, 15, 0, 0, 0, 0, time.UTC),
		MatchTime:  "19:30",
	}

	id, err := svc.CreateMatch(ctx, input)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	assertMatchErrorCode(t, err, derrors.ErrorCodeBadRequest)
	if id != "" {
		t.Fatalf("expected empty ID on error, got %q", id)
	}
}

func TestMatchService_CreateMatch_PastDate(t *testing.T) {
	svc, _, _, _ := setupMatchService(t)
	ctx := context.Background()
	input := &domain.Match{
		HomeTeamID: "team-1",
		AwayTeamID: "team-2",
		MatchDate:  time.Now().AddDate(0, 0, -1), // yesterday
		MatchTime:  "19:00",
	}

	id, err := svc.CreateMatch(ctx, input)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	assertMatchErrorCode(t, err, derrors.ErrorCodeBadRequest)
	if id != "" {
		t.Fatalf("expected empty ID on error, got %q", id)
	}
}

func TestMatchService_CreateMatch_RepoError(t *testing.T) {
	svc, mockMatchRepo, _, _ := setupMatchService(t)
	ctx := context.Background()
	input := &domain.Match{
		HomeTeamID: "team-1",
		AwayTeamID: "team-2",
		MatchDate:  time.Date(2026, 3, 15, 0, 0, 0, 0, time.UTC),
		MatchTime:  "19:30",
		Stadium:    "Gelora Bung Karno",
	}

	mockMatchRepo.EXPECT().Create(ctx, gomock.Any()).Return(derrors.WrapErrorf(errors.New("db error"), derrors.ErrorCodeInternal, "failed to create match"))

	id, err := svc.CreateMatch(ctx, input)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	assertMatchErrorCode(t, err, derrors.ErrorCodeInternal)
	if id != "" {
		t.Fatalf("expected empty ID on error, got %q", id)
	}
}

// ---------------------------------------------------------------------------
// GetMatchByID
// ---------------------------------------------------------------------------

func TestMatchService_GetMatchByID_Success(t *testing.T) {
	svc, mockMatchRepo, _, _ := setupMatchService(t)
	ctx := context.Background()
	expected := &domain.Match{ID: "match-1", HomeTeamID: "team-1", AwayTeamID: "team-2"}

	mockMatchRepo.EXPECT().FindByID(ctx, "match-1").Return(expected, nil)

	match, err := svc.GetMatchByID(ctx, "match-1")

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if match == nil {
		t.Fatal("expected match, got nil")
	}
	if match.ID != "match-1" {
		t.Fatalf("expected match.ID = %q, got %q", "match-1", match.ID)
	}
}

func TestMatchService_GetMatchByID_NotFound(t *testing.T) {
	svc, mockMatchRepo, _, _ := setupMatchService(t)
	ctx := context.Background()

	mockMatchRepo.EXPECT().FindByID(ctx, "nonexistent").Return(nil, domain.ErrMatchNotFound)

	match, err := svc.GetMatchByID(ctx, "nonexistent")

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, domain.ErrMatchNotFound) {
		t.Fatalf("expected ErrMatchNotFound, got: %v", err)
	}
	if match != nil {
		t.Fatalf("expected nil match on error, got %+v", match)
	}
}

// ---------------------------------------------------------------------------
// GetAllMatches
// ---------------------------------------------------------------------------

func TestMatchService_GetAllMatches_Success(t *testing.T) {
	svc, mockMatchRepo, _, _ := setupMatchService(t)
	ctx := context.Background()
	expected := []domain.Match{
		{ID: "match-1", HomeTeamID: "team-1", AwayTeamID: "team-2"},
		{ID: "match-2", HomeTeamID: "team-3", AwayTeamID: "team-4"},
	}

	mockMatchRepo.EXPECT().FindAll(ctx).Return(expected, nil)

	matches, err := svc.GetAllMatches(ctx)

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if len(matches) != 2 {
		t.Fatalf("expected 2 matches, got %d", len(matches))
	}
}

func TestMatchService_GetAllMatches_RepoError(t *testing.T) {
	svc, mockMatchRepo, _, _ := setupMatchService(t)
	ctx := context.Background()

	mockMatchRepo.EXPECT().FindAll(ctx).Return(nil, derrors.WrapErrorf(errors.New("db error"), derrors.ErrorCodeInternal, "failed to fetch matches"))

	matches, err := svc.GetAllMatches(ctx)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	assertMatchErrorCode(t, err, derrors.ErrorCodeInternal)
	if matches != nil {
		t.Fatalf("expected nil matches on error, got %d items", len(matches))
	}
}

// ---------------------------------------------------------------------------
// ReportResult
// ---------------------------------------------------------------------------

func TestMatchService_ReportResult_Success(t *testing.T) {
	svc, mockMatchRepo, mockResultRepo, _ := setupMatchService(t)
	ctx := context.Background()
	matchID := "match-1"
	result := &domain.MatchResult{
		HomeScore: 2,
		AwayScore: 1,
		Goals: []domain.Goal{
			{PlayerID: "player-1", TeamID: "team-1", GoalMinute: 15},
			{PlayerID: "player-2", TeamID: "team-1", GoalMinute: 45},
			{PlayerID: "player-3", TeamID: "team-2", GoalMinute: 70},
		},
	}

	mockMatchRepo.EXPECT().FindByID(ctx, matchID).Return(&domain.Match{ID: matchID, HomeTeamID: "team-1", AwayTeamID: "team-2", Stadium: "Gelora Bung Karno"}, nil)
	mockResultRepo.EXPECT().ExistsByMatchID(ctx, matchID).Return(false, nil)
	mockResultRepo.EXPECT().Create(ctx, gomock.Any()).Return(nil)

	id, err := svc.ReportResult(ctx, matchID, result)

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if id == "" {
		t.Fatal("expected non-empty result ID, got empty string")
	}
}

func TestMatchService_ReportResult_MatchNotFound(t *testing.T) {
	svc, mockMatchRepo, _, _ := setupMatchService(t)
	ctx := context.Background()
	matchID := "nonexistent"
	result := &domain.MatchResult{HomeScore: 1, AwayScore: 0, Goals: []domain.Goal{
		{PlayerID: "player-1", TeamID: "team-1", GoalMinute: 10},
	}}

	mockMatchRepo.EXPECT().FindByID(ctx, matchID).Return(nil, domain.ErrMatchNotFound)

	id, err := svc.ReportResult(ctx, matchID, result)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, domain.ErrMatchNotFound) {
		t.Fatalf("expected ErrMatchNotFound, got: %v", err)
	}
	if id != "" {
		t.Fatalf("expected empty ID on error, got %q", id)
	}
}

func TestMatchService_ReportResult_AlreadyExists(t *testing.T) {
	svc, mockMatchRepo, mockResultRepo, _ := setupMatchService(t)
	ctx := context.Background()
	matchID := "match-1"
	result := &domain.MatchResult{HomeScore: 1, AwayScore: 0, Goals: []domain.Goal{
		{PlayerID: "player-1", TeamID: "team-1", GoalMinute: 10},
	}}

	mockMatchRepo.EXPECT().FindByID(ctx, matchID).Return(&domain.Match{ID: matchID, HomeTeamID: "team-1", AwayTeamID: "team-2", Stadium: "Gelora Bung Karno"}, nil)
	mockResultRepo.EXPECT().ExistsByMatchID(ctx, matchID).Return(true, nil)

	id, err := svc.ReportResult(ctx, matchID, result)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, domain.ErrResultAlreadyExists) {
		t.Fatalf("expected ErrResultAlreadyExists, got: %v", err)
	}
	assertMatchErrorCode(t, err, derrors.ErrorCodeDuplicate)
	if id != "" {
		t.Fatalf("expected empty ID on error, got %q", id)
	}
}

func TestMatchService_ReportResult_ExistsCheckRepoError(t *testing.T) {
	svc, mockMatchRepo, mockResultRepo, _ := setupMatchService(t)
	ctx := context.Background()
	matchID := "match-1"
	result := &domain.MatchResult{HomeScore: 1, AwayScore: 0, Goals: []domain.Goal{
		{PlayerID: "player-1", TeamID: "team-1", GoalMinute: 10},
	}}

	mockMatchRepo.EXPECT().FindByID(ctx, matchID).Return(&domain.Match{ID: matchID, HomeTeamID: "team-1", AwayTeamID: "team-2", Stadium: "Gelora Bung Karno"}, nil)
	mockResultRepo.EXPECT().ExistsByMatchID(ctx, matchID).Return(false, derrors.WrapErrorf(errors.New("db error"), derrors.ErrorCodeInternal, "failed to check result"))

	id, err := svc.ReportResult(ctx, matchID, result)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	assertMatchErrorCode(t, err, derrors.ErrorCodeInternal)
	if id != "" {
		t.Fatalf("expected empty ID on error, got %q", id)
	}
}

func TestMatchService_ReportResult_ScoreMismatch(t *testing.T) {
	svc, mockMatchRepo, mockResultRepo, _ := setupMatchService(t)
	ctx := context.Background()
	matchID := "match-1"
	// 2-1 = 3 goals expected, but only 2 goals provided
	result := &domain.MatchResult{
		HomeScore: 2,
		AwayScore: 1,
		Goals: []domain.Goal{
			{PlayerID: "player-1", TeamID: "team-1", GoalMinute: 15},
			{PlayerID: "player-2", TeamID: "team-1", GoalMinute: 45},
		},
	}

	mockMatchRepo.EXPECT().FindByID(ctx, matchID).Return(&domain.Match{ID: matchID, HomeTeamID: "team-1", AwayTeamID: "team-2", Stadium: "Gelora Bung Karno"}, nil)
	mockResultRepo.EXPECT().ExistsByMatchID(ctx, matchID).Return(false, nil)

	id, err := svc.ReportResult(ctx, matchID, result)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	assertMatchErrorCode(t, err, derrors.ErrorCodeBadRequest)
	if id != "" {
		t.Fatalf("expected empty ID on error, got %q", id)
	}
}

func TestMatchService_ReportResult_NegativeScore(t *testing.T) {
	svc, mockMatchRepo, mockResultRepo, _ := setupMatchService(t)
	ctx := context.Background()
	matchID := "match-1"
	result := &domain.MatchResult{
		HomeScore: -1,
		AwayScore: 0,
		Goals:     []domain.Goal{},
	}

	mockMatchRepo.EXPECT().FindByID(ctx, matchID).Return(&domain.Match{ID: matchID, HomeTeamID: "team-1", AwayTeamID: "team-2", Stadium: "Gelora Bung Karno"}, nil)
	mockResultRepo.EXPECT().ExistsByMatchID(ctx, matchID).Return(false, nil)

	id, err := svc.ReportResult(ctx, matchID, result)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	assertMatchErrorCode(t, err, derrors.ErrorCodeBadRequest)
	if id != "" {
		t.Fatalf("expected empty ID on error, got %q", id)
	}
}

func TestMatchService_ReportResult_InvalidGoalMinute(t *testing.T) {
	svc, mockMatchRepo, mockResultRepo, _ := setupMatchService(t)
	ctx := context.Background()
	matchID := "match-1"
	result := &domain.MatchResult{
		HomeScore: 1,
		AwayScore: 0,
		Goals: []domain.Goal{
			{PlayerID: "player-1", TeamID: "team-1", GoalMinute: 0}, // invalid
		},
	}

	mockMatchRepo.EXPECT().FindByID(ctx, matchID).Return(&domain.Match{ID: matchID, HomeTeamID: "team-1", AwayTeamID: "team-2", Stadium: "Gelora Bung Karno"}, nil)
	mockResultRepo.EXPECT().ExistsByMatchID(ctx, matchID).Return(false, nil)

	id, err := svc.ReportResult(ctx, matchID, result)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	assertMatchErrorCode(t, err, derrors.ErrorCodeBadRequest)
	if id != "" {
		t.Fatalf("expected empty ID on error, got %q", id)
	}
}

func TestMatchService_ReportResult_GoalMinuteExceedsMax(t *testing.T) {
	svc, mockMatchRepo, mockResultRepo, _ := setupMatchService(t)
	ctx := context.Background()
	matchID := "match-1"
	result := &domain.MatchResult{
		HomeScore: 1,
		AwayScore: 0,
		Goals: []domain.Goal{
			{PlayerID: "player-1", TeamID: "team-1", GoalMinute: 200}, // exceeds 150
		},
	}

	mockMatchRepo.EXPECT().FindByID(ctx, matchID).Return(&domain.Match{ID: matchID, HomeTeamID: "team-1", AwayTeamID: "team-2", Stadium: "Gelora Bung Karno"}, nil)
	mockResultRepo.EXPECT().ExistsByMatchID(ctx, matchID).Return(false, nil)

	id, err := svc.ReportResult(ctx, matchID, result)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	assertMatchErrorCode(t, err, derrors.ErrorCodeBadRequest)
	if id != "" {
		t.Fatalf("expected empty ID on error, got %q", id)
	}
}

func TestMatchService_ReportResult_MissingPlayerID(t *testing.T) {
	svc, mockMatchRepo, mockResultRepo, _ := setupMatchService(t)
	ctx := context.Background()
	matchID := "match-1"
	result := &domain.MatchResult{
		HomeScore: 1,
		AwayScore: 0,
		Goals: []domain.Goal{
			{PlayerID: "", TeamID: "team-1", GoalMinute: 10},
		},
	}

	mockMatchRepo.EXPECT().FindByID(ctx, matchID).Return(&domain.Match{ID: matchID, HomeTeamID: "team-1", AwayTeamID: "team-2", Stadium: "Gelora Bung Karno"}, nil)
	mockResultRepo.EXPECT().ExistsByMatchID(ctx, matchID).Return(false, nil)

	id, err := svc.ReportResult(ctx, matchID, result)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	assertMatchErrorCode(t, err, derrors.ErrorCodeBadRequest)
	if id != "" {
		t.Fatalf("expected empty ID on error, got %q", id)
	}
}

func TestMatchService_ReportResult_MissingTeamID(t *testing.T) {
	svc, mockMatchRepo, mockResultRepo, _ := setupMatchService(t)
	ctx := context.Background()
	matchID := "match-1"
	result := &domain.MatchResult{
		HomeScore: 1,
		AwayScore: 0,
		Goals: []domain.Goal{
			{PlayerID: "player-1", TeamID: "", GoalMinute: 10},
		},
	}

	mockMatchRepo.EXPECT().FindByID(ctx, matchID).Return(&domain.Match{ID: matchID, HomeTeamID: "team-1", AwayTeamID: "team-2", Stadium: "Gelora Bung Karno"}, nil)
	mockResultRepo.EXPECT().ExistsByMatchID(ctx, matchID).Return(false, nil)

	id, err := svc.ReportResult(ctx, matchID, result)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	assertMatchErrorCode(t, err, derrors.ErrorCodeBadRequest)
	if id != "" {
		t.Fatalf("expected empty ID on error, got %q", id)
	}
}

func TestMatchService_ReportResult_DrawWithNoGoals(t *testing.T) {
	svc, mockMatchRepo, mockResultRepo, _ := setupMatchService(t)
	ctx := context.Background()
	matchID := "match-1"
	result := &domain.MatchResult{
		HomeScore: 0,
		AwayScore: 0,
		Goals:     []domain.Goal{},
	}

	mockMatchRepo.EXPECT().FindByID(ctx, matchID).Return(&domain.Match{ID: matchID, HomeTeamID: "team-1", AwayTeamID: "team-2", Stadium: "Gelora Bung Karno"}, nil)
	mockResultRepo.EXPECT().ExistsByMatchID(ctx, matchID).Return(false, nil)
	mockResultRepo.EXPECT().Create(ctx, gomock.Any()).Return(nil)

	id, err := svc.ReportResult(ctx, matchID, result)

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if id == "" {
		t.Fatal("expected non-empty result ID, got empty string")
	}
}

func TestMatchService_ReportResult_HomeScoreMismatch(t *testing.T) {
	svc, mockMatchRepo, mockResultRepo, _ := setupMatchService(t)
	ctx := context.Background()
	matchID := "match-1"
	// Home score is 2, but only 1 home goal provided
	result := &domain.MatchResult{
		HomeScore: 2,
		AwayScore: 0,
		Goals: []domain.Goal{
			{PlayerID: "p1", TeamID: "team-1", GoalMinute: 10},
			{PlayerID: "p2", TeamID: "team-2", GoalMinute: 20}, // This makes it 1-1, but reported as 2-0
		},
	}

	mockMatchRepo.EXPECT().FindByID(ctx, matchID).Return(&domain.Match{ID: matchID, HomeTeamID: "team-1", AwayTeamID: "team-2", Stadium: "Gelora Bung Karno"}, nil)
	mockResultRepo.EXPECT().ExistsByMatchID(ctx, matchID).Return(false, nil)

	id, err := svc.ReportResult(ctx, matchID, result)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	assertMatchErrorCode(t, err, derrors.ErrorCodeBadRequest)
	if id != "" {
		t.Fatalf("expected empty ID on error, got %q", id)
	}
}

func TestMatchService_ReportResult_AwayScoreMismatch(t *testing.T) {
	svc, mockMatchRepo, mockResultRepo, _ := setupMatchService(t)
	ctx := context.Background()
	matchID := "match-1"
	// Away score is 1, but goals given to home team
	result := &domain.MatchResult{
		HomeScore: 0,
		AwayScore: 1,
		Goals: []domain.Goal{
			{PlayerID: "p1", TeamID: "team-1", GoalMinute: 10}, // Should have been team-2
		},
	}

	mockMatchRepo.EXPECT().FindByID(ctx, matchID).Return(&domain.Match{ID: matchID, HomeTeamID: "team-1", AwayTeamID: "team-2", Stadium: "Gelora Bung Karno"}, nil)
	mockResultRepo.EXPECT().ExistsByMatchID(ctx, matchID).Return(false, nil)

	id, err := svc.ReportResult(ctx, matchID, result)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	assertMatchErrorCode(t, err, derrors.ErrorCodeBadRequest)
	if id != "" {
		t.Fatalf("expected empty ID on error, got %q", id)
	}
}

func TestMatchService_ReportResult_InvalidGoalTeamID(t *testing.T) {
	svc, mockMatchRepo, mockResultRepo, _ := setupMatchService(t)
	ctx := context.Background()
	matchID := "match-1"
	// Goal for "team-3" which is not in the match
	result := &domain.MatchResult{
		HomeScore: 1,
		AwayScore: 0,
		Goals: []domain.Goal{
			{PlayerID: "p1", TeamID: "team-3", GoalMinute: 10},
		},
	}

	mockMatchRepo.EXPECT().FindByID(ctx, matchID).Return(&domain.Match{ID: matchID, HomeTeamID: "team-1", AwayTeamID: "team-2", Stadium: "Gelora Bung Karno"}, nil)
	mockResultRepo.EXPECT().ExistsByMatchID(ctx, matchID).Return(false, nil)

	id, err := svc.ReportResult(ctx, matchID, result)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	assertMatchErrorCode(t, err, derrors.ErrorCodeBadRequest)
	if id != "" {
		t.Fatalf("expected empty ID on error, got %q", id)
	}
}

func TestMatchService_ReportResult_SaveRepoError(t *testing.T) {
	svc, mockMatchRepo, mockResultRepo, _ := setupMatchService(t)
	ctx := context.Background()
	matchID := "match-1"
	result := &domain.MatchResult{
		HomeScore: 0,
		AwayScore: 0,
		Goals:     []domain.Goal{},
	}

	mockMatchRepo.EXPECT().FindByID(ctx, matchID).Return(&domain.Match{ID: matchID, HomeTeamID: "team-1", AwayTeamID: "team-2", Stadium: "Gelora Bung Karno"}, nil)
	mockResultRepo.EXPECT().ExistsByMatchID(ctx, matchID).Return(false, nil)
	mockResultRepo.EXPECT().Create(ctx, gomock.Any()).Return(derrors.WrapErrorf(errors.New("db error"), derrors.ErrorCodeInternal, "failed to save result"))

	id, err := svc.ReportResult(ctx, matchID, result)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	assertMatchErrorCode(t, err, derrors.ErrorCodeInternal)
	if id != "" {
		t.Fatalf("expected empty ID on error, got %q", id)
	}
}

// ---------------------------------------------------------------------------
// GetMatchReport
// ---------------------------------------------------------------------------

func TestMatchService_GetMatchReport_Success(t *testing.T) {
	svc, _, _, mockReportRepo := setupMatchService(t)
	ctx := context.Background()
	expected := &domain.MatchReportView{
		MatchID:      "match-1",
		HomeTeamName: "Team A",
		AwayTeamName: "Team B",
		HomeScore:    2,
		AwayScore:    1,
		MatchStatus:  "Home Win",
		TopScorer:    "Messi",
	}

	mockReportRepo.EXPECT().GetMatchReport(ctx, "match-1").Return(expected, nil)

	report, err := svc.GetMatchReport(ctx, "match-1")

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if report == nil {
		t.Fatal("expected report, got nil")
	}
	if report.MatchStatus != "Home Win" {
		t.Fatalf("expected status %q, got %q", "Home Win", report.MatchStatus)
	}
	if report.TopScorer != "Messi" {
		t.Fatalf("expected top scorer %q, got %q", "Messi", report.TopScorer)
	}
}

func TestMatchService_GetMatchReport_NotFound(t *testing.T) {
	svc, _, _, mockReportRepo := setupMatchService(t)
	ctx := context.Background()

	mockReportRepo.EXPECT().GetMatchReport(ctx, "nonexistent").Return(nil, domain.ErrMatchResultNotFound)

	report, err := svc.GetMatchReport(ctx, "nonexistent")

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, domain.ErrMatchResultNotFound) {
		t.Fatalf("expected ErrMatchResultNotFound, got: %v", err)
	}
	if report != nil {
		t.Fatalf("expected nil report on error, got %+v", report)
	}
}

// ---------------------------------------------------------------------------
// GetAllMatchReports
// ---------------------------------------------------------------------------

func TestMatchService_GetAllMatchReports_Success(t *testing.T) {
	svc, _, _, mockReportRepo := setupMatchService(t)
	ctx := context.Background()
	expected := []domain.MatchReportView{
		{MatchID: "match-1", HomeTeamName: "Team A", AwayTeamName: "Team B", MatchStatus: "Home Win"},
		{MatchID: "match-2", HomeTeamName: "Team C", AwayTeamName: "Team D", MatchStatus: "Draw"},
	}

	mockReportRepo.EXPECT().GetAllMatchReports(ctx).Return(expected, nil)

	reports, err := svc.GetAllMatchReports(ctx)

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if len(reports) != 2 {
		t.Fatalf("expected 2 reports, got %d", len(reports))
	}
}

func TestMatchService_GetAllMatchReports_RepoError(t *testing.T) {
	svc, _, _, mockReportRepo := setupMatchService(t)
	ctx := context.Background()

	mockReportRepo.EXPECT().GetAllMatchReports(ctx).Return(nil, derrors.WrapErrorf(errors.New("db error"), derrors.ErrorCodeInternal, "failed to fetch reports"))

	reports, err := svc.GetAllMatchReports(ctx)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	assertMatchErrorCode(t, err, derrors.ErrorCodeInternal)
	if reports != nil {
		t.Fatalf("expected nil reports on error, got %d items", len(reports))
	}
}
