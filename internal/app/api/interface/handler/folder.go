package handler

import (
	"file-server/internal/app/api/interface/request"
	"file-server/internal/app/api/interface/response"
	"file-server/internal/app/api/usecase"
	"file-server/internal/app/api/usecase/dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

type FolderHandler interface {
	Create(*gin.Context)
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
		c.JSON(err.Code, err.Message)
		return
	}

	c.JSON(http.StatusOK, fh.dtoToResponse(folderDTO))
}

func (fh *folderHandler) dtoToResponse(folder *dto.FolderDTO) *response.FolderResponse {
	return &response.FolderResponse{
		ID:             folder.ID,
		ParentFolderID: folder.ParentFolderID,
		Name:           folder.Name,
		Path:           folder.Path,
		IsHide:         folder.IsHide,
		CreatedAt:      folder.CreatedAt,
		UpdatedAt:      folder.UpdatedAt,
	}
}
