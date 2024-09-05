package infrastructure

import (
	"file-server/internal/app/api/domain/entity"
	"file-server/internal/app/api/domain/repository"
	"file-server/internal/pkg/config"
	"os"
)

type folderBodyInfrastructure struct{}

func NewFolderBodyInfrastructure() repository.FolderBodyRepository {
	return &folderBodyInfrastructure{}
}

func (fi *folderBodyInfrastructure) Create(folder *entity.FolderBody) error {
	info, err := os.Lstat("./")
	if err != nil {
		return err
	}

	if err := os.MkdirAll(config.STORAGE_PATH+folder.GetPath(), info.Mode()); err != nil {
		return err
	}

	folders := folder.GetFolders()
	if 0 < len(folders) {
		for _, v := range folders {
			if err := fi.Create(&v); err != nil {
				return err
			}
		}
	}

	files := folder.GetFiles()
	if 0 < len(files) {
		for _, v := range files {
			if err := os.WriteFile(config.STORAGE_PATH+v.GetPath(), v.GetBody(), info.Mode()); err != nil {
				return err
			}
		}
	}

	return nil
}

func (fi *folderBodyInfrastructure) Update(oldPath string, newPath string) error {
	return os.Rename(config.STORAGE_PATH+oldPath, config.STORAGE_PATH+newPath)
}

func (fi *folderBodyInfrastructure) Read(path string) (*entity.FolderBody, error) {
	entry, err := os.ReadDir(config.STORAGE_PATH + path)
	if err != nil {
		return nil, err
	}

	folder := entity.NewFolderBody(path)

	var folders []entity.FolderBody
	var files []entity.FileBody

	for _, v := range entry {
		if v.IsDir() {
			f, err := fi.Read(path + v.Name() + "/")
			if err != nil {
				return nil, err
			}
			folders = append(folders, *f)
		} else {
			body, err := os.ReadFile(config.STORAGE_PATH + path + v.Name())
			if err != nil {
				return nil, err
			}
			file := entity.NewFileBody(config.STORAGE_PATH+path+v.Name(), body)
			files = append(files, *file)
		}
	}

	folder.SetFolders(folders)
	folder.SetFiles(files)

	return folder, nil
}
