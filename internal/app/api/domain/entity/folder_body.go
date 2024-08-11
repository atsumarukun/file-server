package entity

type FolderBody struct {
	path string
}

func NewFolderBody(path string) *FolderBody {
	folder := &FolderBody{}

	folder.SetPath(path)

	return folder
}

func (f *FolderBody) GetPath() string {
	return f.path
}

func (f *FolderBody) SetPath(path string) {
	f.path = path
}
