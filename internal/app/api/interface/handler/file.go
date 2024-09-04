package handler

import (
	"errors"
	"file-server/internal/app/api/interface/request"
	"file-server/internal/app/api/interface/response"
	"file-server/internal/app/api/usecase"
	"file-server/internal/app/api/usecase/dto"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type FileHandler interface {
	Create(*gin.Context)
	Update(*gin.Context)
	Remove(*gin.Context)
	Move(*gin.Context)
	Copy(*gin.Context)
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

	fileDTO, err := fh.usecase.Create(file.FolderID, header.Filename, file.IsHide, body)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, err.Error())
		} else {
			c.JSON(http.StatusInternalServerError, err.Error())
		}
		return
	}

	c.JSON(http.StatusOK, fh.dtoToResponse(fileDTO))
}

func (fh *fileHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	var file request.UpdateFileRequest
	if err := c.ShouldBindJSON(&file); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	fileDTO, err := fh.usecase.Update(id, file.Name, file.IsHide)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, err.Error())
		} else {
			c.JSON(http.StatusInternalServerError, err.Error())
		}
		return
	}

	c.JSON(http.StatusOK, fh.dtoToResponse(fileDTO))
}

func (fh *fileHandler) Remove(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if err := fh.usecase.Remove(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, err.Error())
		} else {
			c.JSON(http.StatusInternalServerError, err.Error())
		}
		return
	}

	c.Status(http.StatusNoContent)
}

func (fh *fileHandler) Move(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	var file request.MoveFileRequest
	if err := c.ShouldBindJSON(&file); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	fileDTO, err := fh.usecase.Move(id, file.FolderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, err.Error())
		} else {
			c.JSON(http.StatusInternalServerError, err.Error())
		}
		return
	}

	c.JSON(http.StatusOK, fh.dtoToResponse(fileDTO))
}

func (fh *fileHandler) Copy(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	var file request.CopyFileRequest
	if err := c.ShouldBindJSON(&file); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	fileDTO, err := fh.usecase.Copy(id, file.FolderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, err.Error())
		} else {
			c.JSON(http.StatusInternalServerError, err.Error())
		}
		return
	}

	c.JSON(http.StatusOK, fh.dtoToResponse(fileDTO))
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
