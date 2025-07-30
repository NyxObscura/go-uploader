package handler

import (
	"crypto/rand"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"nyx-uploader/uploader/config" // Sesuaikan path
	"nyx-uploader/uploader/utils"  // Sesuaikan path
)

// UploadFile menangani logika unggah file untuk publik.
func UploadFile(w http.ResponseWriter, r *http.Request) {
	// Logika pengecekan API Key telah dihapus.

	// 1. Batasi ukuran request body
	r.Body = http.MaxBytesReader(w, r.Body, config.MaxUploadSize)
	if err := r.ParseMultipartForm(config.MaxUploadSize); err != nil {
		msg := fmt.Sprintf("File is too large. Max size is %d MB.", config.MaxUploadSize/1024/1024)
		if err.Error() != "http: request body too large" {
			msg = "Invalid request."
		}
		utils.WriteJSON(w, http.StatusRequestEntityTooLarge, utils.JSONResponse{Success: false, Message: msg})
		return
	}

	// 2. Ambil file dari form (gunakan nama 'file')
	file, handler, err := r.FormFile("file")
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.JSONResponse{Success: false, Message: "Failed to get file from form. Ensure the field name is 'file'."})
		return
	}
	defer file.Close()

	// 3. Validasi Tipe MIME
	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, utils.JSONResponse{Success: false, Message: "Failed to read file header."})
		return
	}
	file.Seek(0, 0)

	mimeType := http.DetectContentType(buffer)
	if !config.AllowedMimeTypes[mimeType] {
		msg := fmt.Sprintf("File type '%s' is not allowed.", mimeType)
		utils.WriteJSON(w, http.StatusUnsupportedMediaType, utils.JSONResponse{Success: false, Message: msg})
		return
	}
	slog.Info("File validated", slog.String("filename", handler.Filename), slog.String("mime_type", mimeType))

	// 4. Buat nama file unik
	safeOriginalFilename := filepath.Base(handler.Filename)
	ext := filepath.Ext(safeOriginalFilename)
	randBytes := make([]byte, 8)
	rand.Read(randBytes)
	newFileName := fmt.Sprintf("%d-%x%s", time.Now().UnixNano(), randBytes, ext)
	targetPath := filepath.Join(config.UploadDir, newFileName)
	
	// 5. Pastikan direktori uploads ada
	if err := os.MkdirAll(config.UploadDir, os.ModePerm); err != nil {
		slog.Error("Failed to create upload directory", "error", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.JSONResponse{Success: false, Message: "Server error while processing file."})
		return
	}
	
	// 6. Simpan file
	dst, err := os.Create(targetPath)
	if err != nil {
		slog.Error("Failed to create file on server", "error", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.JSONResponse{Success: false, Message: "Failed to save file on server."})
		return
	}
	defer dst.Close()
	
	if _, err := io.Copy(dst, file); err != nil {
		slog.Error("Failed to copy file content", "error", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.JSONResponse{Success: false, Message: "Failed to save file content."})
		return
	}
	
	// 7. Kirim respons sukses
	fullPath := fmt.Sprintf("%s/uploads/%s", config.BaseURL, newFileName)
	utils.WriteJSON(w, http.StatusOK, utils.JSONResponse{
		Success: true,
		Message: "File uploaded successfully!",
		Data:    map[string]string{"path": fullPath},
	})
}

