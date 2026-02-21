package handler

import (
	"net/http"
	"path/filepath"

	common "github.com/ZyoGo/ayo-indonesia-footbal/pkg/http"
	"github.com/ZyoGo/ayo-indonesia-footbal/pkg/upload"
	"github.com/gin-gonic/gin"
)

type UploadHandler struct {
	uploader *upload.Uploader
}

func NewUploadHandler(uploader *upload.Uploader) *UploadHandler {
	return &UploadHandler{uploader: uploader}
}

type UploadResponse struct {
	URL      string `json:"url"`
	Type     string `json:"type"`
	Filename string `json:"filename"`
}

func (h *UploadHandler) Upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
		return
	}

	uploadType := upload.UploadType(c.PostForm("type"))
	if !uploadType.IsValid() {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid upload type",
			"data":    "type must be one of: team-logo, player-photo, document",
		})
		return
	}

	url, err := h.uploader.UploadWithType(file, uploadType)
	if err != nil {
		resp := common.RenderErrorResponse(err)
		c.JSON(resp.Code, resp)
		return
	}

	filename := filepath.Base(url)
	c.JSON(http.StatusOK, common.NewSuccessResponseWithData(UploadResponse{
		URL:      url,
		Type:     string(uploadType),
		Filename: filename,
	}))
}
