package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"siakad-mini/helper"
)

// RequireJSON memastikan request POST/PUT memakai Content-Type JSON.
func RequireJSON(c *fiber.Ctx) error {
	if c.Method() == "POST" || c.Method() == "PUT" {
		if !strings.Contains(c.Get("Content-Type"), "application/json") {
			return helper.Error(c, 422, "Validasi gagal", fiber.Map{
				"body": []string{"Content-Type harus application/json"},
			})
		}
	}
	return c.Next()
}
