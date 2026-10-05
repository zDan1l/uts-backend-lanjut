package database

import (
	"database/sql"
	"fmt"

	"siakad-mini/helper"
)

// Seed mengisi 1 admin, 10 MK, 20 mahasiswa.
// Dilewati jika tabel users sudah ada isinya.
func Seed(db *sql.DB) error {
	var n int
	if err := db.QueryRow(`SELECT COUNT(*) FROM users`).Scan(&n); err != nil {
		return err
	}
	if n > 0 {
		fmt.Println("seeder dilewati: data sudah ada")
		return nil
	}

	adminPass, _ := helper.HashPassword("password123")
	if _, err := db.Exec(`INSERT INTO users (email, password, role) VALUES ($1,$2,'admin')`,
		"admin@siakad.ac.id", adminPass); err != nil {
		return fmt.Errorf("seed admin: %w", err)
	}

	courses := []struct {
		kode, nama       string
		sks, smt, kuota  int
	}{
		{"MK001", "Algoritma dan Pemrograman", 3, 1, 40},
		{"MK002", "Basis Data", 3, 2, 40},
		{"MK003", "Pemrograman Web", 3, 3, 30},
		{"MK004", "Pemrograman Backend Lanjut", 3, 5, 30},
		{"MK005", "Sistem Operasi", 3, 3, 35},
		{"MK006", "Jaringan Komputer", 2, 4, 35},
		{"MK007", "Rekayasa Perangkat Lunak", 3, 4, 30},
		{"MK008", "Kecerdasan Buatan", 3, 6, 25},
		{"MK009", "Statistika", 2, 1, 50},
		{"MK010", "Interaksi Manusia dan Komputer", 2, 3, 40},
	}
	for _, c := range courses {
		if _, err := db.Exec(`INSERT INTO courses (kode_mk, nama_mk, sks, semester, kuota)
			VALUES ($1,$2,$3,$4,$5)`, c.kode, c.nama, c.sks, c.smt, c.kuota); err != nil {
			return fmt.Errorf("seed course: %w", err)
		}
	}

	prodi := []string{"Sistem Informasi", "Informatika", "Teknik Komputer"}
	names := []string{
		"Rina Putri", "Budi Santoso", "Siti Aminah", "Andi Wijaya", "Dewi Lestari",
		"Agus Pratama", "Maya Anggraini", "Rudi Hartono", "Nina Kurnia", "Doni Saputra",
		"Fitri Handayani", "Yoga Nugroho", "Lina Marlina", "Eko Setiawan", "Rani Puspita",
		"Hendra Gunawan", "Tania Rahma", "Fajar Ramadhan", "Intan Permata", "Galih Prasetyo",
	}
	ipks := []float64{3.45, 2.80, 2.30, 3.80, 3.10, 2.60, 1.90, 3.00, 2.49, 3.75,
		2.90, 3.20, 2.10, 3.60, 2.75, 3.90, 2.40, 3.25, 2.95, 3.55}

	for i := 0; i < 20; i++ {
		nim := fmt.Sprintf("187221%06d", i+1) // 12 digit
		email := fmt.Sprintf("mhs%02d@siakad.ac.id", i+1)
		hash, _ := helper.HashPassword(nim)
		var uid int
		err := db.QueryRow(`INSERT INTO users (email, password, role) VALUES ($1,$2,'mahasiswa') RETURNING id`,
			email, hash).Scan(&uid)
		if err != nil {
			return fmt.Errorf("seed user mhs: %w", err)
		}
		if _, err := db.Exec(`INSERT INTO students (user_id, nim, nama, prodi, angkatan, ipk_terakhir)
			VALUES ($1,$2,$3,$4,$5,$6)`,
			uid, nim, names[i], prodi[i%len(prodi)], 2022+(i%3), ipks[i]); err != nil {
			return fmt.Errorf("seed student: %w", err)
		}
	}

	fmt.Println("seeder OK: 1 admin + 20 mahasiswa + 10 MK")
	fmt.Println("admin login -> admin@siakad.ac.id / password123")
	fmt.Println("mhs login   -> mhs01@siakad.ac.id / <nim-nya>, dst.")
	return nil
}
