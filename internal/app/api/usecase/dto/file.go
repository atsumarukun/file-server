package dto

import "time"

type FileInfoDTO struct {
	ID        uint64
	FolderID  uint64
	Name      string
	Path      string
	MimeType  string
	IsHide    bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewFileInfoDTO(id uint64, folderID uint64, name string, path string, mimeType string, isHide bool, createdAt time.Time, updatedAt time.Time) *FileInfoDTO {
	return &FileInfoDTO{
		ID:        id,
		FolderID:  folderID,
		Name:      name,
		Path:      path,
		MimeType:  mimeType,
		IsHide:    isHide,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}

type FileBodyDTO struct {
	MimeType string
	Body     []byte
}

func NewFileBodyDTO(mimeType string, body []byte) *FileBodyDTO {
	return &FileBodyDTO{
		MimeType: mimeType,
		Body:     body,
	}
}
