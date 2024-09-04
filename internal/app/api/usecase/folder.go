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
	Update(uint64, string, bool) (*dto.FolderDTO, error)
	Remove(uint64) error
	Move(uint64, uint64) (*dto.FolderDTO, error)
	Copy(uint64, uint64) (*dto.FolderDTO, error)
	FindOne(string) (*dto.FolderDTO, error)
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
	var folder *entity.FolderInfo
	if err := fu.db.Transaction(func(tx *gorm.DB) error {
		parentFolder, err := fu.folderInfoRepository.FindOneByID(tx, parentFolderID)
		if err != nil {
			return err
		}

		path := parentFolder.GetPath() + name + "/"

		folderInfo, err := entity.NewFolderInfo(&parentFolderID, name, path, isHide)
		if err != nil {
			return err
		}

		if err := fu.folderInfoService.Exists(tx, folderInfo); err != nil {
			return err
		}

		folder, err = fu.folderInfoRepository.Create(tx, folderInfo)
		if err != nil {
			return err
		}

		folderBody := entity.NewFolderBody(path)

		return fu.folderBodyRepository.Create(folderBody)
	}); err != nil {
		return nil, err
	}

	return fu.entityToDTO(folder), nil
}

func (fu *folderUsecase) Update(id uint64, name string, isHide bool) (*dto.FolderDTO, error) {
	var folder *entity.FolderInfo
	if err := fu.db.Transaction(func(tx *gorm.DB) error {
		folderInfo, err := fu.folderInfoRepository.FindOneByIDWithLower(tx, id)
		if err != nil {
			return err
		}

		oldName := folderInfo.GetName()
		folderInfo.SetIsHide(isHide)

		if name != oldName {
			folderInfo.SetName(name)
			oldPath := folderInfo.GetPath()
			path := oldPath[:strings.LastIndex(oldPath, oldName)] + name + "/"

			if err := folderInfo.Move(oldPath, path); err != nil {
				return err
			}

			if err := fu.folderInfoService.Exists(tx, folderInfo); err != nil {
				return err
			}

			if err := fu.folderBodyRepository.Update(oldPath, path); err != nil {
				return err
			}
		}

		folder, err = fu.folderInfoRepository.Update(tx, folderInfo)
		return err
	}); err != nil {
		return nil, err
	}

	return fu.entityToDTO(folder), nil
}

func (fu *folderUsecase) Remove(id uint64) error {
	if err := fu.db.Transaction(func(tx *gorm.DB) error {
		folderInfo, err := fu.folderInfoRepository.FindOneByIDWithLower(tx, id)
		if err != nil {
			return err
		}

		return fu.folderInfoRepository.Remove(tx, folderInfo)
	}); err != nil {
		return err
	}

	return nil
}

func (fu *folderUsecase) Move(id uint64, parentFolderID uint64) (*dto.FolderDTO, error) {
	var folder *entity.FolderInfo
	if err := fu.db.Transaction(func(tx *gorm.DB) error {
		folderInfo, err := fu.folderInfoRepository.FindOneByIDWithLower(tx, id)
		if err != nil {
			return err
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

		if err := fu.folderInfoService.Exists(tx, folderInfo); err != nil {
			return err
		}

		if err := fu.folderBodyRepository.Update(oldPath, path); err != nil {
			return err
		}

		folder, err = fu.folderInfoRepository.Update(tx, folderInfo)
		return err
	}); err != nil {
		return nil, err
	}

	return fu.entityToDTO(folder), nil
}

func (fu *folderUsecase) Copy(id uint64, parentFolderID uint64) (*dto.FolderDTO, error) {
	var folder *entity.FolderInfo
	if err := fu.db.Transaction(func(tx *gorm.DB) error {
		sourceFolderInfo, err := fu.folderInfoRepository.FindOneByIDWithLower(tx, id)
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

		folder, err = fu.folderInfoRepository.Create(tx, targetFolderInfo)
		return err
	}); err != nil {
		return nil, err
	}

	return fu.entityToDTO(folder), nil
}

func (fu *folderUsecase) FindOne(path string) (*dto.FolderDTO, error) {
	folder, err := fu.folderInfoRepository.FindOneByPathWithChildren(fu.db, path)
	if err != nil {
		return nil, err
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
