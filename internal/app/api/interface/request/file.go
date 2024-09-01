package request

type CreateFileRequest struct {
	FolderID int64 `form:"folder_id"`
	IsHide   bool  `form:"is_hide"`
}

type UpdateFileRequest struct {
	Name   string `json:"name"`
	IsHide bool   `json:"is_hide"`
}

type MoveFileRequest struct {
	FolderID int64 `json:"folder_id"`
}
