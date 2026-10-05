package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"siakad-mini/helper"
)

// Wajib ada header Authorization: Bearer <token>
func AuthRequired(secret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		h := c.Get("Authorization")
		if h == "" || !strings.HasPrefix(h, "Bearer ") {
			return helper.Error(c, 401, "Token tidak ada atau tidak valid")
		}
		claims, err := helper.ParseToken(secret, strings.TrimSpace(strings.TrimPrefix(h, "Bearer ")))
		if err != nil {
			return helper.Error(c, 401, "Token tidak ada, salah, atau kedaluwarsa")
		}
		c.Locals("userID", claims.UserID)
		c.Locals("userRole", claims.Role)
		c.Locals("userEmail", claims.Email)
		return c.Next()
	}
}

// Batasi role tertentu, misal admin saja
func RequireRole(roles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		role, _ := c.Locals("userRole").(string)
		for _, r := range roles {
			if role == r {
				return c.Next()
			}
		}
		return helper.Error(c, 403, "Akses ditolak")
	}
}

func CurrentUser(c *fiber.Ctx) (id int, role string, email string) {
	id, _ = c.Locals("userID").(int)
	role, _ = c.Locals("userRole").(string)
	email, _ = c.Locals("userEmail").(string)
	return
}
