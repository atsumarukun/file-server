package dto

import "time"

type FolderDTO struct {
	ID             int64
	ParentFolderID *int64
	Name           string
	Path           string
	IsHide         bool
	Folders        []FolderDTO
	Files          []FileDTO
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
