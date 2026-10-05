package repository

import (
	"database/sql"

	"siakad-mini/app/model"
)

// Cari user dari email (untuk login).
// Mahasiswa yang sudah soft-delete dianggap tidak ada (tidak dapat login).
func FindUserByEmail(db *sql.DB, email string) (*model.User, *model.Student, error) {
	var u model.User
	err := db.QueryRow(
		`SELECT id, email, password, role FROM users WHERE email = $1`, email,
	).Scan(&u.ID, &u.Email, &u.Password, &u.Role)
	if err != nil {
		return nil, nil, err
	}
	if u.Role == "mahasiswa" {
		var s model.Student
		err := db.QueryRow(
			`SELECT id, user_id, nim, nama, prodi, angkatan, ipk_terakhir
			 FROM students WHERE user_id = $1 AND deleted_at IS NULL`, u.ID,
		).Scan(&s.ID, &s.UserID, &s.NIM, &s.Nama, &s.Prodi, &s.Angkatan, &s.IPKTerakhir)
		if err != nil {
			return nil, nil, err // dianggap tidak ditemukan -> 401
		}
		return &u, &s, nil
	}
	return &u, nil, nil
}

func FindUserByID(db *sql.DB, id int) (*model.User, error) {
	var u model.User
	err := db.QueryRow(
		`SELECT id, email, role FROM users WHERE id = $1`, id,
	).Scan(&u.ID, &u.Email, &u.Role)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func FindStudentByUserID(db *sql.DB, userID int) (*model.Student, error) {
	var s model.Student
	err := db.QueryRow(
		`SELECT id, user_id, nim, nama, prodi, angkatan, ipk_terakhir
		 FROM students WHERE user_id = $1 AND deleted_at IS NULL`, userID,
	).Scan(&s.ID, &s.UserID, &s.NIM, &s.Nama, &s.Prodi, &s.Angkatan, &s.IPKTerakhir)
	if err != nil {
		return nil, err
	}
	return &s, nil
}
