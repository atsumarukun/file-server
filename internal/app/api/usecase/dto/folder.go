package dto

import "time"

type FolderDTO struct {
	ID             uint64
	ParentFolderID *uint64
	Name           string
	Path           string
	IsHide         bool
	Folders        []FolderDTO
	Files          []FileDTO
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
