package repository

import (
	"file-server/internal/app/api/domain/entity"

	"gorm.io/gorm"
)

type FileInfoRepository interface {
	Create(*gorm.DB, *entity.FileInfo) (*entity.FileInfo, error)
	Update(*gorm.DB, *entity.FileInfo) (*entity.FileInfo, error)
	Remove(*gorm.DB, *entity.FileInfo) error
	FindOneByID(*gorm.DB, uint64) (*entity.FileInfo, error)
	FindOneByPath(*gorm.DB, string) (*entity.FileInfo, error)
}
