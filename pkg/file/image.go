package file

import (
	"image"
	_ "image/gif"  // Register GIF format
	_ "image/jpeg" // Register JPEG format
	_ "image/png"  // Register PNG format
	"log/slog"
	"net/http"
	"path/filepath"
)

func ValidateImageFromRequest(r *http.Request, fieldName string) bool {
	file, fileHeader, err := r.FormFile(fieldName)
	if err != nil {
		return false
	}
	defer file.Close()

	if fileHeader.Size == 0 {
		slog.Error("Uploaded file is empty", "fieldName", fieldName)
		return false
	}

	_, _, err = image.DecodeConfig(file)
	if err != nil {
		slog.Error("Failed to decode image config", "error", err, "fieldName", fieldName)
		return false
	}

	ext := filepath.Ext(fileHeader.Filename)
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".gif" && ext != ".webp" && ext != ".bmp" {
		slog.Error("Unsupported image file extension", "extension", ext, "fieldName", fieldName)
		return false
	}

	contentType := fileHeader.Header.Get("Content-Type")
	slog.Info("Checking content type for image upload", "contentType", contentType)
	switch contentType {
	case "image/jpeg", "image/png", "image/gif", "image/webp", "image/bmp", "image/jpg", "application/octet-stream":
		return true
	default:
		return false
	}
}
