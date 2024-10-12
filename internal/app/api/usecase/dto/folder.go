package dto

import "time"

type FolderInfoDTO struct {
	ID             uint64
	ParentFolderID *uint64
	Name           string
	Path           string
	IsHide         bool
	Folders        []FolderInfoDTO
	Files          []FileInfoDTO
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func NewFolderInfoDTO(id uint64, parentFolderID *uint64, name string, path string, isHide bool, folders []FolderInfoDTO, files []FileInfoDTO, createdAt time.Time, updatedAt time.Time) *FolderInfoDTO {
	return &FolderInfoDTO{
		ID:             id,
		ParentFolderID: parentFolderID,
		Name:           name,
		Path:           path,
		IsHide:         isHide,
		Folders:        folders,
		Files:          files,
		CreatedAt:      createdAt,
		UpdatedAt:      updatedAt,
	}
}

type FolderBodyDTO struct {
	MimeType string
	Body     []byte
}

func NewFolderBodyDTO(mimeType string, body []byte) *FolderBodyDTO {
	return &FolderBodyDTO{
		MimeType: mimeType,
		Body:     body,
	}
}
