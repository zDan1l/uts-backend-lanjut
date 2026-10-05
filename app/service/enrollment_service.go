package service

import (
	"regexp"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"siakad-mini/app/model"
	"siakad-mini/app/repository"
	"siakad-mini/config"
	"siakad-mini/helper"
	"siakad-mini/middleware"
)

// EnrollmentService menangani KRS. Menerima fiber.Ctx langsung.
type EnrollmentService struct {
	App *config.App
}

var tahunRe = regexp.MustCompile(`^\d{4}/\d{4}-(Ganjil|Genap)$`)

// POST /enrollments — mahasiswa saja
func (s *EnrollmentService) Create(c *fiber.Ctx) error {
	uid, _, _ := middleware.CurrentUser(c)
	mine, err := repository.FindStudentByUserID(s.App.DB, uid)
	if err != nil {
		return helper.Error(c, 403, "Akses ditolak")
	}

	var b model.EnrollRequest
	if err := c.BodyParser(&b); err != nil {
		return helper.Error(c, 422, "Validasi gagal", fiber.Map{
			"course_id": []string{"Course ID wajib diisi"},
		})
	}

	errs := map[string][]string{}
	if b.CourseID <= 0 {
		errs["course_id"] = []string{"Course ID wajib diisi dan harus ada"}
	}
	if b.TahunAkademik == "" {
		errs["tahun_akademik"] = []string{"Tahun akademik wajib diisi"}
	} else if !tahunRe.MatchString(b.TahunAkademik) {
		errs["tahun_akademik"] = []string{"Format tahun akademik tidak valid (contoh: 2026/2027-Ganjil)"}
	}
	if len(errs) > 0 {
		return helper.Error(c, 422, "Validasi gagal", errs)
	}

	newID, err := repository.CreateEnrollmentTx(s.App.DB, mine.ID, b.CourseID, b.TahunAkademik)
	if err != nil {
		if ee, ok := err.(*repository.EnrollError); ok {
			if ee.Status == 409 {
				return helper.Error(c, 409, ee.Message)
			}
			return helper.Error(c, ee.Status, ee.Message, ee.Errors)
		}
		return helper.Error(c, 500, "Terjadi kesalahan pada server")
	}

	e, _ := repository.FindEnrollmentByID(s.App.DB, newID)
	return helper.Success(c, 201, "Mata kuliah berhasil diambil", e)
}

// DELETE /enrollments/{id} — mahasiswa pemilik saja
func (s *EnrollmentService) Delete(c *fiber.Ctx) error {
	uid, _, _ := middleware.CurrentUser(c)
	mine, err := repository.FindStudentByUserID(s.App.DB, uid)
	if err != nil {
		return helper.Error(c, 403, "Akses ditolak")
	}
	id, _ := strconv.Atoi(c.Params("id"))
	if id <= 0 {
		return helper.Error(c, 404, "Enrollment tidak ditemukan")
	}
	e, err := repository.FindEnrollmentByID(s.App.DB, id)
	if err != nil {
		return helper.Error(c, 404, "Enrollment tidak ditemukan")
	}
	if e.StudentID != mine.ID {
		return helper.Error(c, 403, "Akses ditolak")
	}
	if err := repository.DeleteEnrollment(s.App.DB, id); err != nil {
		return helper.Error(c, 500, "Terjadi kesalahan pada server")
	}
	return c.SendStatus(204)
}
