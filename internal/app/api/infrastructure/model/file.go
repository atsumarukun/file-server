package model

import (
	"time"

	"gorm.io/gorm"
)

type FileModel struct {
	ID        uint64
	FolderID  uint64
	Name      string
	Path      string
	MimeType  string
	IsHide    bool
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}

func (fm *FileModel) TableName() string {
	return "files"
}