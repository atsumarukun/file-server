package service

import (
	"file-server/internal/app/api/domain/entity"
	"file-server/internal/app/api/domain/repository"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type FolderInfoService interface {
	Exists(*gorm.DB, *entity.FolderInfo) error
	Move(*gorm.DB, *entity.FolderInfo, string) error
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

func (fs *folderInfoService) Move(db *gorm.DB, folder *entity.FolderInfo, path string) error {
	oldPath := folder.GetPath()

	folder.SetPath(path)
	if err := fs.Exists(db, folder); err != nil {
		return err
	}

	lowerFolders, err := fs.folderInfoRepository.FindByIDNotAndPathLike(db, folder.GetID(), oldPath)
	if err != nil {
		return err
	}
	for i := 0; i < len(lowerFolders); i++ {
		lowerFolders[i].SetPath(strings.Replace(lowerFolders[i].GetPath(), oldPath, path, 1))
	}
	_, err = fs.folderInfoRepository.Saves(db, lowerFolders)
	return err
}
