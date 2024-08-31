package dto

import "time"

type FileDTO struct {
	ID        int64
	FolderID  int64
	Name      string
	Path      string
	MimeType  string
	IsHide    bool
	CreatedAt time.Time
	UpdatedAt time.Time
}
