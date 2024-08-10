package usecase

import (
	"errors"
	"file-server/internal/app/api/domain/entity"
	"file-server/internal/app/api/domain/repository"
	"file-server/internal/app/api/usecase/dto"
	apiError "file-server/internal/pkg/errors"
	"net/http"

	"gorm.io/gorm"
)

type FolderUsecase interface {
	Create(int64, string, bool) (*dto.FolderDTO, *apiError.Error)
}

type folderUsecase struct {
	db               *gorm.DB
	folderRepository repository.FolderRepository
}

func NewFolderUsecase(db *gorm.DB, folderRepository repository.FolderRepository) FolderUsecase {
	return &folderUsecase{
		db:               db,
		folderRepository: folderRepository,
	}
}

func (fu *folderUsecase) Create(parentFolderID int64, name string, isHide bool) (*dto.FolderDTO, *apiError.Error) {
	var folder *entity.Folder
	if err := fu.db.Transaction(func(tx *gorm.DB) error {
		parentFolder, err := fu.folderRepository.FindOneByID(tx, parentFolderID)
		if err != nil {
			return err
		}
		if parentFolder == nil {
			return apiError.ErrNotFound
		}

		path := parentFolder.GetPath() + name + "/"
		folder, err = entity.NewFolder(&parentFolderID, name, path, isHide)
		if err != nil {
			return err
		}

		folder, err = fu.folderRepository.Create(tx, folder)
		return err
	}); err != nil {
		if errors.Is(err, apiError.ErrNotFound) {
			return nil, apiError.NewError(http.StatusNotFound, err.Error())
		} else {
			return nil, apiError.NewError(http.StatusInternalServerError, err.Error())
		}
	}

	return fu.entityToDTO(folder), nil
}

func (fu *folderUsecase) entityToDTO(folder *entity.Folder) *dto.FolderDTO {
	return &dto.FolderDTO{
		ID:             folder.GetID(),
		ParentFolderID: folder.GetParentFolderID(),
		Name:           folder.GetName(),
		Path:           folder.GetPath(),
		IsHide:         folder.GetIsHide(),
		CreatedAt:      folder.GetCreatedAt(),
		UpdatedAt:      folder.GetUpdatedAt(),
	}
}
