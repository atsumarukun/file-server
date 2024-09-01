package dto

import "time"

type FileDTO struct {
	ID        uint64
	FolderID  uint64
	Name      string
	Path      string
	MimeType  string
	IsHide    bool
	CreatedAt time.Time
	UpdatedAt time.Time
}
