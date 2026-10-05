package helper

import "golang.org/x/crypto/bcrypt"

// HashPassword untuk menyimpan password (dipakai saat buat user / seeder).
func HashPassword(plain string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	return string(b), err
}

// CheckPassword untuk verifikasi saat login.
func CheckPassword(plain, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain)) == nil
}

// BatasSKS menentukan batas SKS dari IPK terakhir.
func BatasSKS(ipk float64) int {
	if ipk >= 3.00 {
		return 24
	}
	if ipk >= 2.50 {
		return 21
	}
	return 18
}
