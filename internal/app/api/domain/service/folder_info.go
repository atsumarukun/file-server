package service

import (
	"errors"
	"file-server/internal/app/api/domain/entity"
	"file-server/internal/app/api/domain/repository"

	"gorm.io/gorm"
)

type FolderInfoService interface {
	IsExists(*gorm.DB, *entity.FolderInfo) (bool, error)
}

type folderInfoService struct {
	folderInfoRepository repository.FolderInfoRepository
}

func NewFolderInfoService(folderInfoRepository repository.FolderInfoRepository) FolderInfoService {
	return &folderInfoService{
		folderInfoRepository: folderInfoRepository,
	}
}

func (fs *folderInfoService) IsExists(db *gorm.DB, folder *entity.FolderInfo) (bool, error) {
	path := folder.GetPath()
	if _, err := fs.folderInfoRepository.FindOneByPath(db, path); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		} else {
			return false, err
		}
	}
	return true, nil
}
