package usecase

import (
	"file-server/internal/app/api/domain/repository"

	"gorm.io/gorm"
)

type FolderUsecase interface{}

type folderUsecase struct {
	db               *gorm.DB
	folderRepository repository.FolderRepository
}

func NewFolderUsecase(db *gorm.DB, folderRepository repository.FolderRepository) FolderUsecase {
	return &folderUsecase{
		db:               db,
		folderRepository: folderRepository,
	}
}
