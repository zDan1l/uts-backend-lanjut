package service

import (
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"siakad-mini/app/model"
	"siakad-mini/app/repository"
	"siakad-mini/config"
	"siakad-mini/helper"
	"siakad-mini/middleware"
)

// StudentService menangani CRUD mahasiswa. Menerima fiber.Ctx langsung.
type StudentService struct {
	App *config.App
}

// GET /students — admin saja
func (s *StudentService) List(c *fiber.Ctx) error {
	page := helper.QueryInt(c, "page", 1)
	perPage := helper.QueryInt(c, "per_page", 10)
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 10
	}
	if perPage > 50 {
		perPage = 50
	}
	f := repository.StudentFilter{
		Page: page, PerPage: perPage,
		Prodi:    helper.QueryStr(c, "prodi", ""),
		Angkatan: helper.QueryInt(c, "angkatan", 0),
		Search:   helper.QueryStr(c, "search", ""),
		Sort:     helper.QueryStr(c, "sort", ""),
	}

	data, total, err := repository.ListStudents(s.App.DB, f)
	if err != nil {
		return helper.Error(c, 500, "Terjadi kesalahan pada server")
	}
	lastPage := (total + perPage - 1) / perPage
	if lastPage == 0 {
		lastPage = 1
	}
	return helper.Success(c, 200, "Data mahasiswa berhasil diambil", data, fiber.Map{
		"current_page": page, "per_page": perPage,
		"total": total, "last_page": lastPage,
	})
}

