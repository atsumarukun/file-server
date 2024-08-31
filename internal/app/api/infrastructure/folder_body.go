package infrastructure

import (
	"file-server/internal/app/api/domain/repository"
	"file-server/internal/pkg/config"
	"os"
)

type folderBodyInfrastructure struct{}

func NewFolderBodyInfrastructure() repository.FolderBodyRepository {
	return &folderBodyInfrastructure{}
}

func (fi *folderBodyInfrastructure) Create(path string) error {
	info, err := os.Lstat(config.STORAGE_PATH)
	if err != nil {
		return err
	}
	return os.MkdirAll(config.STORAGE_PATH+path, info.Mode())
}

func (fi *folderBodyInfrastructure) Update(oldPath string, newPath string) error {
	return os.Rename(config.STORAGE_PATH+oldPath, config.STORAGE_PATH+newPath)
}
