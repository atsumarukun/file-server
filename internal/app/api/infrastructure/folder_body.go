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
	return os.MkdirAll(config.STORAGE_PATH+folder.GetPath(), info.Mode())
}
