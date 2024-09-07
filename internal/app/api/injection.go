package api

import (
	"file-server/internal/app/api/domain/repository"
	"file-server/internal/app/api/domain/service"
	"file-server/internal/app/api/infrastructure"
	"file-server/internal/app/api/interface/handler"
	"file-server/internal/app/api/usecase"

	"gorm.io/gorm"
)

var (
	credentialRepository repository.CredentialRepository
	folderInfoRepository repository.FolderInfoRepository
	folderBodyRepository repository.FolderBodyRepository
	fileInfoRepository   repository.FileInfoRepository
	fileBodyRepository   repository.FileBodyRepository

	folderInfoService service.FolderInfoService
	fileInfoService   service.FileInfoService

	authUsecase   usecase.AuthUsecase
	folderUsecase usecase.FolderUsecase
	fileUsecase   usecase.FileUsecase

	authHandler   handler.AuthHandler
	folderHandler handler.FolderHandler
	fileHandler   handler.FileHandler
)

func inject(db *gorm.DB) {
	credentialRepository = infrastructure.NewCredentialInfrastructure()
	folderInfoRepository = infrastructure.NewFolderInfoInfrastructure()
	folderBodyRepository = infrastructure.NewFolderBodyInfrastructure()
	fileInfoRepository = infrastructure.NewFileInfoInfrastructure()
	fileBodyRepository = infrastructure.NewFileBodyInfrastructure()

	folderInfoService = service.NewFolderInfoService(folderInfoRepository)
	fileInfoService = service.NewFileInfoService(fileInfoRepository)

	authUsecase = usecase.NewAuthUsecase(db, credentialRepository)
	folderUsecase = usecase.NewFolderUsecase(db, folderInfoRepository, folderBodyRepository, folderInfoService)
	fileUsecase = usecase.NewFileUsecase(db, fileInfoRepository, fileBodyRepository, folderInfoRepository, fileInfoService)

	authHandler = handler.NewAuthHandler(authUsecase)
	folderHandler = handler.NewFolderHandler(folderUsecase)
	fileHandler = handler.NewFileHandler(fileUsecase)
}
