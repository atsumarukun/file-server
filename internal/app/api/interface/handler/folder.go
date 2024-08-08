package handler

import "file-server/internal/app/api/usecase"

type FolderHandler interface{}

type folderHandler struct {
	usecase usecase.FolderUsecase
}

func NewFolderHandler(usecase usecase.FolderUsecase) FolderHandler {
	return &folderHandler{
		usecase: usecase,
	}
}
