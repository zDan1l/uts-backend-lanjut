package model

import "time"

// ===== Entitas (tabel database) =====

type User struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

type Student struct {
	ID          int     `json:"id"`
	UserID      int     `json:"user_id,omitempty"`
	NIM         string  `json:"nim"`
	Nama        string  `json:"nama"`
	Prodi       string  `json:"prodi"`
	Angkatan    int     `json:"angkatan"`
	IPKTerakhir float64 `json:"ipk_terakhir"`
	DeletedAt   *string `json:"deleted_at,omitempty"`
}

type Course struct {
	ID        int    `json:"id"`
	KodeMK    string `json:"kode_mk"`
	NamaMK    string `json:"nama_mk"`
	SKS       int    `json:"sks"`
	Semester  int    `json:"semester"`
	Kuota     int    `json:"kuota"`
	Terisi    int    `json:"terisi,omitempty"`
	SisaKuota int    `json:"sisa_kuota,omitempty"`
}

type Enrollment struct {
	ID            int    `json:"id"`
	StudentID     int    `json:"student_id"`
	CourseID      int    `json:"course_id"`
	TahunAkademik string `json:"tahun_akademik"`
}

// ===== Struct response =====

type StudentDetail struct {
	Student
	MataKuliah []CourseEnrollment `json:"mata_kuliah"`
	TotalSKS   int                `json:"total_sks"`
	BatasSKS   int                `json:"batas_sks"`
}

type CourseEnrollment struct {
	EnrollmentID  int    `json:"enrollment_id"`
	CourseID      int    `json:"course_id"`
	KodeMK        string `json:"kode_mk"`
	NamaMK        string `json:"nama_mk"`
	SKS           int    `json:"sks"`
	TahunAkademik string `json:"tahun_akademik"`
}

// ===== Struct request =====

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type StudentCreateRequest struct {
	NIM         string  `json:"nim"`
	Nama        string  `json:"nama"`
	Email       string  `json:"email"`
	Prodi       string  `json:"prodi"`
	Angkatan    int     `json:"angkatan"`
	IPKTerakhir float64 `json:"ipk_terakhir"`
}

type StudentUpdateRequest struct {
	Nama        string   `json:"nama"`
	Prodi       string   `json:"prodi"`
	Angkatan    int      `json:"angkatan"`
	IPKTerakhir *float64 `json:"ipk_terakhir"`
}

type EnrollRequest struct {
	CourseID      int    `json:"course_id"`
	TahunAkademik string `json:"tahun_akademik"`
}
