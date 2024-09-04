package handler

import (
	"errors"
	"file-server/internal/app/api/interface/request"
	"file-server/internal/app/api/interface/response"
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
	var folder request.CreateFolderRequest
	if err := c.ShouldBindJSON(&folder); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	folderDTO, err := fh.usecase.Create(folder.ParentFolderID, folder.Name, folder.IsHide)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, err.Error())
		} else {
			c.JSON(http.StatusInternalServerError, err.Error())
		}
		return
	}

	c.JSON(http.StatusOK, fh.dtoToResponse(folderDTO))
}

func (fh *folderHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	var folder request.UpdateFolderRequest
	if err := c.ShouldBindJSON(&folder); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	folderDTO, err := fh.usecase.Update(id, folder.Name, folder.IsHide)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, err.Error())
		} else {
			c.JSON(http.StatusInternalServerError, err.Error())
		}
		return
	}

	c.JSON(http.StatusOK, fh.dtoToResponse(folderDTO))
}

func (fh *folderHandler) Remove(c *gin.Context) {
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

func (fh *folderHandler) Move(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	var folder request.MoveFolderRequest
	if err := c.ShouldBindJSON(&folder); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	folderDTO, err := fh.usecase.Move(id, folder.ParentFolderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, err.Error())
		} else {
			c.JSON(http.StatusInternalServerError, err.Error())
		}
		return
	}

	c.JSON(http.StatusOK, fh.dtoToResponse(folderDTO))
}

func (fh *folderHandler) Copy(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	var folder request.CopyFolderRequest
	if err := c.ShouldBindJSON(&folder); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	folderDTO, err := fh.usecase.Copy(id, folder.ParentFolderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, err.Error())
		} else {
			c.JSON(http.StatusInternalServerError, err.Error())
		}
		return
	}

	c.JSON(http.StatusOK, fh.dtoToResponse(folderDTO))
}

func (fh *folderHandler) FindOne(c *gin.Context) {
	path := c.Param("path")

	folderDTO, err := fh.usecase.FindOne(path)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, err.Error())
		} else {
			c.JSON(http.StatusInternalServerError, err.Error())
		}
		return
	}

	c.JSON(http.StatusOK, fh.dtoToResponse(folderDTO))
}

func (fh *folderHandler) dtoToResponse(folder *dto.FolderDTO) *response.FolderResponse {
	var folders []response.FolderResponse
	if folder.Folders != nil {
		folders = make([]response.FolderResponse, len(folder.Folders))
		for i, v := range folder.Folders {
			folders[i] = *fh.dtoToResponse(&v)
		}
	}
	var files []response.FileResponse
	if folder.Files != nil {
		files = make([]response.FileResponse, len(folder.Files))
		for i, v := range folder.Files {
			files[i] = response.FileResponse{
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
	return &response.FolderResponse{
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
