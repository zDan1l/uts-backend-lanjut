package helper

import (
	"github.com/gofiber/fiber/v2"
)

// Format sukses seragam: {success, message, data, meta}
func Success(c *fiber.Ctx, status int, message string, data interface{}, meta ...interface{}) error {
	body := fiber.Map{
		"success": true,
		"message": message,
		"data":    data,
	}
	if len(meta) > 0 && meta[0] != nil {
		body["meta"] = meta[0]
	}
	return c.Status(status).JSON(body)
}

// Format error seragam: {success:false, message} atau + errors
func Error(c *fiber.Ctx, status int, message string, errors ...interface{}) error {
	body := fiber.Map{
		"success": false,
		"message": message,
	}
	if len(errors) > 0 && errors[0] != nil {
		body["errors"] = errors[0]
	}
	return c.Status(status).JSON(body)
}
