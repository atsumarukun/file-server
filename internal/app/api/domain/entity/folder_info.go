package entity

import (
	"fmt"
	"strings"
	"time"
)

type FolderInfo struct {
	id             int64
	parentFolderID *int64
	name           string
	path           string
	isHide         bool
	folders        []FolderInfo
	createdAt      time.Time
	updatedAt      time.Time
}

func NewFolderInfo(parentFolderID *int64, name string, path string, isHide bool) (*FolderInfo, error) {
	folder := &FolderInfo{}

	folder.SetParentFolderID(parentFolderID)
	folder.SetIsHide(isHide)
	if err := folder.SetName(name); err != nil {
		return nil, err
	}
	if err := folder.SetPath(path); err != nil {
		return nil, err
	}

	return folder, nil
}

func (f *FolderInfo) GetID() int64 {
	return f.id
}

func (f *FolderInfo) SetID(id int64) {
	f.id = id
}

func (f *FolderInfo) GetParentFolderID() *int64 {
	return f.parentFolderID
}

func (f *FolderInfo) SetParentFolderID(parentFolderID *int64) {
	f.parentFolderID = parentFolderID
}

func (f *FolderInfo) GetName() string {
	return f.name
}

func (f *FolderInfo) SetName(name string) error {
	var invalidStrings = []string{"\\", "/", ":", "*", "?", "\"", "<", ">", "|"}
	for _, v := range invalidStrings {
		if strings.Contains(name, v) {
			return fmt.Errorf("set folder name: invalid name")
		}
	}
	if 128 < len(name) {
		return fmt.Errorf("set folder name: name is too long")
	}
	f.name = name
	return nil
}

func (f *FolderInfo) GetPath() string {
	return f.path
}

func (f *FolderInfo) SetPath(path string) error {
	if 255 < len(path) {
		return fmt.Errorf("set folder path: path is too long")
	}
	if path[len(path)-1:] != "/" {
		return fmt.Errorf("set folder path: invalid path")
	}
	f.path = path
	return nil
}

func (f *FolderInfo) GetIsHide() bool {
	return f.isHide
}

func (f *FolderInfo) SetIsHide(isHide bool) {
	f.isHide = isHide
}

func (f *FolderInfo) GetFolders() []FolderInfo {
	return f.folders
}

func (f *FolderInfo) SetFolders(folders []FolderInfo) {
	f.folders = folders
}

func (f *FolderInfo) GetCreatedAt() time.Time {
	return f.createdAt
}

func (f *FolderInfo) SetCreatedAt(createdAt time.Time) {
	f.createdAt = createdAt
}

func (f *FolderInfo) GetUpdatedAt() time.Time {
	return f.updatedAt
}

func (f *FolderInfo) SetUpdatedAt(updatedAt time.Time) {
	f.updatedAt = updatedAt
}

func (f *FolderInfo) Move(oldPath string, newPath string) error {
	return f.SetPath(strings.Replace(f.path, oldPath, newPath, 1))
}
