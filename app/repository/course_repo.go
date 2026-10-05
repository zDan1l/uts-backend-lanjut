package repository

import (
	"database/sql"
	"fmt"

	"siakad-mini/app/model"
)

type CourseFilter struct {
	Semester  int // 0 = semua
	Search    string
	Available bool
}

// Daftar MK + terisi & sisa kuota (dihitung dari enrollments)
func ListCourses(db *sql.DB, f CourseFilter) ([]model.Course, error) {
	q := `SELECT c.id, c.kode_mk, c.nama_mk, c.sks, c.semester, c.kuota,
		COALESCE((SELECT COUNT(*) FROM enrollments e WHERE e.course_id = c.id), 0) AS terisi
		FROM courses c WHERE 1=1`
	args := []interface{}{}
	n := 1

	if f.Semester > 0 {
		q += fmt.Sprintf(" AND c.semester = $%d", n)
		args = append(args, f.Semester)
		n++
	}
	if f.Search != "" {
		q += fmt.Sprintf(" AND (c.kode_mk ILIKE $%d OR c.nama_mk ILIKE $%d)", n, n+1)
		like := "%" + f.Search + "%"
		args = append(args, like, like)
		n += 2
	}
	q += " ORDER BY c.id ASC"

	rows, err := db.Query(q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := []model.Course{}
	for rows.Next() {
		var c model.Course
		if err := rows.Scan(&c.ID, &c.KodeMK, &c.NamaMK, &c.SKS, &c.Semester, &c.Kuota, &c.Terisi); err != nil {
			return nil, err
		}
		c.SisaKuota = c.Kuota - c.Terisi
		if c.SisaKuota < 0 {
			c.SisaKuota = 0
		}
		if f.Available && c.SisaKuota <= 0 {
			continue
		}
		out = append(out, c)
	}
	if out == nil {
		out = []model.Course{}
	}
	return out, nil
}
