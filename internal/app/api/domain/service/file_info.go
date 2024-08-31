package service

import (
	"file-server/internal/app/api/domain/entity"
	"file-server/internal/app/api/domain/repository"
	"fmt"

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
	if checkFile, err := fs.fileInfoRepository.FindOneByPath(db, path); err != nil {
		return err
	} else if checkFile != nil {
		return fmt.Errorf("%s is already exists", path)
	}
	return nil
}
