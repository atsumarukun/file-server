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
	info, err := os.Lstat(config.STORAGE_PATH)
	if err != nil {
		return err
	}
	return os.WriteFile(config.STORAGE_PATH+file.GetPath(), file.GetBody(), info.Mode())
}
