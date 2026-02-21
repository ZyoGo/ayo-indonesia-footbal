package middleware

import (
	"net/http"
	"strings"

	common "github.com/ZyoGo/ayo-indonesia-footbal/pkg/http"
	"github.com/ZyoGo/ayo-indonesia-footbal/pkg/jwt"
	"github.com/gin-gonic/gin"
)

var (
	UserAttr     = "userAttr"
	PrefixHeader = "Bearer "
	InvalidToken = "Invalid / expired token"
)

type AuthGuard struct {
	j *jwt.Service
}

func NewAuthGuard(j *jwt.Service) *AuthGuard {
	return &AuthGuard{j: j}
}

// Guard is the middleware function to verify JWT token.
func (g *AuthGuard) Guard() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if !strings.HasPrefix(authHeader, PrefixHeader) {
			c.JSON(http.StatusUnauthorized, common.NewUnauthorizedResponse("Authorization header missing/invalid"))
			c.Abort()
			return
		}

		attr, err := g.j.ParseAndVerify(strings.TrimPrefix(authHeader, PrefixHeader))
		if err != nil {
			c.JSON(http.StatusUnauthorized, common.NewUnauthorizedResponse(InvalidToken))
			c.Abort()
			return
		}

		c.Set(UserAttr, attr)
		c.Next()
	}
}
