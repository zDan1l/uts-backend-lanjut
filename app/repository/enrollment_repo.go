package repository

import (
	"database/sql"
	"fmt"

	"siakad-mini/app/model"
	"siakad-mini/helper"
)

// Ambil 1 enrollment (untuk cek kepemilikan saat delete)
func FindEnrollmentByID(db *sql.DB, id int) (*model.Enrollment, error) {
	var e model.Enrollment
	err := db.QueryRow(`SELECT id, student_id, course_id, tahun_akademik FROM enrollments WHERE id=$1`, id).
		Scan(&e.ID, &e.StudentID, &e.CourseID, &e.TahunAkademik)
	if err != nil {
		return nil, err
	}
	return &e, nil
}

// Tambah KRS dalam 1 transaction + row locking.
// Cek: duplikasi (409), kuota penuh (422), batas SKS (422).
type EnrollError struct {
	Status  int
	Message string
	Errors  map[string][]string
}

func (e *EnrollError) Error() string { return e.Message }

func CreateEnrollmentTx(db *sql.DB, studentID, courseID int, tahun string) (int, error) {
	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	// Kunci baris course agar kuota aman dari race
	var sks, kuota int
	var kode string
	err = tx.QueryRow(`SELECT sks, kuota, kode_mk FROM courses WHERE id=$1 FOR UPDATE`, courseID).
		Scan(&sks, &kuota, &kode)
	if err == sql.ErrNoRows {
		return 0, &EnrollError{Status: 422, Message: "Validasi gagal",
			Errors: map[string][]string{"course_id": {"Mata kuliah tidak ditemukan"}}}
	}
	if err != nil {
		return 0, err
	}

	// Duplikasi: MK sama di tahun akademik sama
	var dup int
	_ = tx.QueryRow(`SELECT COUNT(*) FROM enrollments WHERE student_id=$1 AND course_id=$2 AND tahun_akademik=$3`,
		studentID, courseID, tahun).Scan(&dup)
	if dup > 0 {
		return 0, &EnrollError{Status: 409, Message: "Mata kuliah sudah diambil pada tahun akademik ini"}
	}

	// Kuota
	var terisi int
	_ = tx.QueryRow(`SELECT COUNT(*) FROM enrollments WHERE course_id=$1`, courseID).Scan(&terisi)
	if terisi >= kuota {
		return 0, &EnrollError{Status: 422, Message: "Validasi gagal",
			Errors: map[string][]string{"course_id": {"Kuota mata kuliah penuh"}}}
	}

	// Batas SKS: total SKS student di tahun tsb + sks baru
	var ipk float64
	if err := tx.QueryRow(`SELECT ipk_terakhir FROM students WHERE id=$1 AND deleted_at IS NULL`, studentID).Scan(&ipk); err != nil {
		return 0, err
	}
	batas := helper.BatasSKS(ipk)
	var total int
	_ = tx.QueryRow(`SELECT COALESCE(SUM(c.sks),0) FROM enrollments e JOIN courses c ON c.id=e.course_id
		WHERE e.student_id=$1 AND e.tahun_akademik=$2`, studentID, tahun).Scan(&total)
	if total+sks > batas {
		sisa := batas - total
		if sisa < 0 {
			sisa = 0
		}
		return 0, &EnrollError{Status: 422, Message: "Validasi gagal",
			Errors: map[string][]string{
				"sks": {fmt.Sprintf("Total SKS melebihi batas (%d SKS). Sisa SKS Anda: %d SKS", batas, sisa)},
			}}
	}

	var newID int
	err = tx.QueryRow(`INSERT INTO enrollments (student_id, course_id, tahun_akademik)
		VALUES ($1,$2,$3) RETURNING id`, studentID, courseID, tahun).Scan(&newID)
	if err != nil {
		return 0, err
	}
	return newID, tx.Commit()
}

func DeleteEnrollment(db *sql.DB, id int) error {
	_, err := db.Exec(`DELETE FROM enrollments WHERE id=$1`, id)
	return err
}
