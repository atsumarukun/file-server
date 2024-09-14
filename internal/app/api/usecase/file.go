package usecase

import (
	"file-server/internal/app/api/domain/entity"
	"file-server/internal/app/api/domain/repository"
	"file-server/internal/app/api/domain/service"
	"file-server/internal/app/api/usecase/dto"
	"file-server/internal/pkg/types"
	"fmt"
	"net/http"
	"strings"

	"gorm.io/gorm"
)

type FileUsecase interface {
	Create(uint64, bool, []types.File) ([]dto.FileDTO, error)
	Update(uint64, string, bool, bool) (*dto.FileDTO, error)
	Remove(uint64, bool) error
	Move(uint64, uint64, bool) (*dto.FileDTO, error)
	Copy(uint64, uint64, bool) (*dto.FileDTO, error)
}

type fileUsecase struct {
	db                   *gorm.DB
	fileInfoRepository   repository.FileInfoRepository
	fileBodyRepository   repository.FileBodyRepository
	folderInfoRepository repository.FolderInfoRepository
	fileInfoService      service.FileInfoService
}

func NewFileUsecase(db *gorm.DB, fileInfoRepository repository.FileInfoRepository, fileBodyRepository repository.FileBodyRepository, folderInfoRepository repository.FolderInfoRepository, fileInfoService service.FileInfoService) FileUsecase {
	return &fileUsecase{
		db:                   db,
		fileInfoRepository:   fileInfoRepository,
		fileBodyRepository:   fileBodyRepository,
		folderInfoRepository: folderInfoRepository,
		fileInfoService:      fileInfoService,
	}
}

func (fu *fileUsecase) Create(folderID uint64, isHide bool, files []types.File) ([]dto.FileDTO, error) {
	fileInfos := make([]entity.FileInfo, len(files))
	if err := fu.db.Transaction(func(tx *gorm.DB) error {
		parentFolder, err := fu.folderInfoRepository.FindOneByID(tx, folderID)
		if err != nil {
			return err
		}

		for i, file := range files {
			path := parentFolder.GetPath() + file.Name
			mimeType := http.DetectContentType(file.Body)

			fileInfo, err := entity.NewFileInfo(folderID, file.Name, path, mimeType, isHide)
			if err != nil {
				return err
			}
			fileInfos[i] = *fileInfo

			if isExists, err := fu.fileInfoService.IsExists(tx, fileInfo); err != nil {
				return err
			} else if isExists {
				return fmt.Errorf("%s is already exists", fileInfo.GetPath())
			}

			fileBody := entity.NewFileBody(path, file.Body)
			if err := fu.fileBodyRepository.Create(fileBody); err != nil {
				return err
			}
		}

		fileInfos, err = fu.fileInfoRepository.Creates(tx, fileInfos)
		return err
	}); err != nil {
		return nil, err
	}

	return fu.entitiesToDTOs(fileInfos), nil
}

func (fu *fileUsecase) Update(id uint64, name string, isHide bool, isDisplayHiddenObject bool) (*dto.FileDTO, error) {
	var fileInfo *entity.FileInfo
	if err := fu.db.Transaction(func(tx *gorm.DB) error {
		var err error
		if isDisplayHiddenObject {
			fileInfo, err = fu.fileInfoRepository.FindOneByID(tx, id)
		} else {
			fileInfo, err = fu.fileInfoRepository.FindOneByIDAndIsHide(tx, id, false)
		}
		if err != nil {
			return err
		}

		fileInfo.SetIsHide(isHide)

		if name != fileInfo.GetName() {
			oldPath := fileInfo.GetPath()
			path := oldPath[:strings.LastIndex(oldPath, fileInfo.GetName())] + name

			if err := fileInfo.SetName(name); err != nil {
				return err
			}

			if err := fileInfo.Move(oldPath, path); err != nil {
				return err
			}

			if err := fu.fileBodyRepository.Update(oldPath, path); err != nil {
				return err
			}
		}

		fileInfo, err = fu.fileInfoRepository.Update(tx, fileInfo)
		return err
	}); err != nil {
		return nil, err
	}

	return fu.entityToDTO(fileInfo), nil
}

