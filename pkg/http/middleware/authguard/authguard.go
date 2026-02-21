package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

var (
	UserAttr     = "userAttr"
	PrefixHeader = "Bearer "
	InvalidToken = "Invalid / expired token"
)

type jwt interface {
	ParseAndVerify(accessToken string) (JwtAttr, error)
}

type JwtAttr struct {
	Email    string
}

type AuthGuard struct {
	j jwt
}

func NewAuthGuard(j jwt) *AuthGuard {
	return &AuthGuard{j: j}
}

// Guard is the middleware function to verify JWT token.
func (g *AuthGuard) Guard() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if !strings.HasPrefix(authHeader, PrefixHeader) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing/invalid"})
			c.Abort()
			return
		}

		attr, err := g.j.ParseAndVerify(strings.TrimPrefix(authHeader, PrefixHeader))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": InvalidToken})
			c.Abort()
			return
		}

		c.Set(UserAttr, attr)
		c.Next()
	}
}
