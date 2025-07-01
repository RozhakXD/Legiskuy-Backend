package middleware

import (
	"os"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func Protected() fiber.Handler {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "zf9i95p(x_@^72-j_tj&=&(j^j&_7ku2d^vdpiotmu#gb*&h(3"
	}

	return jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(jwtSecret)},
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized",
			})
		},
	})
}

func RequireRole(requiredRole string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := c.Locals("user").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)

		role, ok := claims["role"].(string)
		if !ok {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Role not found in token",
			})
		}

		if role != requiredRole {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Forbidden: insufficient permissions",
			})
		}

		return c.Next()
	}
}
