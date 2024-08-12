package repository

import (
	"file-server/internal/app/api/domain/entity"

	"gorm.io/gorm"
)

type FolderInfoRepository interface {
	Save(*gorm.DB, *entity.FolderInfo) (*entity.FolderInfo, error)
	Saves(*gorm.DB, []entity.FolderInfo) ([]entity.FolderInfo, error)
	Remove(*gorm.DB, *entity.FolderInfo) error
	FindOneByID(*gorm.DB, int64) (*entity.FolderInfo, error)
	FindOneByPath(*gorm.DB, string) (*entity.FolderInfo, error)
	FindOneByPathWithRelationship(*gorm.DB, string) (*entity.FolderInfo, error)
	FindByIDNotAndPathLike(*gorm.DB, int64, string) ([]entity.FolderInfo, error)
}
