package usecase

import (
	"errors"
	"file-server/internal/app/api/domain/entity"
	"file-server/internal/app/api/domain/repository"
	"file-server/internal/app/api/domain/service"
	"file-server/internal/app/api/usecase/dto"
	apiError "file-server/internal/pkg/errors"
	"net/http"

	"gorm.io/gorm"
)

type FileUsecase interface {
	Create(int64, string, bool, []byte) (*dto.FileDTO, *apiError.Error)
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

func (fu *fileUsecase) Create(folderID int64, name string, isHide bool, body []byte) (*dto.FileDTO, *apiError.Error) {
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
