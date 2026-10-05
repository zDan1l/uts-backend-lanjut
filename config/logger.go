package config

import (
	"io"
	"log"
	"os"
)

// SetupLogger menulis log ke logs/app.log sekaligus ke stdout.
// Kembalikan logger + writer (writer dipakai juga oleh HTTP logger).
func SetupLogger() (*log.Logger, io.Writer) {
	_ = os.MkdirAll("logs", 0755)
	f, err := os.OpenFile("logs/app.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		lg := log.New(os.Stdout, "[siakad] ", log.LstdFlags)
		lg.Printf("gagal buka logs/app.log: %v", err)
		return lg, os.Stdout
	}
	w := io.MultiWriter(os.Stdout, f)
	return log.New(w, "[siakad] ", log.LstdFlags), w
}
