package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"github.com/gorilla/mux"
	_ "github.com/joho/godotenv/autoload"
	"nyx-uploader/uploader/config"
	"nyx-uploader/uploader/handler"
	"nyx-uploader/uploader/middleware"
)

func main() {
	// Setup structured logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	slog.Info("🚀 Starting Nyx Uploader Server...")

	r := mux.NewRouter()

	// Middleware
	r.Use(middleware.StructuredLogger)
	r.Use(middleware.RateLimit)
	corsHandler := middleware.CorsHandler().Handler(r)

	// 1. Define API routes first
	r.HandleFunc("/upload", handler.UploadFile).Methods("POST")
	
	// 2. Define routes for uploaded files
	uploadsDir := http.StripPrefix("/uploads/", http.FileServer(http.Dir(config.UploadDir)))
	r.PathPrefix("/uploads/").Handler(uploadsDir)
	
	// 3. Use a single File Server for all static assets in the 'public' directory

	staticFileServer := http.FileServer(http.Dir("./public"))
	r.PathPrefix("/").Handler(staticFileServer)
	
	// Server configuration (remains the same)
	srv := &http.Server{
		Addr:         config.Port,
		Handler:      corsHandler,
		ReadTimeout:  config.ReadTimeout,
		WriteTimeout: config.WriteTimeout,
		IdleTimeout:  config.IdleTimeout,
	}

	// Run server & graceful shutdown
	go func() {
		slog.Info("✅ Server running", "url", config.BaseURL)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Failed to start server", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("🛑 Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("Server shutdown failed", "error", err)
	}
	slog.Info("👋 Server successfully shut down.")
}

