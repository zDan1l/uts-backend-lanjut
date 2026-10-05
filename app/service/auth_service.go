package service

import (
	"strings"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"siakad-mini/app/model"
	"siakad-mini/app/repository"
	"siakad-mini/config"
	"siakad-mini/helper"
	"siakad-mini/middleware"
)

// AuthService menangani login + profil. Menerima fiber.Ctx langsung.
type AuthService struct {
	App *config.App
}

// Rate limit sederhana: max 5x gagal login per menit per IP
var (
	failMu   sync.Mutex
	failLog  = map[string][]time.Time{}
	failMax  = 5
	failWind = time.Minute
)

func failCount(ip string) int {
	failMu.Lock()
	defer failMu.Unlock()
	now := time.Now()
	list := failLog[ip]
	fresh := list[:0]
	for _, t := range list {
		if now.Sub(t) < failWind {
			fresh = append(fresh, t)
		}
	}
	failLog[ip] = fresh
	return len(fresh)
}

func failAdd(ip string) {
	failMu.Lock()
	defer failMu.Unlock()
	failLog[ip] = append(failLog[ip], time.Now())
}

func failReset(ip string) {
	failMu.Lock()
	defer failMu.Unlock()
	delete(failLog, ip)
}

func (s *AuthService) Login(c *fiber.Ctx) error {
	if failCount(c.IP()) >= failMax {
		return helper.Error(c, 429, "Terlalu banyak percobaan login. Coba lagi nanti")
	}

	var b model.LoginRequest
	if err := c.BodyParser(&b); err != nil {
		return helper.Error(c, 422, "Validasi gagal", fiber.Map{
			"email": []string{"Format request tidak valid"},
		})
	}
	b.Email = strings.TrimSpace(b.Email)

	errs := map[string][]string{}
	if b.Email == "" {
		errs["email"] = append(errs["email"], "Email wajib diisi")
	} else if !strings.Contains(b.Email, "@") || !strings.Contains(b.Email, ".") {
		errs["email"] = append(errs["email"], "Format email tidak valid")
	}
	if b.Password == "" {
		errs["password"] = append(errs["password"], "Password wajib diisi")
	} else if len(b.Password) < 8 {
		errs["password"] = append(errs["password"], "Password minimal 8 karakter")
	}
	if len(errs) > 0 {
		return helper.Error(c, 422, "Validasi gagal", errs)
	}

	user, _, err := repository.FindUserByEmail(s.App.DB, b.Email)
	if err != nil || !helper.CheckPassword(b.Password, user.Password) {
		failAdd(c.IP())
		if failCount(c.IP()) >= failMax {
			return helper.Error(c, 429, "Terlalu banyak percobaan login. Coba lagi nanti")
		}
		return helper.Error(c, 401, "Email atau password salah")
	}

	failReset(c.IP())
	token, expIn, err := helper.GenerateToken(s.App.Env.JWTSecret, s.App.Env.JWTExpireH, user.ID, user.Email, user.Role)
	if err != nil {
		return helper.Error(c, 500, "Terjadi kesalahan pada server")
	}
	return helper.Success(c, 200, "Login berhasil", fiber.Map{
		"access_token": token,
		"token_type":   "Bearer",
		"expires_in":   expIn,
		"user": fiber.Map{
			"id":    user.ID,
			"email": user.Email,
			"role":  user.Role,
		},
	})
}

func (s *AuthService) Me(c *fiber.Ctx) error {
	id, _, _ := middleware.CurrentUser(c)

	user, err := repository.FindUserByID(s.App.DB, id)
	if err != nil {
		return helper.Error(c, 401, "Token tidak ada, salah, atau kedaluwarsa")
	}
	data := fiber.Map{
		"id":    user.ID,
		"email": user.Email,
		"role":  user.Role,
	}
	if user.Role == "mahasiswa" {
		st, err := repository.FindStudentByUserID(s.App.DB, user.ID)
		if err != nil {
			return helper.Error(c, 401, "Akun mahasiswa tidak ditemukan")
		}
		data["student"] = fiber.Map{
			"id": st.ID, "nim": st.NIM, "nama": st.Nama,
			"prodi": st.Prodi, "angkatan": st.Angkatan,
		}
	}
	return helper.Success(c, 200, "Profil berhasil diambil", data)
}
