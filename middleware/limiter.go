package middleware

import (
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/patrickmn/go-cache"
	"golang.org/x/time/rate"
)

// VisitorCache menyimpan limiter untuk setiap IP.
// Menggunakan cache dengan expiration untuk mencegah memory leak.
var visitorCache = cache.New(3*time.Minute, 5*time.Minute)

func getVisitorLimiter(ip string) *rate.Limiter {
	limiter, found := visitorCache.Get(ip)
	if !found {
		// Izinkan 5 request per detik, dengan burst hingga 10 request
		newLimiter := rate.NewLimiter(5, 10)
		visitorCache.Set(ip, newLimiter, cache.DefaultExpiration)
		return newLimiter
	}
	return limiter.(*rate.Limiter)
}

// RateLimit middleware membatasi jumlah request per IP.
func RateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			slog.Warn("Gagal membaca remote address", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		limiter := getVisitorLimiter(ip)
		if !limiter.Allow() {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

