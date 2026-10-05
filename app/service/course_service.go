package service

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"siakad-mini/app/repository"
	"siakad-mini/config"
	"siakad-mini/helper"
)

// CourseService menangani daftar mata kuliah.
type CourseService struct {
	App *config.App
}

// GET /courses — semua role yang login
func (s *CourseService) List(c *fiber.Ctx) error {
	f := repository.CourseFilter{
		Semester:  helper.QueryInt(c, "semester", 0),
		Search:    strings.TrimSpace(helper.QueryStr(c, "search", "")),
		Available: strings.ToLower(helper.QueryStr(c, "available", "")) == "true",
	}
	data, err := repository.ListCourses(s.App.DB, f)
	if err != nil {
		return helper.Error(c, 500, "Terjadi kesalahan pada server")
	}
	return helper.Success(c, 200, "Data mata kuliah berhasil diambil", data)
}
