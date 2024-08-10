package model

import (
	"time"

	"gorm.io/gorm"
)

type FolderModel struct {
	ID             int64
	ParentFolderID *int64
	Name           string
	Path           string
	IsHide         bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt
}

func (fm *FolderModel) TableName() string {
	return "folders"
}
