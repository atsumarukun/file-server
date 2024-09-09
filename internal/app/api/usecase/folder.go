package usecase

import (
	"file-server/internal/app/api/domain/entity"
	"file-server/internal/app/api/domain/repository"
	"file-server/internal/app/api/domain/service"
	"file-server/internal/app/api/usecase/dto"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type FolderUsecase interface {
	Create(uint64, string, bool) (*dto.FolderDTO, error)
	Update(uint64, string, bool, bool) (*dto.FolderDTO, error)
	Remove(uint64, bool) error
	Move(uint64, uint64, bool) (*dto.FolderDTO, error)
	Copy(uint64, uint64, bool) (*dto.FolderDTO, error)
	FindOne(string, bool) (*dto.FolderDTO, error)
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

func (fu *folderUsecase) Create(parentFolderID uint64, name string, isHide bool) (*dto.FolderDTO, error) {
	var folderInfo *entity.FolderInfo
	if err := fu.db.Transaction(func(tx *gorm.DB) error {
		parentFolder, err := fu.folderInfoRepository.FindOneByID(tx, parentFolderID)
		if err != nil {
			return err
		}

		path := parentFolder.GetPath() + name + "/"

		folderInfo, err = entity.NewFolderInfo(&parentFolderID, name, path, isHide)
		if err != nil {
			return err
		}

		if isExists, err := fu.folderInfoService.IsExists(tx, folderInfo); err != nil {
			return err
		} else if isExists {
			return fmt.Errorf("%s is already exists", folderInfo.GetPath())
		}

		folderInfo, err = fu.folderInfoRepository.Create(tx, folderInfo)
		if err != nil {
			return err
		}

		folderBody := entity.NewFolderBody(path)

		return fu.folderBodyRepository.Create(folderBody)
	}); err != nil {
		return nil, err
	}

	return fu.entityToDTO(folderInfo), nil
}

func (fu *folderUsecase) Update(id uint64, name string, isHide bool, isDisplayHiddenObject bool) (*dto.FolderDTO, error) {
	var folderInfo *entity.FolderInfo
	if err := fu.db.Transaction(func(tx *gorm.DB) error {
		var err error
		if isDisplayHiddenObject {
			folderInfo, err = fu.folderInfoRepository.FindOneByIDWithLower(tx, id)
		} else {
			folderInfo, err = fu.folderInfoRepository.FindOneByIDAndIsHideWithLower(tx, id, false)
		}
		if err != nil {
			return err
		}

		if folderInfo.IsRoot() {
			return fmt.Errorf("root directory is not updatable")
		}

		folderInfo.SetIsHide(isHide)

		if name != folderInfo.GetName() {
			oldPath := folderInfo.GetPath()
			path := oldPath[:strings.LastIndex(oldPath, folderInfo.GetName())] + name + "/"

			if err := folderInfo.SetName(name); err != nil {
				return err
			}

			if err := folderInfo.Move(oldPath, path); err != nil {
				return err
			}

			if isExists, err := fu.folderInfoService.IsExists(tx, folderInfo); err != nil {
				return err
			} else if isExists {
				return fmt.Errorf("%s is already exists", folderInfo.GetPath())
			}

			if err := fu.folderBodyRepository.Update(oldPath, path); err != nil {
				return err
			}
		}

		folderInfo, err = fu.folderInfoRepository.Update(tx, folderInfo)
		return err
	}); err != nil {
		return nil, err
	}

	return fu.entityToDTO(folderInfo), nil
}

func (fu *folderUsecase) Remove(id uint64, isDisplayHiddenObject bool) error {
	var folderInfo *entity.FolderInfo
	if err := fu.db.Transaction(func(tx *gorm.DB) error {
		var err error
		if isDisplayHiddenObject {
			folderInfo, err = fu.folderInfoRepository.FindOneByIDWithLower(tx, id)
		} else {
			folderInfo, err = fu.folderInfoRepository.FindOneByIDAndIsHideWithLower(tx, id, false)
		}
		if err != nil {
			return err
		}

		if folderInfo.IsRoot() {
			return fmt.Errorf("root directory is not removable")
		}

		if err := fu.folderBodyRepository.Remove(folderInfo.GetPath()); err != nil {
			return err
		}

		return fu.folderInfoRepository.Remove(tx, folderInfo)
	}); err != nil {
		return err
	}

	return nil
}

func (fu *folderUsecase) Move(id uint64, parentFolderID uint64, isDisplayHiddenObject bool) (*dto.FolderDTO, error) {
	var folderInfo *entity.FolderInfo
	if err := fu.db.Transaction(func(tx *gorm.DB) error {
		var err error
		if isDisplayHiddenObject {
			folderInfo, err = fu.folderInfoRepository.FindOneByIDWithLower(tx, id)
		} else {
			folderInfo, err = fu.folderInfoRepository.FindOneByIDAndIsHideWithLower(tx, id, false)
		}
		if err != nil {
			return err
		}

		if folderInfo.IsRoot() {
			return fmt.Errorf("root directory is not updatable")
		}

		parentFolder, err := fu.folderInfoRepository.FindOneByID(tx, parentFolderID)
		if err != nil {
			return err
		}

		oldPath := folderInfo.GetPath()
		if strings.Contains(parentFolder.GetPath(), oldPath) {
			return fmt.Errorf("cannot move to lower directory")
		}
		path := parentFolder.GetPath() + folderInfo.GetName() + "/"

		if err := folderInfo.Move(oldPath, path); err != nil {
			return err
		}
		folderInfo.SetParentFolderID(&parentFolderID)

		if isExists, err := fu.folderInfoService.IsExists(tx, folderInfo); err != nil {
			return err
		} else if isExists {
			return fmt.Errorf("%s is already exists", folderInfo.GetPath())
		}

		if err := fu.folderBodyRepository.Update(oldPath, path); err != nil {
			return err
		}

		folderInfo, err = fu.folderInfoRepository.Update(tx, folderInfo)
		return err
	}); err != nil {
		return nil, err
	}

	return fu.entityToDTO(folderInfo), nil
}

func (fu *folderUsecase) Copy(id uint64, parentFolderID uint64, isDisplayHiddenObject bool) (*dto.FolderDTO, error) {
	var folderInfo *entity.FolderInfo
	if err := fu.db.Transaction(func(tx *gorm.DB) error {
		var sourceFolderInfo *entity.FolderInfo
		var err error
		if isDisplayHiddenObject {
			sourceFolderInfo, err = fu.folderInfoRepository.FindOneByIDWithLower(tx, id)
		} else {
			sourceFolderInfo, err = fu.folderInfoRepository.FindOneByIDAndIsHideWithLower(tx, id, false)
		}
		if err != nil {
			return err
		}

		sourceFolderBody, err := fu.folderBodyRepository.Read(sourceFolderInfo.GetPath())
		if err != nil {
			return err
		}

		parentFolder, err := fu.folderInfoRepository.FindOneByID(tx, parentFolderID)
		if err != nil {
			return err
		}

		path := parentFolder.GetPath() + sourceFolderInfo.GetName() + "/"
		targetFolderInfo, err := sourceFolderInfo.Copy(path)
		if err != nil {
			return err
		}
		targetFolderInfo.SetParentFolderID(&parentFolderID)

		targetFolderBody := sourceFolderBody.Copy(path)
		if err := fu.folderBodyRepository.Create(targetFolderBody); err != nil {
			return err
		}

		folderInfo, err = fu.folderInfoRepository.Create(tx, targetFolderInfo)
		return err
	}); err != nil {
		return nil, err
	}

	return fu.entityToDTO(folderInfo), nil
}

func (fu *folderUsecase) FindOne(path string, isDisplayHiddenObject bool) (*dto.FolderDTO, error) {
	var folderInfo *entity.FolderInfo
	var err error
	if isDisplayHiddenObject {
		folderInfo, err = fu.folderInfoRepository.FindOneByPathWithChildren(fu.db, path)
	} else {
		folderInfo, err = fu.folderInfoRepository.FindOneByPathAndIsHideWithChildren(fu.db, path, false)
	}
	if err != nil {
		return nil, err
	}

	return fu.entityToDTO(folderInfo), nil
}

func (fu *folderUsecase) entityToDTO(folder *entity.FolderInfo) *dto.FolderDTO {
	var folders []dto.FolderDTO
	if folder.GetFolders() != nil {
		folders = make([]dto.FolderDTO, len(folder.GetFolders()))
		for i, v := range folder.GetFolders() {
			folders[i] = *fu.entityToDTO(&v)
		}
	}
	var files []dto.FileDTO
	if folder.GetFiles() != nil {
		files = make([]dto.FileDTO, len(folder.GetFiles()))
		for i, v := range folder.GetFiles() {
			files[i] = dto.FileDTO{
				ID:        v.GetID(),
				FolderID:  v.GetFolderID(),
				Name:      v.GetName(),
				Path:      v.GetPath(),
				MimeType:  v.GetMimeType(),
				IsHide:    v.GetIsHide(),
				CreatedAt: v.GetCreatedAt(),
				UpdatedAt: v.GetUpdatedAt(),
			}
		}
	}
	return &dto.FolderDTO{
		ID:             folder.GetID(),
		ParentFolderID: folder.GetParentFolderID(),
		Name:           folder.GetName(),
		Path:           folder.GetPath(),
		IsHide:         folder.GetIsHide(),
		Folders:        folders,
		Files:          files,
		CreatedAt:      folder.GetCreatedAt(),
		UpdatedAt:      folder.GetUpdatedAt(),
	}
}
