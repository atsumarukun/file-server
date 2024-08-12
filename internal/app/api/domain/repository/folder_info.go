package repository

import (
	"file-server/internal/app/api/domain/entity"

	"gorm.io/gorm"
)

type FolderInfoRepository interface {
	Save(*gorm.DB, *entity.FolderInfo) (*entity.FolderInfo, error)
	FindOneByID(*gorm.DB, int64) (*entity.FolderInfo, error)
	FindOneByPathWithRelationship(*gorm.DB, string) (*entity.FolderInfo, error)
}
