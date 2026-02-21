package jwt

import (
	"crypto/rsa"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JwtAttr represents JWT attributes
type JwtAttr struct {
	Email string
}

// Service handles JWT signing and verification using RSA
type Service struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
	issuer     string
	subject    string
}

// NewService creates a new JWT service with RSA keys
func NewService(privateKey *rsa.PrivateKey, publicKey *rsa.PublicKey, issuer, subject string) *Service {
	return &Service{
		privateKey: privateKey,
		publicKey:  publicKey,
		issuer:     issuer,
		subject:    subject,
	}
}

type Claims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

// GenerateToken creates a new JWT token with RSA signing (RS256)
func (s *Service) GenerateToken(attr JwtAttr, expiry time.Duration) (string, error) {
	now := time.Now()

	claims := Claims{
		Email: attr.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(expiry)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    s.issuer,
			Subject:   s.subject,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(s.privateKey)
}

// ParseAndVerify validates a JWT token and returns the claims
func (s *Service) ParseAndVerify(tokenString string) (JwtAttr, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return s.publicKey, nil
	})

	if err != nil {
		return JwtAttr{}, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return JwtAttr{
			Email: claims.Email,
		}, nil
	}

	return JwtAttr{}, errors.New("invalid token")
}
