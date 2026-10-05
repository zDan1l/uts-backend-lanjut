package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"siakad-mini/app/model"
	"siakad-mini/helper"
)

type StudentFilter struct {
	Page     int
	PerPage  int
	Prodi    string
	Angkatan int // 0 = semua
	Search   string
	Sort     string // "nama" atau "-ipk_terakhir"
}

// Daftar mahasiswa (tanpa yang soft-delete) + total
func ListStudents(db *sql.DB, f StudentFilter) ([]model.Student, int, error) {
	where := []string{"deleted_at IS NULL"}
	args := []interface{}{}
	n := 1

	if f.Prodi != "" {
		where = append(where, fmt.Sprintf("prodi = $%d", n))
		args = append(args, f.Prodi)
		n++
	}
	if f.Angkatan > 0 {
		where = append(where, fmt.Sprintf("angkatan = $%d", n))
		args = append(args, f.Angkatan)
		n++
	}
	if f.Search != "" {
		where = append(where, fmt.Sprintf("(nim ILIKE $%d OR nama ILIKE $%d)", n, n+1))
		like := "%" + f.Search + "%"
		args = append(args, like, like)
		n += 2
	}
	w := strings.Join(where, " AND ")

	var total int
	if err := db.QueryRow(`SELECT COUNT(*) FROM students WHERE `+w, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	order := "ORDER BY id ASC"
	if f.Sort == "nama" {
		order = "ORDER BY nama ASC"
	} else if f.Sort == "-ipk_terakhir" {
		order = "ORDER BY ipk_terakhir DESC"
	}

	offset := (f.Page - 1) * f.PerPage
	args = append(args, f.PerPage, offset)
	q := fmt.Sprintf(`SELECT id, user_id, nim, nama, prodi, angkatan, ipk_terakhir
		FROM students WHERE %s %s LIMIT $%d OFFSET $%d`, w, order, n, n+1)

	rows, err := db.Query(q, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	out := []model.Student{}
	for rows.Next() {
		var s model.Student
		if err := rows.Scan(&s.ID, &s.UserID, &s.NIM, &s.Nama, &s.Prodi, &s.Angkatan, &s.IPKTerakhir); err != nil {
			return nil, 0, err
		}
		out = append(out, s)
	}
	if out == nil {
		out = []model.Student{}
	}
	return out, total, nil
}

// Detail 1 mahasiswa (404 jika soft-delete) + KRS + total & batas SKS
func GetStudentDetail(db *sql.DB, id int) (*model.StudentDetail, error) {
	var s model.Student
	err := db.QueryRow(`SELECT id, user_id, nim, nama, prodi, angkatan, ipk_terakhir
		FROM students WHERE id = $1 AND deleted_at IS NULL`, id,
	).Scan(&s.ID, &s.UserID, &s.NIM, &s.Nama, &s.Prodi, &s.Angkatan, &s.IPKTerakhir)
	if err != nil {
		return nil, err
	}

	rows, err := db.Query(`SELECT e.id, c.id, c.kode_mk, c.nama_mk, c.sks, e.tahun_akademik
		FROM enrollments e JOIN courses c ON c.id = e.course_id
		WHERE e.student_id = $1 ORDER BY e.id ASC`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	mk := []model.CourseEnrollment{}
	total := 0
	for rows.Next() {
		var m model.CourseEnrollment
		if err := rows.Scan(&m.EnrollmentID, &m.CourseID, &m.KodeMK, &m.NamaMK, &m.SKS, &m.TahunAkademik); err != nil {
			return nil, err
		}
		mk = append(mk, m)
		total += m.SKS
	}
	if mk == nil {
		mk = []model.CourseEnrollment{}
	}
	return &model.StudentDetail{
		Student:    s,
		MataKuliah: mk,
		TotalSKS:   total,
		BatasSKS:   helper.BatasSKS(s.IPKTerakhir),
	}, nil
}

// Buat users + students dalam 1 transaction. Password awal = NIM (di-hash di service).
func CreateStudentTx(db *sql.DB, email, passHash, nim, nama, prodi string, angkatan int, ipk float64) (userID, studentID int, err error) {
	tx, err := db.Begin()
	if err != nil {
		return 0, 0, err
	}
	defer tx.Rollback()

	err = tx.QueryRow(`INSERT INTO users (email, password, role) VALUES ($1, $2, 'mahasiswa') RETURNING id`,
		email, passHash).Scan(&userID)
	if err != nil {
		return 0, 0, err
	}
	err = tx.QueryRow(`INSERT INTO students (user_id, nim, nama, prodi, angkatan, ipk_terakhir)
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
		userID, nim, nama, prodi, angkatan, ipk).Scan(&studentID)
	if err != nil {
		return 0, 0, err
	}
	return userID, studentID, tx.Commit()
}

func UpdateStudent(db *sql.DB, id int, nama, prodi string, angkatan int, ipk float64) error {
	res, err := db.Exec(`UPDATE students SET nama=$1, prodi=$2, angkatan=$3, ipk_terakhir=$4
		WHERE id=$5 AND deleted_at IS NULL`, nama, prodi, angkatan, ipk, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func SoftDeleteStudent(db *sql.DB, id int) error {
	res, err := db.Exec(`UPDATE students SET deleted_at = NOW() WHERE id=$1 AND deleted_at IS NULL`, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func NIMExists(db *sql.DB, nim string) bool {
	var c int
	_ = db.QueryRow(`SELECT COUNT(*) FROM students WHERE nim=$1`, nim).Scan(&c)
	return c > 0
}

func EmailExists(db *sql.DB, email string) bool {
	var c int
	_ = db.QueryRow(`SELECT COUNT(*) FROM users WHERE email=$1`, email).Scan(&c)
	return c > 0
}
