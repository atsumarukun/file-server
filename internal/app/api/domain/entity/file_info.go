package entity

import (
	"fmt"
	"strings"
	"time"
)

type FileInfo struct {
	id        uint64
	folderID  uint64
	name      string
	path      string
	mimeType  string
	isHide    bool
	createdAt time.Time
	updatedAt time.Time
}

func NewFileInfo(folderID uint64, name string, path string, mimeType string, isHide bool) (*FileInfo, error) {
	file := &FileInfo{}

	file.SetFolderID(folderID)
	file.SetIsHide(isHide)
	if err := file.SetName(name); err != nil {
		return nil, err
	}
	if err := file.SetPath(path); err != nil {
		return nil, err
	}
	if err := file.SetMimeType(mimeType); err != nil {
		return nil, err
	}

	return file, nil
}

func (f *FileInfo) GetID() uint64 {
	return f.id
}

func (f *FileInfo) SetID(id uint64) {
	f.id = id
}

func (f *FileInfo) GetFolderID() uint64 {
	return f.folderID
}

func (f *FileInfo) SetFolderID(folderID uint64) {
	f.folderID = folderID
}

func (f *FileInfo) GetName() string {
	return f.name
}

func (f *FileInfo) SetName(name string) error {
	var invalidStrings = []string{"\\", "/", ":", "*", "?", "\"", "<", ">", "|"}
	for _, v := range invalidStrings {
		if strings.Contains(name, v) {
			return fmt.Errorf("set file name: invalid name")
		}
	}
	if 128 < len(name) {
		return fmt.Errorf("set file name: name is too long")
	}
	f.name = name
	return nil
}

func (f *FileInfo) GetPath() string {
	return f.path
}

func (f *FileInfo) SetPath(path string) error {
	if 255 < len(path) {
		return fmt.Errorf("set file path: path is too long")
	}
	f.path = path
	return nil
}

func (f *FileInfo) GetMimeType() string {
	return f.mimeType
}

func (f *FileInfo) SetMimeType(mimeType string) error {
	if !strings.Contains(mimeType, "/") {
		return fmt.Errorf("set file name: invalid name")
	}
	if 64 < len(mimeType) {
		return fmt.Errorf("set file MIME type: MIME type is too long")
	}
	f.mimeType = mimeType
	return nil
}

func (f *FileInfo) GetIsHide() bool {
	return f.isHide
}

func (f *FileInfo) SetIsHide(isHide bool) {
	f.isHide = isHide
}

func (f *FileInfo) GetCreatedAt() time.Time {
	return f.createdAt
}

func (f *FileInfo) SetCreatedAt(createdAt time.Time) {
	f.createdAt = createdAt
}

func (f *FileInfo) GetUpdatedAt() time.Time {
	return f.updatedAt
}

func (f *FileInfo) SetUpdatedAt(updatedAt time.Time) {
	f.updatedAt = updatedAt
}

func (f *FileInfo) Move(oldPath string, newPath string) error {
	return f.SetPath(strings.Replace(f.path, oldPath, newPath, 1))
}

func (f *FileInfo) Copy(path string) (*FileInfo, error) {
	return NewFileInfo(0, f.name, path, f.mimeType, f.isHide)
}
