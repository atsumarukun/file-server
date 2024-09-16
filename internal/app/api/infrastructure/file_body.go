package infrastructure

import (
	"file-server/internal/app/api/domain/entity"
	"file-server/internal/app/api/domain/repository"
	"file-server/internal/pkg/config"
	"os"
)

type fileBodyInfrastructure struct{}

func NewFileBodyInfrastructure() repository.FileBodyRepository {
	return &fileBodyInfrastructure{}
}

func (fi *fileBodyInfrastructure) Create(file *entity.FileBody) error {
	info, err := os.Lstat("./")
	if err != nil {
		return err
	}
	return os.WriteFile(config.STORAGE_PATH+file.Path, file.Body, info.Mode())
}

func (fi *fileBodyInfrastructure) Update(oldPath string, newPath string) error {
	return os.Rename(config.STORAGE_PATH+oldPath, config.STORAGE_PATH+newPath)
}

func (fi *fileBodyInfrastructure) Remove(path string) error {
	return os.Remove(config.STORAGE_PATH + path)
}

func (fi *fileBodyInfrastructure) Read(path string) (*entity.FileBody, error) {
	body, err := os.ReadFile(config.STORAGE_PATH + path)
	if err != nil {
		return nil, err
	}
	return entity.NewFileBody(path, body), nil
}
