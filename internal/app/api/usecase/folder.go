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
	FindOne(string) (*dto.FolderDTO, *apiError.Error)
}

type folderUsecase struct {
	db                   *gorm.DB
	folderInfoRepository repository.FolderInfoRepository
}

func NewFolderUsecase(db *gorm.DB, folderInfoRepository repository.FolderInfoRepository) FolderUsecase {
	return &folderUsecase{
		db:                   db,
		folderInfoRepository: folderInfoRepository,
	}
}

func (fu *folderUsecase) Create(parentFolderID int64, name string, isHide bool) (*dto.FolderDTO, *apiError.Error) {
	var folder *entity.FolderInfo
	if err := fu.db.Transaction(func(tx *gorm.DB) error {
		parentFolder, err := fu.folderInfoRepository.FindOneByID(tx, parentFolderID)
		if err != nil {
			return err
		}
		if parentFolder == nil {
			return apiError.ErrNotFound
		}

		path := parentFolder.GetPath() + name + "/"
		folder, err = entity.NewFolderInfo(&parentFolderID, name, path, isHide)
		if err != nil {
			return err
		}

		folder, err = fu.folderInfoRepository.Create(tx, folder)
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

func (fu *folderUsecase) FindOne(path string) (*dto.FolderDTO, *apiError.Error) {
	folder, err := fu.folderInfoRepository.FindOneByPathWithRelationship(fu.db, path)
	if err != nil {
		return nil, apiError.NewError(http.StatusNotFound, err.Error())
	}
	if folder == nil {
		return nil, apiError.NewError(http.StatusNotFound, apiError.ErrNotFound.Error())
	}
	return fu.entityToDTO(folder), nil
}

func (fu *folderUsecase) entityToDTO(folder *entity.FolderInfo) *dto.FolderDTO {
	var folders []dto.FolderDTO
	if folder.GetFolders() != nil {
		folders = make([]dto.FolderDTO, 0)
		for _, v := range folder.GetFolders() {
			folders = append(folders, *fu.entityToDTO(&v))
		}
	}
	return &dto.FolderDTO{
		ID:             folder.GetID(),
		ParentFolderID: folder.GetParentFolderID(),
		Name:           folder.GetName(),
		Path:           folder.GetPath(),
		IsHide:         folder.GetIsHide(),
		Folders:        folders,
		CreatedAt:      folder.GetCreatedAt(),
		UpdatedAt:      folder.GetUpdatedAt(),
	}
}
