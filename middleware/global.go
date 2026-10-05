package middleware

import (
	"io"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

// Register memasang middleware global: recover + request logger ke logs/app.log.
func Register(app *fiber.App, out io.Writer) {
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Output: out,
		Format: "${time} ${method} ${path} ${status} ${latency}\n",
	}))
}
