# SIAKAD Mini — RESTful API (Go + Fiber + PostgreSQL)

Struktur sederhana mengikuti contoh:
```
app/
  model/       entitas + struct request-response
  repository/  query manual ke PostgreSQL
  service/     business rules, menerima fiber.Ctx langsung
config/        app.go, env.go, logger.go
database/      postgres.go, migrate.go, seed.go, schema.sql
helper/        response envelope, query string, jwt, hash, batas SKS
logs/          output log (tidak di-commit)
middleware/    global (recover+logger), auth JWT, RequireJSON
route/         pendaftaran 10 endpoint
.env           (tidak di-commit)
main.go
```

## Cara jalan
1. Buat database: `CREATE DATABASE siakad_mini;`
2. Isi `DB_PASSWORD` di `.env`.
3. Seed (migrate + 1 admin, 20 mhs, 10 MK): `go run . --seed`
4. Jalan: `go run .` → http://localhost:3000

## Akun awal
- Admin: `admin@siakad.ac.id` / `password123`
- Mahasiswa: `mhs01@siakad.ac.id` / nim-nya (`187221000001`), dst. sampai mhs20.

## Aturan SKS
IPK ≥3.00 → 24 SKS, 2.50–2.99 → 21 SKS, <2.50 → 18 SKS.
