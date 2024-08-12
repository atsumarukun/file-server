package service

import (
	"file-server/internal/app/api/domain/repository"
)

type FolderInfoService interface{}

type folderInfoService struct {
	folderInfoRepository repository.FolderInfoRepository
}

func NewFolderInfoService(folderInfoRepository repository.FolderInfoRepository) FolderInfoService {
	return &folderInfoService{
		folderInfoRepository: folderInfoRepository,
	}
}
