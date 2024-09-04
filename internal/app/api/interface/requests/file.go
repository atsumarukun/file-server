package requests

type CreateFileRequest struct {
	FolderID uint64 `form:"folder_id"`
	IsHide   bool   `form:"is_hide"`
}

type UpdateFileRequest struct {
	Name   string `json:"name"`
	IsHide bool   `json:"is_hide"`
}

type MoveFileRequest struct {
	FolderID uint64 `json:"folder_id"`
}

type CopyFileRequest struct {
	FolderID uint64 `json:"folder_id"`
}
