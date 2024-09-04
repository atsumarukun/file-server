package entity

import "strings"

type FolderBody struct {
	path    string
	folders []FolderBody
	files   []FileBody
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

func (f *FolderBody) GetFolders() []FolderBody {
	return f.folders
}

func (f *FolderBody) SetFolders(folders []FolderBody) {
	f.folders = folders
}

func (f *FolderBody) GetFiles() []FileBody {
	return f.files
}

func (f *FolderBody) SetFiles(files []FileBody) {
	f.files = files
}

func (f *FolderBody) Copy(path string) *FolderBody {
	folder := NewFolderBody(path)

	folders := make([]FolderBody, len(f.folders))
	if 0 < len(f.folders) {
		for i, v := range f.folders {
			folderPath := v.GetPath()
			folderName := folderPath[strings.LastIndex(folderPath[:len(folderPath)-1], "/")+1:]
			folder := v.Copy(path + folderName + "/")
			folders[i] = *folder
		}
		folder.folders = folders
	}

	files := make([]FileBody, len(f.files))
	if 0 < len(f.files) {
		for i, v := range f.files {
			filePath := v.GetPath()
			fileName := filePath[strings.LastIndex(filePath, "/")+1:]
			file := v.Copy(path + fileName)
			files[i] = *file
		}
		folder.files = files
	}

	return folder
}