func isDigits(str string) bool {
	if str == "" {
		return false
	}
	for _, r := range str {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

// POST /students — admin saja
func (s *StudentService) Create(c *fiber.Ctx) error {
	var b model.StudentCreateRequest
	if err := c.BodyParser(&b); err != nil {
		return helper.Error(c, 422, "Validasi gagal", fiber.Map{"nim": []string{"Format request tidak valid"}})
	}
	b.NIM = strings.TrimSpace(b.NIM)
	b.Nama = strings.TrimSpace(b.Nama)
	b.Email = strings.TrimSpace(b.Email)
	b.Prodi = strings.TrimSpace(b.Prodi)

	errs := map[string][]string{}
	if b.NIM == "" {
		errs["nim"] = []string{"NIM wajib diisi"}
	} else if len(b.NIM) != 12 || !isDigits(b.NIM) {
		errs["nim"] = []string{"NIM harus 12 digit angka"}
	} else if repository.NIMExists(s.App.DB, b.NIM) {
		errs["nim"] = []string{"NIM sudah terdaftar"}
	}
	if b.Nama == "" {
		errs["nama"] = []string{"Nama wajib diisi"}
	}
	if b.Email == "" {
		errs["email"] = []string{"Email wajib diisi"}
	} else if !strings.Contains(b.Email, "@") || !strings.Contains(b.Email, ".") {
		errs["email"] = []string{"Format email tidak valid"}
	} else if repository.EmailExists(s.App.DB, b.Email) {
		errs["email"] = []string{"Email sudah terdaftar"}
	}
	if b.Prodi == "" {
		errs["prodi"] = []string{"Prodi wajib diisi"}
	}
	nowYear := time.Now().Year()
	if b.Angkatan == 0 {
		errs["angkatan"] = []string{"Angkatan wajib diisi"}
	} else if len(strconv.Itoa(b.Angkatan)) != 4 || b.Angkatan > nowYear {
		errs["angkatan"] = []string{"Angkatan harus 4 digit dan tidak melebihi tahun berjalan"}
	}
	if b.IPKTerakhir < 0 || b.IPKTerakhir > 4 {
		errs["ipk_terakhir"] = []string{"IPK harus antara 0,00–4,00"}
	}
	if len(errs) > 0 {
		return helper.Error(c, 422, "Validasi gagal", errs)
	}

	// Business rule: password awal = NIM
	hash, err := helper.HashPassword(b.NIM)
	if err != nil {
		return helper.Error(c, 500, "Terjadi kesalahan pada server")
	}
	_, sid, err := repository.CreateStudentTx(s.App.DB, b.Email, hash, b.NIM, b.Nama, b.Prodi, b.Angkatan, b.IPKTerakhir)
	if err != nil {
		return helper.Error(c, 500, "Terjadi kesalahan pada server")
	}
	detail, _ := repository.GetStudentDetail(s.App.DB, sid)
	return helper.Success(c, 201, "Mahasiswa berhasil ditambahkan", detail)
}

// GET /students/{id} — admin atau pemilik data
func (s *StudentService) Detail(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	if id <= 0 {
		return helper.Error(c, 404, "Mahasiswa tidak ditemukan")
	}
	uid, role, _ := middleware.CurrentUser(c)
	if role == "mahasiswa" {
		mine, err := repository.FindStudentByUserID(s.App.DB, uid)
		if err != nil || mine.ID != id {
			return helper.Error(c, 403, "Akses ditolak")
		}
	}
	detail, err := repository.GetStudentDetail(s.App.DB, id)
	if err != nil {
		return helper.Error(c, 404, "Mahasiswa tidak ditemukan")
	}
	return helper.Success(c, 200, "Detail mahasiswa berhasil diambil", detail)
}

// PUT /students/{id} — admin saja (nim tidak boleh diubah)
func (s *StudentService) Update(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	if id <= 0 {
		return helper.Error(c, 404, "Mahasiswa tidak ditemukan")
	}
	old, err := repository.GetStudentDetail(s.App.DB, id)
	if err != nil {
		return helper.Error(c, 404, "Mahasiswa tidak ditemukan")
	}
	var b model.StudentUpdateRequest
	if err := c.BodyParser(&b); err != nil {
		return helper.Error(c, 422, "Validasi gagal")
	}
	nama := old.Nama
	if strings.TrimSpace(b.Nama) != "" {
		nama = strings.TrimSpace(b.Nama)
	}
	prodi := old.Prodi
	if strings.TrimSpace(b.Prodi) != "" {
		prodi = strings.TrimSpace(b.Prodi)
	}
	angkatan := old.Angkatan
	if b.Angkatan != 0 {
		angkatan = b.Angkatan
	}
	ipk := old.IPKTerakhir
	if b.IPKTerakhir != nil {
		ipk = *b.IPKTerakhir
	}

	errs := map[string][]string{}
	if strings.TrimSpace(nama) == "" {
		errs["nama"] = []string{"Nama wajib diisi"}
	}
	if strings.TrimSpace(prodi) == "" {
		errs["prodi"] = []string{"Prodi wajib diisi"}
	}
	if len(strconv.Itoa(angkatan)) != 4 || angkatan > time.Now().Year() {
		errs["angkatan"] = []string{"Angkatan harus 4 digit dan tidak melebihi tahun berjalan"}
	}
	if ipk < 0 || ipk > 4 {
		errs["ipk_terakhir"] = []string{"IPK harus antara 0,00–4,00"}
	}
	if len(errs) > 0 {
		return helper.Error(c, 422, "Validasi gagal", errs)
	}

	if err := repository.UpdateStudent(s.App.DB, id, nama, prodi, angkatan, ipk); err != nil {
		return helper.Error(c, 404, "Mahasiswa tidak ditemukan")
	}
	detail, _ := repository.GetStudentDetail(s.App.DB, id)
	return helper.Success(c, 200, "Mahasiswa berhasil diperbarui", detail)
}

// DELETE /students/{id} — admin saja (soft delete)
func (s *StudentService) Delete(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	if id <= 0 {
		return helper.Error(c, 404, "Mahasiswa tidak ditemukan")
	}
	if err := repository.SoftDeleteStudent(s.App.DB, id); err != nil {
		return helper.Error(c, 404, "Mahasiswa tidak ditemukan")
	}
	return c.SendStatus(204)
}
