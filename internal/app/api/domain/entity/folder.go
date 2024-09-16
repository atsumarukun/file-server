package entity

import (
	"fmt"
	"strings"
	"time"
)

type FolderName struct {
	Value string
}

func NewFolderName(name string) (*FolderName, error) {
	var invalidStrings = []string{"\\", "/", ":", "*", "?", "\"", "<", ">", "|"}
	for _, v := range invalidStrings {
		if strings.Contains(name, v) {
			return nil, fmt.Errorf("invalid filder name")
		}
	}
	if 128 < len(name) {
		return nil, fmt.Errorf("folder name is too long")
	}
	return &FolderName{
		Value: name,
	}, nil
}

type FolderPath struct {
	Value string
}

func NewFolderPath(path string) (*FolderPath, error) {
	if path[:1] != "/" || path[len(path)-1:] != "/" {
		return nil, fmt.Errorf("invalid folder path")
	}
	if 255 < len(path) {
		return nil, fmt.Errorf("folder path is too long")
	}
	return &FolderPath{
		Value: path,
	}, nil
}

type FolderInfo struct {
	ID             uint64
	ParentFolderID *uint64
	Name           FolderName
	Path           FolderPath
	IsHide         bool
	Folders        []FolderInfo
	Files          []FileInfo
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func NewFolderInfo(parentFolderID *uint64, name string, path string, isHide bool) (*FolderInfo, error) {
	folderName, err := NewFolderName(name)
	if err != nil {
		return nil, err
	}
	folderPath, err := NewFolderPath(path)
	if err != nil {
		return nil, err
	}

	return &FolderInfo{
		ParentFolderID: parentFolderID,
		Name:           *folderName,
		Path:           *folderPath,
		IsHide:         isHide,
	}, nil
}

func (f *FolderInfo) SetName(name string) error {
	folderName, err := NewFolderName(name)
	if err != nil {
		return err
	}
	f.Name = *folderName
	return nil
}

func (f *FolderInfo) SetPath(path string) error {
	folderPath, err := NewFolderPath(path)
	if err != nil {
		return err
	}
	f.Path = *folderPath
	return nil
}

func (f *FolderInfo) IsRoot() bool {
	return f.ParentFolderID == nil
}

func (f *FolderInfo) Move(oldPath string, newPath string) error {
	if 0 < len(f.Folders) {
		for i := 0; i < len(f.Folders); i++ {
			if err := f.Folders[i].Move(oldPath, newPath); err != nil {
				return err
			}
		}
	}
	if 0 < len(f.Files) {
		for i := 0; i < len(f.Files); i++ {
			if err := f.Files[i].Move(oldPath, newPath); err != nil {
				return err
			}
		}
	}
	return f.SetPath(strings.Replace(f.Path.Value, oldPath, newPath, 1))
}

func (f *FolderInfo) Copy(path string) (*FolderInfo, error) {
	folder, err := NewFolderInfo(nil, f.Name.Value, path, f.IsHide)
	if err != nil {
		return nil, err
	}

	if 0 < len(f.Folders) {
		folders := make([]FolderInfo, len(f.Folders))
		for i, v := range f.Folders {
			f, err := v.Copy(path + v.Name.Value + "/")
			if err != nil {
				return nil, err
			}
			folders[i] = *f
		}
		folder.Folders = folders
	}

	if 0 < len(f.Files) {
		files := make([]FileInfo, len(f.Files))
		for i, v := range f.Files {
			f, err := v.Copy(path + v.Name.Value)
			if err != nil {
				return nil, err
			}
			files[i] = *f
		}
		folder.Files = files
	}

	return folder, nil
}

type FolderBody struct {
	Path    string
	Folders []FolderBody
	Files   []FileBody
}

func NewFolderBody(path string) *FolderBody {
	return &FolderBody{
		Path: path,
	}
}

func (f *FolderBody) Copy(path string) *FolderBody {
	folder := NewFolderBody(path)

	folders := make([]FolderBody, len(f.Folders))
	if 0 < len(f.Folders) {
		for i, v := range f.Folders {
			folderPath := v.Path
			folderName := folderPath[strings.LastIndex(folderPath[:len(folderPath)-1], "/")+1:]
			folder := v.Copy(path + folderName + "/")
			folders[i] = *folder
		}
		folder.Folders = folders
	}

	files := make([]FileBody, len(f.Files))
	if 0 < len(f.Files) {
		for i, v := range f.Files {
			filePath := v.Path
			fileName := filePath[strings.LastIndex(filePath, "/")+1:]
			file := v.Copy(path + fileName)
			files[i] = *file
		}
		folder.Files = files
	}

	return folder
}
