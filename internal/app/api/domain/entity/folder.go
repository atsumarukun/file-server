package entity

import (
	"fmt"
	"strings"
	"time"
)

type Folder struct {
	id             int64
	parentFolderID *int64
	name           string
	path           string
	isHide         bool
	createdAt      time.Time
	updatedAt      time.Time
}

func NewFolder(parentFolderID *int64, name string, path string, isHide bool) (*Folder, error) {
	folder := &Folder{}

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

func (f *Folder) GetID() int64 {
	return f.id
}

func (f *Folder) SetID(id int64) {
	f.id = id
}

func (f *Folder) GetParentFolderID() *int64 {
	return f.parentFolderID
}

func (f *Folder) SetParentFolderID(parentFolderID *int64) {
	f.parentFolderID = parentFolderID
}

func (f *Folder) GetName() string {
	return f.name
}

func (f *Folder) SetName(name string) error {
	if 128 < len(name) {
		return fmt.Errorf("set folder name: name is too long")
	}
	if strings.Contains(name, "/") {
		return fmt.Errorf("set folder name: invalid name")
	}
	f.name = name
	return nil
}

func (f *Folder) GetPath() string {
	return f.path
}

func (f *Folder) SetPath(path string) error {
	if 255 < len(path) {
		return fmt.Errorf("set folder path: path is too long")
	}
	f.path = path
	return nil
}

func (f *Folder) GetIsHide() bool {
	return f.isHide
}

func (f *Folder) SetIsHide(isHide bool) {
	f.isHide = isHide
}

func (f *Folder) GetCreatedAt() time.Time {
	return f.createdAt
}

func (f *Folder) SetCreatedAt(createdAt time.Time) {
	f.createdAt = createdAt
}

func (f *Folder) GetUpdatedAt() time.Time {
	return f.updatedAt
}

func (f *Folder) SetUpdatedAt(updatedAt time.Time) {
	f.updatedAt = updatedAt
}
