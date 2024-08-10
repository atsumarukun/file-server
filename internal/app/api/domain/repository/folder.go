package repository

import (
	"file-server/internal/app/api/domain/entity"

	"gorm.io/gorm"
)

type FolderRepository interface {
	Create(*gorm.DB, *entity.Folder) (*entity.Folder, error)
	FindOneByID(*gorm.DB, int64) (*entity.Folder, error)
}