func (fu *fileUsecase) Remove(id uint64, isDisplayHiddenObject bool) error {
	var fileInfo *entity.FileInfo
	if err := fu.db.Transaction(func(tx *gorm.DB) error {
		var err error
		if isDisplayHiddenObject {
			fileInfo, err = fu.fileInfoRepository.FindOneByID(tx, id)
		} else {
			fileInfo, err = fu.fileInfoRepository.FindOneByIDAndIsHide(tx, id, false)
		}
		if err != nil {
			return err
		}

		if err := fu.fileBodyRepository.Remove(fileInfo.GetPath()); err != nil {
			return err
		}

		return fu.fileInfoRepository.Remove(tx, fileInfo)
	}); err != nil {
		return err
	}

	return nil
}

func (fu *fileUsecase) Move(id uint64, folderID uint64, isDisplayHiddenObject bool) (*dto.FileDTO, error) {
	var fileInfo *entity.FileInfo
	if err := fu.db.Transaction(func(tx *gorm.DB) error {
		var err error
		if isDisplayHiddenObject {
			fileInfo, err = fu.fileInfoRepository.FindOneByID(tx, id)
		} else {
			fileInfo, err = fu.fileInfoRepository.FindOneByIDAndIsHide(tx, id, false)
		}
		if err != nil {
			return err
		}

		parentFolder, err := fu.folderInfoRepository.FindOneByID(tx, folderID)
		if err != nil {
			return err
		}

		oldPath := fileInfo.GetPath()
		path := parentFolder.GetPath() + fileInfo.GetName()

		if err := fileInfo.Move(oldPath, path); err != nil {
			return err
		}
		fileInfo.SetFolderID(folderID)

		if isExists, err := fu.fileInfoService.IsExists(tx, fileInfo); err != nil {
			return err
		} else if isExists {
			return fmt.Errorf("%s is already exists", fileInfo.GetPath())
		}

		if err := fu.fileBodyRepository.Update(oldPath, path); err != nil {
			return err
		}

		fileInfo, err = fu.fileInfoRepository.Update(tx, fileInfo)
		return err
	}); err != nil {
		return nil, err
	}

	return fu.entityToDTO(fileInfo), nil
}

func (fu *fileUsecase) Copy(id uint64, folderID uint64, isDisplayHiddenObject bool) (*dto.FileDTO, error) {
	var fileInfo *entity.FileInfo
	if err := fu.db.Transaction(func(tx *gorm.DB) error {
		var sourceFileInfo *entity.FileInfo
		var err error
		if isDisplayHiddenObject {
			sourceFileInfo, err = fu.fileInfoRepository.FindOneByID(tx, id)
		} else {
			sourceFileInfo, err = fu.fileInfoRepository.FindOneByIDAndIsHide(tx, id, false)
		}
		if err != nil {
			return err
		}

		sourceFileBody, err := fu.fileBodyRepository.Read(sourceFileInfo.GetPath())
		if err != nil {
			return err
		}

		parentFolder, err := fu.folderInfoRepository.FindOneByID(tx, folderID)
		if err != nil {
			return err
		}

		path := parentFolder.GetPath() + sourceFileInfo.GetName()
		targetFileInfo, err := sourceFileInfo.Copy(path)
		if err != nil {
			return err
		}
		targetFileInfo.SetFolderID(folderID)

		targetFileBody := sourceFileBody.Copy(path)
		if err := fu.fileBodyRepository.Create(targetFileBody); err != nil {
			return err
		}

		fileInfo, err = fu.fileInfoRepository.Create(tx, targetFileInfo)
		return err
	}); err != nil {
		return nil, err
	}

	return fu.entityToDTO(fileInfo), nil
}

func (fu *fileUsecase) entityToDTO(file *entity.FileInfo) *dto.FileDTO {
	return &dto.FileDTO{
		ID:        file.GetID(),
		FolderID:  file.GetFolderID(),
		Name:      file.GetName(),
		Path:      file.GetPath(),
		MimeType:  file.GetMimeType(),
		IsHide:    file.GetIsHide(),
		CreatedAt: file.GetCreatedAt(),
		UpdatedAt: file.GetUpdatedAt(),
	}
}

func (fu *fileUsecase) entitiesToDTOs(files []entity.FileInfo) []dto.FileDTO {
	dtos := make([]dto.FileDTO, len(files))
	for i, file := range files {
		f := fu.entityToDTO(&file)
		dtos[i] = *f
	}
	return dtos
}
