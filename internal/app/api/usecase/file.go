package usecase

import (
	"errors"
	"file-server/internal/app/api/domain/entity"
	"file-server/internal/app/api/domain/repository"
	"file-server/internal/app/api/domain/service"
	"file-server/internal/app/api/usecase/dto"
	apiError "file-server/internal/pkg/errors"
	"net/http"
	"strings"

	"gorm.io/gorm"
)

type FileUsecase interface {
	Create(uint64, string, bool, []byte) (*dto.FileDTO, *apiError.Error)
	Update(uint64, string, bool) (*dto.FileDTO, *apiError.Error)
	Remove(uint64) *apiError.Error
	Move(uint64, uint64) (*dto.FileDTO, *apiError.Error)
	Copy(uint64, uint64) (*dto.FileDTO, *apiError.Error)
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

func (fu *fileUsecase) Create(folderID uint64, name string, isHide bool, body []byte) (*dto.FileDTO, *apiError.Error) {
	var file *entity.FileInfo
	if err := fu.db.Transaction(func(tx *gorm.DB) error {
		parentFolder, err := fu.folderInfoRepository.FindOneByID(tx, folderID)
		if err != nil {
			return err
		}
		if parentFolder == nil {
			return apiError.ErrNotFound
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

		file, err = fu.fileInfoRepository.Save(tx, fileInfo)
		if err != nil {
			return err
		}

		fileBody := entity.NewFileBody(path, body)

		return fu.fileBodyRepository.Create(fileBody)
	}); err != nil {
		if errors.Is(err, apiError.ErrNotFound) {
			return nil, apiError.NewError(http.StatusNotFound, err.Error())
		} else {
			return nil, apiError.NewError(http.StatusInternalServerError, err.Error())
		}
	}

	return fu.entityToDTO(file), nil
}

func (fu *fileUsecase) Update(id uint64, name string, isHide bool) (*dto.FileDTO, *apiError.Error) {
	var file *entity.FileInfo
	if err := fu.db.Transaction(func(tx *gorm.DB) error {
		fileInfo, err := fu.fileInfoRepository.FindOneByID(tx, id)
		if err != nil {
			return err
		}
		if fileInfo == nil {
			return apiError.ErrNotFound
		}

		fileInfo.SetIsHide(isHide)

		if name != fileInfo.GetName() {
			oldPath := fileInfo.GetPath()
			path := oldPath[:strings.LastIndex(oldPath, fileInfo.GetName())] + name

			if err := fileInfo.SetPath(path); err != nil {
				return err
			}
			if err := fileInfo.SetName(name); err != nil {
				return err
			}

			if err := fu.fileBodyRepository.Update(oldPath, path); err != nil {
				return err
			}
		}

		file, err = fu.fileInfoRepository.Save(tx, fileInfo)
		return err
	}); err != nil {
		if errors.Is(err, apiError.ErrNotFound) {
			return nil, apiError.NewError(http.StatusNotFound, err.Error())
		} else {
			return nil, apiError.NewError(http.StatusInternalServerError, err.Error())
		}
	}

	return fu.entityToDTO(file), nil
}

func (fu *fileUsecase) Remove(id uint64) *apiError.Error {
	if err := fu.db.Transaction(func(tx *gorm.DB) error {
		fileInfo, err := fu.fileInfoRepository.FindOneByID(tx, id)
		if err != nil {
			return err
		}
		if fileInfo == nil {
			return apiError.ErrNotFound
		}

		return fu.fileInfoRepository.Remove(tx, fileInfo)
	}); err != nil {
		if errors.Is(err, apiError.ErrNotFound) {
			return apiError.NewError(http.StatusNotFound, err.Error())
		} else {
			return apiError.NewError(http.StatusInternalServerError, err.Error())
		}
	}

	return nil
}

func (fu *fileUsecase) Move(id uint64, folderID uint64) (*dto.FileDTO, *apiError.Error) {
	var file *entity.FileInfo
	if err := fu.db.Transaction(func(tx *gorm.DB) error {
		fileInfo, err := fu.fileInfoRepository.FindOneByID(tx, id)
		if err != nil {
			return err
		}
		if fileInfo == nil {
			return apiError.ErrNotFound
		}

		parentFolder, err := fu.folderInfoRepository.FindOneByID(tx, folderID)
		if err != nil {
			return err
		}
		if parentFolder == nil {
			return apiError.ErrNotFound
		}

		oldPath := fileInfo.GetPath()
		path := parentFolder.GetPath() + fileInfo.GetName()
		if err := fileInfo.SetPath(path); err != nil {
			return err
		}
		fileInfo.SetFolderID(folderID)

		if err := fu.fileBodyRepository.Update(oldPath, path); err != nil {
			return err
		}

		file, err = fu.fileInfoRepository.Save(tx, fileInfo)
		return err
	}); err != nil {
		if errors.Is(err, apiError.ErrNotFound) {
			return nil, apiError.NewError(http.StatusNotFound, err.Error())
		} else {
			return nil, apiError.NewError(http.StatusInternalServerError, err.Error())
		}
	}

	return fu.entityToDTO(file), nil
}

func (fu *fileUsecase) Copy(id uint64, folderID uint64) (*dto.FileDTO, *apiError.Error) {
	var file *entity.FileInfo
	if err := fu.db.Transaction(func(tx *gorm.DB) error {
		sourceFileInfo, err := fu.fileInfoRepository.FindOneByID(tx, id)
		if err != nil {
			return err
		}
		if sourceFileInfo == nil {
			return apiError.ErrNotFound
		}

		sourceFileBody, err := fu.fileBodyRepository.Read(sourceFileInfo.GetPath())
		if err != nil {
			return err
		}

		parentFolder, err := fu.folderInfoRepository.FindOneByID(tx, folderID)
		if err != nil {
			return err
		}
		if parentFolder == nil {
			return apiError.ErrNotFound
		}

		path := parentFolder.GetPath() + sourceFileInfo.GetName()
		targetFileInfo, err := sourceFileInfo.Copy(path)
		if err != nil {
			return err
		}
		targetFileInfo.SetFolderID(folderID)

		targetFileBody := entity.NewFileBody(path, sourceFileBody.GetBody())
		if err := fu.fileBodyRepository.Create(targetFileBody); err != nil {
			return err
		}

		file, err = fu.fileInfoRepository.Save(tx, targetFileInfo)
		return err
	}); err != nil {
		if errors.Is(err, apiError.ErrNotFound) {
			return nil, apiError.NewError(http.StatusNotFound, err.Error())
		} else {
			return nil, apiError.NewError(http.StatusInternalServerError, err.Error())
		}
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
