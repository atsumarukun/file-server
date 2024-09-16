package entity

import (
	"fmt"
	"strings"
	"time"
)

type FileName struct {
	Value string
}

func NewFileName(name string) (*FileName, error) {
	var invalidStrings = []string{"\\", "/", ":", "*", "?", "\"", "<", ">", "|"}
	for _, v := range invalidStrings {
		if strings.Contains(name, v) {
			return nil, fmt.Errorf("invalid file name")
		}
	}
	if 128 < len(name) {
		return nil, fmt.Errorf("file name is too long")
	}
	return &FileName{
		Value: name,
	}, nil
}

type FilePath struct {
	Value string
}

func NewFilePath(path string) (*FilePath, error) {
	if path[:1] != "/" {
		return nil, fmt.Errorf("invalid file path")
	}
	if 255 < len(path) {
		return nil, fmt.Errorf("file path is too long")
	}
	return &FilePath{
		Value: path,
	}, nil
}

type MimeType struct {
	Value string
}

func NewMimeType(mimeType string) (*MimeType, error) {
	if !strings.Contains(mimeType, "/") {
		return nil, fmt.Errorf("invalid mime type")
	}
	if 64 < len(mimeType) {
		return nil, fmt.Errorf("mime type is too long")
	}
	return &MimeType{
		Value: mimeType,
	}, nil
}

type FileInfo struct {
	ID        uint64
	FolderID  uint64
	Name      FileName
	Path      FilePath
	MimeType  MimeType
	IsHide    bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewFileInfo(folderID uint64, name string, path string, mimeType string, isHide bool) (*FileInfo, error) {
	fileName, err := NewFileName(name)
	if err != nil {
		return nil, err
	}
	filePath, err := NewFilePath(path)
	if err != nil {
		return nil, err
	}
	fileMimeType, err := NewMimeType(mimeType)
	if err != nil {
		return nil, err
	}

	return &FileInfo{
		FolderID: folderID,
		Name:     *fileName,
		Path:     *filePath,
		MimeType: *fileMimeType,
		IsHide:   isHide,
	}, nil
}

func (f *FileInfo) SetName(name string) error {
	fileName, err := NewFileName(name)
	if err != nil {
		return err
	}
	f.Name = *fileName
	return nil
}

func (f *FileInfo) SetPath(path string) error {
	filePath, err := NewFilePath(path)
	if err != nil {
		return err
	}
	f.Path = *filePath
	return nil
}

func (f *FileInfo) SetMimeType(mimeType string) error {
	fileMimeType, err := NewMimeType(mimeType)
	if err != nil {
		return err
	}
	f.MimeType = *fileMimeType
	return nil
}

func (f *FileInfo) Move(oldPath string, newPath string) error {
	return f.SetPath(strings.Replace(f.Path.Value, oldPath, newPath, 1))
}

func (f *FileInfo) Copy(path string) (*FileInfo, error) {
	return NewFileInfo(0, f.Name.Value, path, f.MimeType.Value, f.IsHide)
}

type FileBody struct {
	Path string
	Body []byte
}

func NewFileBody(path string, body []byte) *FileBody {
	return &FileBody{
		Path: path,
		Body: body,
	}
}

func (f *FileBody) Copy(path string) *FileBody {
	return NewFileBody(path, f.Body)
}
