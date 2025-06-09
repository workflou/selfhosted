package file

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

func UploadFromRequest(r *http.Request, fieldName string, uploadsPath string) (string, error) {
	file, fileHeader, err := r.FormFile(fieldName)
	if err != nil {
		return "", err
	}
	defer file.Close()

	err = os.MkdirAll(uploadsPath, os.ModePerm)
	if err != nil {
		return "", fmt.Errorf("failed to create upload directory: %w", err)
	}

	filename := uuid.New().String() + filepath.Ext(fileHeader.Filename)
	path := filepath.Join(uploadsPath, filename)

	dst, err := os.Create(path)
	if err != nil {
		return "", fmt.Errorf("failed to create upload file: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return "", fmt.Errorf("failed to copy uploaded file: %w", err)
	}

	return path, nil
}
