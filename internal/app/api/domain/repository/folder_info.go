package repository

import (
	"file-server/internal/app/api/domain/entity"

	"gorm.io/gorm"
)

type FolderInfoRepository interface {
	Save(*gorm.DB, *entity.FolderInfo) (*entity.FolderInfo, error)
	Saves(*gorm.DB, []entity.FolderInfo) ([]entity.FolderInfo, error)
	Removes(*gorm.DB, []entity.FolderInfo) error
	FindOneByID(*gorm.DB, uint64) (*entity.FolderInfo, error)
	FindOneByPath(*gorm.DB, string) (*entity.FolderInfo, error)
	FindOneByIDWithRelationship(*gorm.DB, uint64) (*entity.FolderInfo, error)
	FindOneByPathWithRelationship(*gorm.DB, string) (*entity.FolderInfo, error)
	FindByIDNotAndPathLike(*gorm.DB, uint64, string) ([]entity.FolderInfo, error)
}
