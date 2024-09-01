package response

import "time"

type FolderResponse struct {
	ID             uint64           `json:"id"`
	ParentFolderID *uint64          `json:"parent_folder_id"`
	Name           string           `json:"name"`
	Path           string           `json:"path"`
	IsHide         bool             `json:"is_hide"`
	Folders        []FolderResponse `json:"folders,omitempty"`
	Files          []FileResponse   `json:"files,omitempty"`
	CreatedAt      time.Time        `json:"created_at"`
	UpdatedAt      time.Time        `json:"updated_at"`
}
