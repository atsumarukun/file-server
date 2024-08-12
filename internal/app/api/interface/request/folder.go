package request

type CreateFolderRequest struct {
	ParentFolderID int64  `json:"parent_folder_id"`
	Name           string `json:"name"`
	IsHide         bool   `json:"is_hide"`
}

type UpdateFolderRequest struct {
	Name   string `json:"name"`
	IsHide bool   `json:"is_hide"`
}
