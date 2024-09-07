package service

import (
	"errors"
	"file-server/internal/app/api/domain/entity"
	"file-server/internal/app/api/domain/repository"

	"gorm.io/gorm"
)

type FileInfoService interface {
	IsExists(*gorm.DB, *entity.FileInfo) (bool, error)
}

type fileInfoService struct {
	fileInfoRepository repository.FileInfoRepository
}

func NewFileInfoService(fileInfoRepository repository.FileInfoRepository) FileInfoService {
	return &fileInfoService{
		fileInfoRepository: fileInfoRepository,
	}
}

func (fs *fileInfoService) IsExists(db *gorm.DB, file *entity.FileInfo) (bool, error) {
	path := file.GetPath()
	if _, err := fs.fileInfoRepository.FindOneByPath(db, path); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		} else {
			return false, err
		}
	}
	return true, nil
}
