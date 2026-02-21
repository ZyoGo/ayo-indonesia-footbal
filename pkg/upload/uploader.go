package upload

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ZyoGo/ayo-indonesia-footbal/pkg/derrors"
	"github.com/oklog/ulid/v2"
)

// Uploader handles file uploads
type Uploader struct {
	BasePath     string
	MaxSize      int64
	AllowedTypes []string
	URLPrefix    string
}

// NewUploader creates a new file uploader
func NewUploader(basePath, urlPrefix string, maxSize int64, allowedTypes []string) *Uploader {
	return &Uploader{
		BasePath:     basePath,
		MaxSize:      maxSize,
		AllowedTypes: allowedTypes,
		URLPrefix:    urlPrefix,
	}
}

// UploadFile uploads a file and returns the file URL
func (u *Uploader) UploadFile(fileHeader *multipart.FileHeader, subdirectory string) (string, error) {
	// Validate file size
	if fileHeader.Size > u.MaxSize {
		return "", derrors.NewErrorf(derrors.ErrorCodeInvalidArgument,
			"file size exceeds maximum allowed size of %d bytes", u.MaxSize)
	}

	// Validate file type
	contentType := fileHeader.Header.Get("Content-Type")
	if !u.isAllowedType(contentType) {
		return "", derrors.NewErrorf(derrors.ErrorCodeInvalidArgument,
			"file type %s is not allowed", contentType)
	}

	// Create upload directory
	uploadDir := filepath.Join(u.BasePath, subdirectory)
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return "", derrors.WrapErrorf(err, derrors.ErrorCodeInternal, "failed to create upload directory")
	}

	// Generate unique filename
	ext := filepath.Ext(fileHeader.Filename)
	filename := fmt.Sprintf("%s%s", ulid.Make().String(), ext)
	filePath := filepath.Join(uploadDir, filename)

	// Open uploaded file
	src, err := fileHeader.Open()
	if err != nil {
		return "", derrors.WrapErrorf(err, derrors.ErrorCodeInternal, "failed to open uploaded file")
	}
	defer src.Close()

	// Create destination file
	dst, err := os.Create(filePath)
	if err != nil {
		return "", derrors.WrapErrorf(err, derrors.ErrorCodeInternal, "failed to create destination file")
	}
	defer dst.Close()

	// Copy file content
	if _, err := io.Copy(dst, src); err != nil {
		return "", derrors.WrapErrorf(err, derrors.ErrorCodeInternal, "failed to save file")
	}

	// Generate URL
	fileURL := fmt.Sprintf("%s/%s/%s", u.URLPrefix, subdirectory, filename)

	return fileURL, nil
}

// UploadWithType uploads a file with type-based validation and returns the file URL
func (u *Uploader) UploadWithType(fileHeader *multipart.FileHeader, uploadType UploadType) (string, error) {
	if !uploadType.IsValid() {
		return "", ErrInvalidUploadType
	}

	if err := ValidateFile(uploadType, fileHeader); err != nil {
		return "", err
	}

	config := uploadType.Config()
	subdirectory := string(uploadType)

	uploadDir := filepath.Join(u.BasePath, subdirectory)
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return "", derrors.WrapErrorf(err, derrors.ErrorCodeInternal, "failed to create upload directory")
	}

	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if !isAllowedExtension(ext, config.AllowedExtensions) {
		return "", ErrFileTypeNotAllowed
	}

	filename := fmt.Sprintf("%s%s", ulid.Make().String(), ext)
	filePath := filepath.Join(uploadDir, filename)

	src, err := fileHeader.Open()
	if err != nil {
		return "", derrors.WrapErrorf(err, derrors.ErrorCodeInternal, "failed to open uploaded file")
	}
	defer src.Close()

	dst, err := os.Create(filePath)
	if err != nil {
		return "", derrors.WrapErrorf(err, derrors.ErrorCodeInternal, "failed to create destination file")
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return "", derrors.WrapErrorf(err, derrors.ErrorCodeInternal, "failed to save file")
	}

	fileURL := fmt.Sprintf("%s/%s/%s", u.URLPrefix, subdirectory, filename)

	return fileURL, nil
}

func isAllowedExtension(ext string, allowed []string) bool {
	for _, a := range allowed {
		if a == ext {
			return true
		}
	}
	return false
}

// DeleteFile deletes a file by URL
func (u *Uploader) DeleteFile(fileURL string) error {
	// Extract filename from URL
	parts := strings.Split(fileURL, "/")
	if len(parts) < 2 {
		return nil // Invalid URL, ignore
	}

	// Reconstruct file path
	relativePath := strings.Join(parts[len(parts)-2:], "/")
	filePath := filepath.Join(u.BasePath, relativePath)

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil // File doesn't exist, nothing to delete
	}

	// Delete file
	if err := os.Remove(filePath); err != nil {
		return derrors.WrapErrorf(err, derrors.ErrorCodeInternal, "failed to delete file")
	}

	return nil
}

// isAllowedType checks if the content type is allowed
func (u *Uploader) isAllowedType(contentType string) bool {
	for _, allowed := range u.AllowedTypes {
		if contentType == allowed {
			return true
		}
	}
	return false
}

// GenerateFilename generates a unique filename with timestamp
func GenerateFilename(originalName string) string {
	ext := filepath.Ext(originalName)
	timestamp := time.Now().UnixNano()
	return fmt.Sprintf("%d_%s%s", timestamp, ulid.Make().String(), ext)
}
