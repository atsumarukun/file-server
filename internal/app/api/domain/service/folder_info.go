package service

import (
	"file-server/internal/app/api/domain/entity"
	"file-server/internal/app/api/domain/repository"
	"fmt"

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
	if checkFolder, err := fs.folderInfoRepository.FindOneByPath(db, path); err != nil {
		return err
	} else if checkFolder != nil {
		return fmt.Errorf("%s is already exists", path)
	}
	return nil
}
