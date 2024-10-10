package responses

import "time"

type FileResponse struct {
	ID        uint64    `json:"id"`
	FolderID  uint64    `json:"folder_id"`
	Name      string    `json:"name"`
	Path      string    `json:"path"`
	MimeType  string    `json:"mime_type"`
	IsHide    bool      `json:"is_hide"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewFileResponse(id uint64, folderID uint64, name string, path string, mimeType string, isHide bool, createdAt time.Time, updatedAt time.Time) *FileResponse {
	return &FileResponse{
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
