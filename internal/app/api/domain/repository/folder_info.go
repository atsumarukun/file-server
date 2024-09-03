package repository

import (
	"file-server/internal/app/api/domain/entity"

	"gorm.io/gorm"
)

type FolderInfoRepository interface {
	Create(*gorm.DB, *entity.FolderInfo) (*entity.FolderInfo, error)
	Update(*gorm.DB, *entity.FolderInfo) (*entity.FolderInfo, error)
	Remove(*gorm.DB, *entity.FolderInfo) error
	FindOneByID(*gorm.DB, uint64) (*entity.FolderInfo, error)
	FindOneByPath(*gorm.DB, string) (*entity.FolderInfo, error)
	FindOneByPathWithChildren(*gorm.DB, string) (*entity.FolderInfo, error)
	FindOneByIDWithLower(*gorm.DB, uint64) (*entity.FolderInfo, error)
}
