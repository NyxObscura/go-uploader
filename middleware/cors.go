package middleware

import "github.com/rs/cors"

// CorsHandler mengkonfigurasi dan menyediakan middleware CORS.
func CorsHandler() *cors.Cors {
	return cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // Lebih aman: ganti dengan domain frontend Anda
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		Debug:            false, // Set true untuk debugging
	})
}

