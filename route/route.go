package route

import (
	"github.com/gofiber/fiber/v2"
	"siakad-mini/app/service"
	"siakad-mini/config"
	"siakad-mini/middleware"
)

// Register mendaftarkan 10 endpoint SIAKAD Mini.
func Register(app *fiber.App, deps *config.App) {
	auth := &service.AuthService{App: deps}
	stud := &service.StudentService{App: deps}
	course := &service.CourseService{App: deps}
	enroll := &service.EnrollmentService{App: deps}

	api := app.Group("/api/v1")

	// 1. Publik
	api.Post("/auth/login", middleware.RequireJSON, auth.Login)

	// 2-10. Butuh token
	authed := api.Group("", middleware.AuthRequired(deps.Env.JWTSecret))
	authed.Get("/auth/me", auth.Me)

	// 3-7. Students
	authed.Get("/students", middleware.RequireRole("admin"), stud.List)
	authed.Post("/students", middleware.RequireJSON, middleware.RequireRole("admin"), stud.Create)
	authed.Get("/students/:id", stud.Detail) // cek pemilik di service
	authed.Put("/students/:id", middleware.RequireJSON, middleware.RequireRole("admin"), stud.Update)
	authed.Delete("/students/:id", middleware.RequireRole("admin"), stud.Delete)

	// 8. Courses (semua role yang login)
	authed.Get("/courses", course.List)

	// 9-10. Enrollments (mahasiswa, cek pemilik di service)
	authed.Post("/enrollments", middleware.RequireJSON, middleware.RequireRole("mahasiswa"), enroll.Create)
	authed.Delete("/enrollments/:id", middleware.RequireRole("mahasiswa"), enroll.Delete)
}
