package usecase

import (
	"file-server/internal/app/api/domain/entity"
	"file-server/internal/app/api/domain/repository"
	"file-server/internal/app/api/domain/service"
	"file-server/internal/app/api/usecase/dto"
	"net/http"
	"strings"

	"gorm.io/gorm"
)

type FileUsecase interface {
	Create(uint64, string, bool, []byte) (*dto.FileDTO, error)
	Update(uint64, string, bool) (*dto.FileDTO, error)
	Remove(uint64) error
	Move(uint64, uint64) (*dto.FileDTO, error)
	Copy(uint64, uint64) (*dto.FileDTO, error)
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

func (fu *fileUsecase) Create(folderID uint64, name string, isHide bool, body []byte) (*dto.FileDTO, error) {
	var file *entity.FileInfo
	if err := fu.db.Transaction(func(tx *gorm.DB) error {
		parentFolder, err := fu.folderInfoRepository.FindOneByID(tx, folderID)
		if err != nil {
			return err
		}

		path := parentFolder.GetPath() + name
		mimeType := http.DetectContentType(body)

		fileInfo, err := entity.NewFileInfo(folderID, name, path, mimeType, isHide)
		if err != nil {
			return err
		}

		if err := fu.fileInfoService.Exists(tx, fileInfo); err != nil {
			return err
		}

		file, err = fu.fileInfoRepository.Create(tx, fileInfo)
		if err != nil {
			return err
		}

		fileBody := entity.NewFileBody(path, body)

		return fu.fileBodyRepository.Create(fileBody)
	}); err != nil {
		return nil, err
	}

	return fu.entityToDTO(file), nil
}

func (fu *fileUsecase) Update(id uint64, name string, isHide bool) (*dto.FileDTO, error) {
	var file *entity.FileInfo
	if err := fu.db.Transaction(func(tx *gorm.DB) error {
		fileInfo, err := fu.fileInfoRepository.FindOneByID(tx, id)
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

		file, err = fu.fileInfoRepository.Update(tx, fileInfo)
		return err
	}); err != nil {
		return nil, err
	}

	return fu.entityToDTO(file), nil
}

func (fu *fileUsecase) Remove(id uint64) error {
	if err := fu.db.Transaction(func(tx *gorm.DB) error {
		fileInfo, err := fu.fileInfoRepository.FindOneByID(tx, id)
		if err != nil {
			return err
		}

		return fu.fileInfoRepository.Remove(tx, fileInfo)
	}); err != nil {
		return err
	}

	return nil
}

func (fu *fileUsecase) Move(id uint64, folderID uint64) (*dto.FileDTO, error) {
	var file *entity.FileInfo
	if err := fu.db.Transaction(func(tx *gorm.DB) error {
		fileInfo, err := fu.fileInfoRepository.FindOneByID(tx, id)
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

		if err := fu.fileInfoService.Exists(tx, fileInfo); err != nil {
			return err
		}

		if err := fu.fileBodyRepository.Update(oldPath, path); err != nil {
			return err
		}

		file, err = fu.fileInfoRepository.Update(tx, fileInfo)
		return err
	}); err != nil {
		return nil, err
	}

	return fu.entityToDTO(file), nil
}

func (fu *fileUsecase) Copy(id uint64, folderID uint64) (*dto.FileDTO, error) {
	var file *entity.FileInfo
	if err := fu.db.Transaction(func(tx *gorm.DB) error {
		sourceFileInfo, err := fu.fileInfoRepository.FindOneByID(tx, id)
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

		file, err = fu.fileInfoRepository.Create(tx, targetFileInfo)
		return err
	}); err != nil {
		return nil, err
	}

	return fu.entityToDTO(file), nil
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
