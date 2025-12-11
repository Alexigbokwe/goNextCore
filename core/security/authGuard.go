package security

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

type AuthGuard struct {
	JwtService *JwtService `inject:"type"`
}

func (g *AuthGuard) CanActivate(c *fiber.Ctx) bool {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return false
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return false
	}

	token := parts[1]
	claims, err := g.JwtService.Verify(token)
	if err != nil {
		return false
	}

	// Store user in locals
	c.Locals("user", claims)
	return true
}
