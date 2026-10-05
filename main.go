package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"siakad-mini/config"
	"siakad-mini/database"
	"siakad-mini/helper"
	"siakad-mini/middleware"
	"siakad-mini/route"
)

func main() {
	env := config.Load()
	lg, logWriter := config.SetupLogger()

	db, err := database.Connect(env)
	if err != nil {
		lg.Fatalf("gagal konek database: %v", err)
	}
	defer db.Close()

	if err := database.Migrate(db); err != nil {
		lg.Fatalf("gagal migrate: %v", err)
	}
	lg.Println("migration OK")

	// Seeder: go run . --seed
	if len(os.Args) > 1 && os.Args[1] == "--seed" {
		if err := database.Seed(db); err != nil {
			lg.Fatalf("gagal seed: %v", err)
		}
		return
	}

	deps := config.New(db, env, lg)

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			if code == 404 {
				return helper.Error(c, 404, "Endpoint tidak ditemukan")
			}
			if env.AppEnv == "production" {
				return helper.Error(c, 500, "Terjadi kesalahan pada server")
			}
			return helper.Error(c, code, err.Error())
		},
	})
	middleware.Register(app, logWriter)

	app.Get("/", func(c *fiber.Ctx) error {
		return helper.Success(c, 200, "SIAKAD Mini API berjalan", fiber.Map{"version": "v1"})
	})

	route.Register(app, deps)

	lg.Printf("server jalan di :%s", env.AppPort)
	log.Fatal(app.Listen(":" + env.AppPort))
}
