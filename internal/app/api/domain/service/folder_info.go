package service

import (
	"errors"
	"file-server/internal/app/api/domain/entity"
	"file-server/internal/app/api/domain/repository"

	"gorm.io/gorm"
)

type FolderInfoService interface {
	Exists(*gorm.DB, *entity.FolderInfo) error
}

type folderInfoService struct {
	folderInfoRepository repository.FolderInfoRepository
}

func NewFolderInfoService(folderInfoRepository repository.FolderInfoRepository) FolderInfoService {
	return &folderInfoService{
		folderInfoRepository: folderInfoRepository,
	}
}

func (fs *folderInfoService) Exists(db *gorm.DB, folder *entity.FolderInfo) error {
	path := folder.GetPath()
	if _, err := fs.folderInfoRepository.FindOneByPath(db, path); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		} else {
			return err
		}
	}
	return nil
}
