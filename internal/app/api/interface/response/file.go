package response

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
