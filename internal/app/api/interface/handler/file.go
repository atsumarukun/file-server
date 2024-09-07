package handler

import (
	"errors"
	"file-server/internal/app/api/interface/requests"
	"file-server/internal/app/api/interface/responses"
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

	var request requests.CreateFileRequest
	if err := c.Bind(&request); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	dto, err := fh.usecase.Create(request.FolderID, header.Filename, request.IsHide, body)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, err.Error())
		} else {
			c.JSON(http.StatusInternalServerError, err.Error())
		}
		return
	}

	c.JSON(http.StatusOK, fh.dtoToResponse(dto))
}

func (fh *fileHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	var request requests.UpdateFileRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	dto, err := fh.usecase.Update(id, request.Name, request.IsHide, fh.getIsDisplayHiddenObject(c))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, err.Error())
		} else {
			c.JSON(http.StatusInternalServerError, err.Error())
		}
		return
	}

	c.JSON(http.StatusOK, fh.dtoToResponse(dto))
}

func (fh *fileHandler) Remove(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if err := fh.usecase.Remove(id, fh.getIsDisplayHiddenObject(c)); err != nil {
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

	var request requests.MoveFileRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	dto, err := fh.usecase.Move(id, request.FolderID, fh.getIsDisplayHiddenObject(c))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, err.Error())
		} else {
			c.JSON(http.StatusInternalServerError, err.Error())
		}
		return
	}

	c.JSON(http.StatusOK, fh.dtoToResponse(dto))
}

func (fh *fileHandler) Copy(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	var request requests.CopyFileRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	dto, err := fh.usecase.Copy(id, request.FolderID, fh.getIsDisplayHiddenObject(c))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, err.Error())
		} else {
			c.JSON(http.StatusInternalServerError, err.Error())
		}
		return
	}

	c.JSON(http.StatusOK, fh.dtoToResponse(dto))
}

func (fh *fileHandler) getIsDisplayHiddenObject(c *gin.Context) bool {
	if v, ok := c.Get("isDisplayHiddenObject"); !ok || v == false {
		return false
	} else {
		return true
	}
}

func (fh *fileHandler) dtoToResponse(file *dto.FileDTO) *responses.FileResponse {
	return &responses.FileResponse{
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
