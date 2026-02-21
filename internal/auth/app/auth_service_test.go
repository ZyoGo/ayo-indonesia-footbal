package app

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"errors"
	"testing"

	"github.com/ZyoGo/ayo-indonesia-footbal/internal/auth/domain"
	mockDomain "github.com/ZyoGo/ayo-indonesia-footbal/internal/auth/mock"
	"github.com/ZyoGo/ayo-indonesia-footbal/pkg/derrors"
	"github.com/ZyoGo/ayo-indonesia-footbal/pkg/jwt"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
)

func setupAuthService(t *testing.T) (*AuthService, *mockDomain.MockUserRepository) {
	t.Helper()
	ctrl := gomock.NewController(t)
	mockRepo := mockDomain.NewMockUserRepository(ctrl)

	// Generate a test RSA key pair
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}

	jwtService := jwt.NewService(privateKey, &privateKey.PublicKey, "test-issuer", "test-subject")
	svc := &AuthService{userRepo: mockRepo, jwtService: jwtService}
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

func hashPassword(t *testing.T, password string) string {
	t.Helper()
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatal(err)
	}
	return string(hash)
}

// ---------------------------------------------------------------------------
// Login
// ---------------------------------------------------------------------------

func TestAuthService_Login_Success(t *testing.T) {
	// Given
	svc, mockRepo := setupAuthService(t)
	ctx := context.Background()

	mockRepo.EXPECT().FindByUsername(ctx, "admin").Return(&domain.User{
		ID:           "user-1",
		Username:     "admin",
		PasswordHash: hashPassword(t, "admin123"),
	}, nil)

	// When
	token, err := svc.Login(ctx, "admin", "admin123")

	// Then
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if token == "" {
		t.Fatal("expected non-empty token, got empty string")
	}
}

func TestAuthService_Login_WrongPassword(t *testing.T) {
	// Given
	svc, mockRepo := setupAuthService(t)
	ctx := context.Background()

	mockRepo.EXPECT().FindByUsername(ctx, "admin").Return(&domain.User{
		ID:           "user-1",
		Username:     "admin",
		PasswordHash: hashPassword(t, "admin123"),
	}, nil)

	// When
	token, err := svc.Login(ctx, "admin", "wrong-password")

	// Then
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	assertErrorCode(t, err, derrors.ErrorCodeUnauthorized)
	if token != "" {
		t.Fatalf("expected empty token on error, got %q", token)
	}
}

func TestAuthService_Login_UserNotFound(t *testing.T) {
	// Given
	svc, mockRepo := setupAuthService(t)
	ctx := context.Background()

	mockRepo.EXPECT().FindByUsername(ctx, "unknown").Return(nil, domain.ErrUserNotFound)

	// When
	token, err := svc.Login(ctx, "unknown", "password")

	// Then
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	assertErrorCode(t, err, derrors.ErrorCodeUnauthorized)
	if token != "" {
		t.Fatalf("expected empty token on error, got %q", token)
	}
}

func TestAuthService_Login_EmptyUsername(t *testing.T) {
	// Given
	svc, _ := setupAuthService(t)
	ctx := context.Background()

	// When
	token, err := svc.Login(ctx, "", "password")

	// Then
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	assertErrorCode(t, err, derrors.ErrorCodeBadRequest)
	if token != "" {
		t.Fatalf("expected empty token on error, got %q", token)
	}
}

func TestAuthService_Login_EmptyPassword(t *testing.T) {
	// Given
	svc, _ := setupAuthService(t)
	ctx := context.Background()

	// When
	token, err := svc.Login(ctx, "admin", "")

	// Then
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	assertErrorCode(t, err, derrors.ErrorCodeBadRequest)
	if token != "" {
		t.Fatalf("expected empty token on error, got %q", token)
	}
}

func TestAuthService_Login_RepoError(t *testing.T) {
	// Given
	svc, mockRepo := setupAuthService(t)
	ctx := context.Background()

	mockRepo.EXPECT().FindByUsername(ctx, "admin").Return(nil,
		derrors.WrapErrorf(errors.New("db error"), derrors.ErrorCodeInternal, "database connection failed"))

	// When
	token, err := svc.Login(ctx, "admin", "password")

	// Then
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	assertErrorCode(t, err, derrors.ErrorCodeInternal)
	if token != "" {
		t.Fatalf("expected empty token on error, got %q", token)
	}
}
