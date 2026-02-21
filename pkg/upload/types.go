package upload

import (
	"errors"
	"mime/multipart"

	"github.com/ZyoGo/ayo-indonesia-footbal/pkg/derrors"
)

var (
	ErrFileTooLarge       = derrors.NewErrorf(derrors.ErrorCodeInvalidArgument, "file size exceeds maximum allowed size")
	ErrFileTypeNotAllowed = derrors.NewErrorf(derrors.ErrorCodeInvalidArgument, "file type is not allowed")
	ErrInvalidUploadType  = errors.New("invalid upload type")
)

type UploadType string

const (
	TypeTeamLogo    UploadType = "team-logo"
	TypePlayerPhoto UploadType = "player-photo"
	TypeDocument    UploadType = "document"
)

var uploadTypeConfigs = map[UploadType]TypeConfig{
	TypeTeamLogo: {
		AllowedMIME: []string{
			"image/jpeg",
			"image/png",
			"image/gif",
			"image/webp",
		},
		AllowedExtensions: []string{".jpg", ".jpeg", ".png", ".gif", ".webp"},
		MaxSize:           5 * 1024 * 1024, // 5MB
	},
	TypePlayerPhoto: {
		AllowedMIME: []string{
			"image/jpeg",
			"image/png",
			"image/gif",
			"image/webp",
		},
		AllowedExtensions: []string{".jpg", ".jpeg", ".png", ".gif", ".webp"},
		MaxSize:           10 * 1024 * 1024, // 10MB
	},
	TypeDocument: {
		AllowedMIME: []string{
			"application/pdf",
			"application/msword",
			"application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		},
		AllowedExtensions: []string{".pdf", ".doc", ".docx"},
		MaxSize:           25 * 1024 * 1024, // 25MB
	},
}

type TypeConfig struct {
	AllowedMIME       []string
	AllowedExtensions []string
	MaxSize           int64
}

func (t UploadType) IsValid() bool {
	_, exists := uploadTypeConfigs[t]
	return exists
}

func (t UploadType) Config() TypeConfig {
	if config, exists := uploadTypeConfigs[t]; exists {
		return config
	}
	return TypeConfig{
		AllowedMIME:       []string{"image/jpeg", "image/png", "image/gif", "image/webp", "application/pdf"},
		AllowedExtensions: []string{".jpg", ".jpeg", ".png", ".gif", ".webp", ".pdf"},
		MaxSize:           5 * 1024 * 1024,
	}
}

func IsValidUploadType(t string) bool {
	_, exists := uploadTypeConfigs[UploadType(t)]
	return exists
}

func GetAllowedTypes() []string {
	types := make([]string, 0, len(uploadTypeConfigs))
	for t := range uploadTypeConfigs {
		types = append(types, string(t))
	}
	return types
}

func GetAllowedMIMETypes(uploadType UploadType) []string {
	return uploadType.Config().AllowedMIME
}

func ValidateFile(uploadType UploadType, fileHeader *multipart.FileHeader) error {
	config := uploadType.Config()

	if fileHeader.Size > config.MaxSize {
		return ErrFileTooLarge
	}

	contentType := fileHeader.Header.Get("Content-Type")
	if !isAllowedMIME(contentType, config.AllowedMIME) {
		return ErrFileTypeNotAllowed
	}

	return nil
}

func isAllowedMIME(mime string, allowed []string) bool {
	for _, a := range allowed {
		if a == mime {
			return true
		}
	}
	return false
}
