package handler

import (
	"file-server/internal/app/api/interface/request"
	"file-server/internal/app/api/interface/response"
	"file-server/internal/app/api/usecase"
	"file-server/internal/app/api/usecase/dto"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type FolderHandler interface {
	Create(*gin.Context)
	Update(*gin.Context)
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

	folderDTO, apiErr := fh.usecase.Create(folder.ParentFolderID, folder.Name, folder.IsHide)
	if apiErr != nil {
		c.JSON(apiErr.Code, apiErr.Message)
		return
	}

	c.JSON(http.StatusOK, fh.dtoToResponse(folderDTO))
}

func (fh *folderHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	var folder request.UpdateFolderRequest
	if err := c.ShouldBindJSON(&folder); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	folderDTO, apiErr := fh.usecase.Update(id, folder.Name, folder.IsHide)
	if apiErr != nil {
		c.JSON(apiErr.Code, apiErr.Message)
		return
	}

	c.JSON(http.StatusOK, fh.dtoToResponse(folderDTO))
}

func (fh *folderHandler) FindOne(c *gin.Context) {
	path := c.Param("path")

	folderDTO, apiErr := fh.usecase.FindOne(path)
	if apiErr != nil {
		c.JSON(apiErr.Code, apiErr.Message)
		return
	}

	c.JSON(http.StatusOK, fh.dtoToResponse(folderDTO))
}

func (fh *folderHandler) dtoToResponse(folder *dto.FolderDTO) *response.FolderResponse {
	var folders []response.FolderResponse
	if folder.Folders != nil {
		folders = make([]response.FolderResponse, 0)
		for _, v := range folder.Folders {
			folders = append(folders, *fh.dtoToResponse(&v))
		}
	}
	return &response.FolderResponse{
		ID:             folder.ID,
		ParentFolderID: folder.ParentFolderID,
		Name:           folder.Name,
		Path:           folder.Path,
		IsHide:         folder.IsHide,
		Folders:        folders,
		CreatedAt:      folder.CreatedAt,
		UpdatedAt:      folder.UpdatedAt,
	}
}
