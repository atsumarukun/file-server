package infrastructure

import "file-server/internal/app/api/domain/repository"

type folderBodyInfrastructure struct{}

func NewFolderBodyInfrastructure() repository.FolderBodyRepository {
	return &folderBodyInfrastructure{}
}
