package handler

import (
	"errors"
	"file-server/internal/app/api/interface/requests"
	"file-server/internal/app/api/interface/responses"
	"file-server/internal/app/api/usecase"
	"file-server/internal/app/api/usecase/dto"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type FolderHandler interface {
	Create(*gin.Context)
	Update(*gin.Context)
	Remove(*gin.Context)
	Move(*gin.Context)
	Copy(*gin.Context)
	FindOne(*gin.Context)
}

type folderHandler struct {
	usecase usecase.FolderUsecase
}

func NewFolderHandler(usecase usecase.FolderUsecase) FolderHandler {
	return &folderHandler{
		usecase: usecase,
	}
}

func (fh *folderHandler) Create(c *gin.Context) {
	var request requests.CreateFolderRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	dto, err := fh.usecase.Create(request.ParentFolderID, request.Name, request.IsHide)
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

func (fh *folderHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	var request requests.UpdateFolderRequest
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

func (fh *folderHandler) Remove(c *gin.Context) {
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

func (fh *folderHandler) Move(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	var request requests.MoveFolderRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	dto, err := fh.usecase.Move(id, request.ParentFolderID, fh.getIsDisplayHiddenObject(c))
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

func (fh *folderHandler) Copy(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	var request requests.CopyFolderRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	dto, err := fh.usecase.Copy(id, request.ParentFolderID, fh.getIsDisplayHiddenObject(c))
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

func (fh *folderHandler) FindOne(c *gin.Context) {
	path := c.Param("path")

	dto, err := fh.usecase.FindOne(path, fh.getIsDisplayHiddenObject(c))
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

func (fh *folderHandler) getIsDisplayHiddenObject(c *gin.Context) bool {
	if v, ok := c.Get("isDisplayHiddenObject"); !ok || v == false {
		return false
	} else {
		return true
	}
}

func (fh *folderHandler) dtoToResponse(folder *dto.FolderDTO) *responses.FolderResponse {
	var folders []responses.FolderResponse
	if folder.Folders != nil {
		folders = make([]responses.FolderResponse, len(folder.Folders))
		for i, v := range folder.Folders {
			folders[i] = *fh.dtoToResponse(&v)
		}
	}
	var files []responses.FileResponse
	if folder.Files != nil {
		files = make([]responses.FileResponse, len(folder.Files))
		for i, v := range folder.Files {
			files[i] = responses.FileResponse{
				ID:        v.ID,
				FolderID:  v.FolderID,
				Name:      v.Name,
				Path:      v.Path,
				MimeType:  v.MimeType,
				IsHide:    v.IsHide,
				CreatedAt: v.CreatedAt,
				UpdatedAt: v.UpdatedAt,
			}
		}
	}
	return &responses.FolderResponse{
		ID:             folder.ID,
		ParentFolderID: folder.ParentFolderID,
		Name:           folder.Name,
		Path:           folder.Path,
		IsHide:         folder.IsHide,
		Folders:        folders,
		Files:          files,
		CreatedAt:      folder.CreatedAt,
		UpdatedAt:      folder.UpdatedAt,
	}
}
