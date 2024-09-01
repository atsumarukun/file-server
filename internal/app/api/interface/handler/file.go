package handler

import (
	"file-server/internal/app/api/interface/request"
	"file-server/internal/app/api/interface/response"
	"file-server/internal/app/api/usecase"
	"file-server/internal/app/api/usecase/dto"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type FileHandler interface {
	Create(*gin.Context)
	Update(*gin.Context)
	Remove(*gin.Context)
}

type fileHandler struct {
	usecase usecase.FileUsecase
}

func NewFileHandler(usecase usecase.FileUsecase) FileHandler {
	return &fileHandler{
		usecase: usecase,
	}
}

func (fh *fileHandler) Create(c *gin.Context) {
	f, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	defer f.Close()
	body, err := io.ReadAll(f)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	var file request.CreateFileRequest
	if err := c.Bind(&file); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	fileDTO, apiErr := fh.usecase.Create(file.FolderID, header.Filename, file.IsHide, body)
	if apiErr != nil {
		c.JSON(apiErr.Code, apiErr.Message)
		return
	}

	c.JSON(http.StatusOK, fh.dtoToResponse(fileDTO))
}

func (fh *fileHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	var file request.UpdateFileRequest
	if err := c.ShouldBindJSON(&file); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	fileDTO, apiErr := fh.usecase.Update(id, file.Name, file.IsHide)
	if apiErr != nil {
		c.JSON(apiErr.Code, apiErr.Message)
		return
	}

	c.JSON(http.StatusOK, fh.dtoToResponse(fileDTO))
}

func (fh *fileHandler) Remove(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if apiErr := fh.usecase.Remove(id); apiErr != nil {
		c.JSON(apiErr.Code, apiErr.Message)
		return
	}

	c.Status(http.StatusNoContent)
}

func (fh *fileHandler) dtoToResponse(file *dto.FileDTO) *response.FileResponse {
	return &response.FileResponse{
		ID:        file.ID,
		FolderID:  file.FolderID,
		Name:      file.Name,
		Path:      file.Path,
		MimeType:  file.MimeType,
		IsHide:    file.IsHide,
		CreatedAt: file.CreatedAt,
		UpdatedAt: file.UpdatedAt,
	}
}
