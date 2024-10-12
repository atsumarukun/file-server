package handler

import (
	"errors"
	"file-server/internal/app/api/interface/requests"
	"file-server/internal/app/api/interface/responses"
	"file-server/internal/app/api/usecase"
	"file-server/internal/app/api/usecase/dto"
	"file-server/internal/pkg/types"
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
	Read(*gin.Context)
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
	form, err := c.MultipartForm()
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	files := make([]types.File, len(form.File["files[]"]))
	for i, file := range form.File["files[]"] {
		f, err := file.Open()
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}
		defer f.Close()
		body, err := io.ReadAll(f)
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}
		files[i] = types.File{
			Name: file.Filename,
			Body: body,
		}
	}

	var request requests.CreateFileRequest
	if err := c.Bind(&request); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	dtos, err := fh.usecase.Create(request.FolderID, request.IsHide, files)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.String(http.StatusNotFound, err.Error())
		} else {
			c.String(http.StatusInternalServerError, err.Error())
		}
		return
	}

	res := make([]responses.FileResponse, len(dtos))
	for i, v := range dtos {
		f := fh.convertToFileResponse(&v)
		res[i] = *f
	}

	c.JSON(http.StatusOK, res)
}

func (fh *fileHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	var request requests.UpdateFileRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	dto, err := fh.usecase.Update(id, request.Name, request.IsHide, fh.getIsDisplayHiddenObject(c))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.String(http.StatusNotFound, err.Error())
		} else {
			c.String(http.StatusInternalServerError, err.Error())
		}
		return
	}

	c.JSON(http.StatusOK, fh.convertToFileResponse(dto))
}

func (fh *fileHandler) Remove(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	if err := fh.usecase.Remove(id, fh.getIsDisplayHiddenObject(c)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.String(http.StatusNotFound, err.Error())
		} else {
			c.String(http.StatusInternalServerError, err.Error())
		}
		return
	}

	c.Status(http.StatusNoContent)
}

func (fh *fileHandler) Move(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	var request requests.MoveFileRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	dto, err := fh.usecase.Move(id, request.FolderID, fh.getIsDisplayHiddenObject(c))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.String(http.StatusNotFound, err.Error())
		} else {
			c.String(http.StatusInternalServerError, err.Error())
		}
		return
	}

	c.JSON(http.StatusOK, fh.convertToFileResponse(dto))
}

func (fh *fileHandler) Copy(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	var request requests.CopyFileRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	dto, err := fh.usecase.Copy(id, request.FolderID, fh.getIsDisplayHiddenObject(c))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.String(http.StatusNotFound, err.Error())
		} else {
			c.String(http.StatusInternalServerError, err.Error())
		}
		return
	}

	c.JSON(http.StatusOK, fh.convertToFileResponse(dto))
}

func (fh *fileHandler) Read(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	dto, err := fh.usecase.Read(id, fh.getIsDisplayHiddenObject(c))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.String(http.StatusNotFound, err.Error())
		} else {
			c.String(http.StatusInternalServerError, err.Error())
		}
		return
	}

	c.Data(http.StatusOK, dto.MimeType, dto.Body)
}

func (fh *fileHandler) getIsDisplayHiddenObject(c *gin.Context) bool {
	if v, ok := c.Get("isDisplayHiddenObject"); ok && v == true {
		return true
	} else {
		return false
	}
}

func (fh *fileHandler) convertToFileResponse(file *dto.FileInfoDTO) *responses.FileResponse {
	return responses.NewFileResponse(file.ID, file.FolderID, file.Name, file.Path, file.MimeType, file.IsHide, file.CreatedAt, file.UpdatedAt)
}
