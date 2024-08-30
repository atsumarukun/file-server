package usecase

import (
	"errors"
	"file-server/internal/app/api/domain/entity"
	"file-server/internal/app/api/domain/repository"
	"file-server/internal/app/api/domain/service"
	"file-server/internal/app/api/usecase/dto"
	apiError "file-server/internal/pkg/errors"
	"fmt"
	"net/http"
	"strings"

	"gorm.io/gorm"
)

type FolderUsecase interface {
	Create(int64, string, bool) (*dto.FolderDTO, *apiError.Error)
	Update(int64, string, bool) (*dto.FolderDTO, *apiError.Error)
	Remove(int64) *apiError.Error
	Move(int64, int64) (*dto.FolderDTO, *apiError.Error)
	FindOne(string) (*dto.FolderDTO, *apiError.Error)
}

type folderUsecase struct {
	db                   *gorm.DB
	folderInfoRepository repository.FolderInfoRepository
	folderBodyRepository repository.FolderBodyRepository
	folderInfoService    service.FolderInfoService
}

func NewFolderUsecase(db *gorm.DB, folderInfoRepository repository.FolderInfoRepository, folderBodyRepository repository.FolderBodyRepository, folderInfoService service.FolderInfoService) FolderUsecase {
	return &folderUsecase{
		db:                   db,
		folderInfoRepository: folderInfoRepository,
		folderBodyRepository: folderBodyRepository,
		folderInfoService:    folderInfoService,
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

		folderInfo, err := entity.NewFolderInfo(&parentFolderID, name, path, isHide)
		if err != nil {
			return err
		}

		if err := fu.folderInfoService.Exists(tx, folderInfo); err != nil {
			return err
		}

		folder, err = fu.folderInfoRepository.Save(tx, folderInfo)
		if err != nil {
			return err
		}

		return fu.folderBodyRepository.Create(path)
	}); err != nil {
		if errors.Is(err, apiError.ErrNotFound) {
			return nil, apiError.NewError(http.StatusNotFound, err.Error())
		} else {
			return nil, apiError.NewError(http.StatusInternalServerError, err.Error())
		}
	}

	return fu.entityToDTO(folder), nil
}

func (fu *folderUsecase) Update(id int64, name string, isHide bool) (*dto.FolderDTO, *apiError.Error) {
	var folder *entity.FolderInfo
	if err := fu.db.Transaction(func(tx *gorm.DB) error {
		folderInfo, err := fu.folderInfoRepository.FindOneByID(tx, id)
		if err != nil {
			return err
		}
		if folderInfo == nil {
			return apiError.ErrNotFound
		}

		oldName := folderInfo.GetName()

		folderInfo.SetName(name)
		folderInfo.SetIsHide(isHide)

		if name != oldName {
			oldPath := folderInfo.GetPath()
			path := oldPath[:strings.LastIndex(oldPath, oldName)] + name + "/"

			if err := fu.folderInfoService.Move(tx, folderInfo, path); err != nil {
				return err
			}

			if err := fu.folderBodyRepository.Update(oldPath, path); err != nil {
				return err
			}
		}

		folder, err = fu.folderInfoRepository.Save(tx, folderInfo)
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

func (fu *folderUsecase) Remove(id int64) *apiError.Error {
	if err := fu.db.Transaction(func(tx *gorm.DB) error {
		folderInfo, err := fu.folderInfoRepository.FindOneByID(tx, id)
		if err != nil {
			return err
		}
		if folderInfo == nil {
			return apiError.ErrNotFound
		}

		lowerFolders, err := fu.folderInfoRepository.FindByIDNotAndPathLike(tx, folderInfo.GetID(), folderInfo.GetPath())
		if err != nil {
			return err
		}

		removeFolders := []entity.FolderInfo{*folderInfo}
		removeFolders = append(removeFolders, lowerFolders...)

		return fu.folderInfoRepository.Removes(tx, removeFolders)
	}); err != nil {
		if errors.Is(err, apiError.ErrNotFound) {
			return apiError.NewError(http.StatusNotFound, err.Error())
		} else {
			return apiError.NewError(http.StatusInternalServerError, err.Error())
		}
	}

	return nil
}

func (fu *folderUsecase) Move(id int64, parentFolderID int64) (*dto.FolderDTO, *apiError.Error) {
	var folder *entity.FolderInfo
	if err := fu.db.Transaction(func(tx *gorm.DB) error {
		folderInfo, err := fu.folderInfoRepository.FindOneByID(tx, id)
		if err != nil {
			return err
		}
		if folderInfo == nil {
			return apiError.ErrNotFound
		}

		parentFolder, err := fu.folderInfoRepository.FindOneByID(tx, parentFolderID)
		if err != nil {
			return err
		}
		if parentFolder == nil {
			return apiError.ErrNotFound
		}

		oldPath := folderInfo.GetPath()

		if strings.Contains(parentFolder.GetPath(), oldPath) {
			return fmt.Errorf("cannot move to lower directory")
		}

		path := parentFolder.GetPath() + folderInfo.GetName() + "/"
		if err := fu.folderInfoService.Move(tx, folderInfo, path); err != nil {
			return err
		}

		if err := fu.folderBodyRepository.Update(oldPath, path); err != nil {
			return err
		}

		folder, err = fu.folderInfoRepository.Save(tx, folderInfo)
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
		folders = make([]dto.FolderDTO, len(folder.GetFolders()))
		for i, v := range folder.GetFolders() {
			folders[i] = *fu.entityToDTO(&v)
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
