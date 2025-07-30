package config

import (
	"nyx-uploader/uploader/utils" // Sesuaikan path
	"time"
)

// Konfigurasi aplikasi yang dibaca dari environment variables.
var (
	Port          = utils.GetEnv("APP_PORT", ":8880")
	BaseURL       = utils.GetEnv("APP_BASE_URL", "http://localhost"+Port)
	UploadDir     = utils.GetEnv("UPLOAD_DIR", "uploads/")
	MaxUploadSize = int64(utils.GetEnvAsInt("MAX_UPLOAD_SIZE_MB", 300) * 1024 * 1024)
	ReadTimeout   = time.Duration(utils.GetEnvAsInt("READ_TIMEOUT", 15)) * time.Second
	WriteTimeout  = time.Duration(utils.GetEnvAsInt("WRITE_TIMEOUT", 15)) * time.Second
	IdleTimeout   = time.Duration(utils.GetEnvAsInt("IDLE_TIMEOUT", 60)) * time.Second
)

// AllowedMimeTypes daftar tipe file yang diizinkan.
var AllowedMimeTypes = map[string]bool{
	// Images
	"image/jpeg": true, "image/png": true, "image/gif": true, "image/webp": true, "image/svg+xml": true, "image/bmp": true, "image/tiff": true,
	// Documents
	"application/pdf": true, "text/plain": true, "text/csv": true, "application/rtf": true,
	"application/msword": true, "application/vnd.openxmlformats-officedocument.wordprocessingml.document":   true, // .docx
	"application/vnd.ms-excel": true, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":       true, // .xlsx
	"application/vnd.ms-powerpoint": true, "application/vnd.openxmlformats-officedocument.presentationml.presentation": true, // .pptx
	// Audio
	"audio/mpeg": true, "audio/wav": true, "audio/ogg": true, "audio/aac": true, "audio/flac": true,
	// Video
	"video/mp4": true, "video/webm": true, "video/ogg": true, "video/x-msvideo": true, // .avi
	"video/quicktime": true, // .mov
	// Archives
	"application/zip": true, "application/x-rar-compressed": true, "application/x-7z-compressed": true, "application/x-tar": true,
}

