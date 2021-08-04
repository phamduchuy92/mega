package middleware

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v2"
	"github.com/golang-jwt/jwt"
	"gitlab.com/emi2/mega/internal/app"
)

// Protected protect routes
func Protected() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey:   []byte(app.Config.String("security.jwt-secret")),
		ErrorHandler: jwtError,
	})
}

// HasAuthority check if the current role has specified authorities
func HasAuthority(authorityName string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := c.Locals("user").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)
		if authorities, ok := claims["authorities"]; ok && (authorities != nil) {
			authorities := claims["authorities"].([]interface{})
			for _, authority := range authorities {
				if fmt.Sprint(authority) == authorityName {
					return c.Next()
				}
			}
		}
		return fiber.NewError(fiber.StatusForbidden, fmt.Sprintf(`Account doesn't have required authority %s to access this resource`, authorityName))
	}
}

func jwtError(c *fiber.Ctx, err error) error {
	if err.Error() == "Missing or malformed JWT" {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"status": "error", "message": "Missing or malformed JWT", "data": nil})
	}
	return c.Status(fiber.StatusUnauthorized).
		JSON(fiber.Map{"status": "error", "message": "Invalid or expired JWT", "data": nil})
}
