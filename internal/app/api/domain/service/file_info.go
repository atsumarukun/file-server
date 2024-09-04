package service

import (
	"errors"
	"file-server/internal/app/api/domain/entity"
	"file-server/internal/app/api/domain/repository"

	"gorm.io/gorm"
)

type FileInfoService interface {
	Exists(*gorm.DB, *entity.FileInfo) error
}

type fileInfoService struct {
	fileInfoRepository repository.FileInfoRepository
}

func NewFileInfoService(fileInfoRepository repository.FileInfoRepository) FileInfoService {
	return &fileInfoService{
		fileInfoRepository: fileInfoRepository,
	}
}

func (fs *fileInfoService) Exists(db *gorm.DB, file *entity.FileInfo) error {
	path := file.GetPath()
	if _, err := fs.fileInfoRepository.FindOneByPath(db, path); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		} else {
			return err
		}
	}
	return nil
}
