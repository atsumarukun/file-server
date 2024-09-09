package model

import "time"

type FolderModel struct {
	ID             uint64
	ParentFolderID *uint64
	Name           string
	Path           string
	IsHide         bool
	Folders        []FolderModel `gorm:"foreignkey:ParentFolderID"`
	Files          []FileModel   `gorm:"foreignkey:FolderID"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (fm *FolderModel) TableName() string {
	return "folders"
}
