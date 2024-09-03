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
	info, err := os.Lstat("./")
	if err != nil {
		return err
	}
	return os.MkdirAll(config.STORAGE_PATH+path, info.Mode())
}

func (fi *folderBodyInfrastructure) Update(oldPath string, newPath string) error {
	return os.Rename(config.STORAGE_PATH+oldPath, config.STORAGE_PATH+newPath)
}

func (fi *folderBodyInfrastructure) Copy(oldPath string, newPath string) error {
	info, err := os.Lstat("./")
	if err != nil {
		return err
	}

	if err := os.Mkdir(config.STORAGE_PATH+newPath, info.Mode()); err != nil {
		return err
	}

	entry, err := os.ReadDir(config.STORAGE_PATH + oldPath)
	if err != nil {
		return err
	}

	for _, v := range entry {
		if v.IsDir() {
			if err := fi.Copy(oldPath+v.Name()+"/", newPath+v.Name()+"/"); err != nil {
				return err
			}
		} else {
			body, err := os.ReadFile(config.STORAGE_PATH + oldPath + v.Name())
			if err != nil {
				return err
			}
			if err := os.WriteFile(config.STORAGE_PATH+newPath+v.Name(), body, info.Mode()); err != nil {
				return err
			}
		}
	}

	return nil
}
