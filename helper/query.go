package helper

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// QueryInt membaca query string angka dengan nilai default.
func QueryInt(c *fiber.Ctx, key string, def int) int {
	v, err := strconv.Atoi(c.Query(key, ""))
	if err != nil {
		return def
	}
	return v
}

// QueryStr membaca query string dengan nilai default.
func QueryStr(c *fiber.Ctx, key, def string) string {
	if v := c.Query(key, ""); v != "" {
		return v
	}
	return def
}
