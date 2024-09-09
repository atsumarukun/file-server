package model

import "time"

type FileModel struct {
	ID        uint64
	FolderID  uint64
	Name      string
	Path      string
	MimeType  string
	IsHide    bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (fm *FileModel) TableName() string {
	return "files"
}
