package request

type CreateFileRequest struct {
	FolderID int64 `form:"folder_id"`
	IsHide   bool  `form:"is_hide"`
}
