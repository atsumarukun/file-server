package repository

import (
	"file-server/internal/app/api/domain/entity"

	"gorm.io/gorm"
)

type FileInfoRepository interface {
	Save(*gorm.DB, *entity.FileInfo) (*entity.FileInfo, error)
	FindOneByID(*gorm.DB, int64) (*entity.FileInfo, error)
	FindOneByPath(*gorm.DB, string) (*entity.FileInfo, error)
}
