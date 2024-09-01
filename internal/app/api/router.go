package api

import (
	"file-server/internal/app/api/domain/service"
	"file-server/internal/app/api/infrastructure"
	"file-server/internal/app/api/interface/handler"
	"file-server/internal/app/api/usecase"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func route(r *gin.Engine, db *gorm.DB) {
	folderInfoRepository := infrastructure.NewFolderInfoInfrastructure()
	folderBodyRepository := infrastructure.NewFolderBodyInfrastructure()
	folderInfoService := service.NewFolderInfoService(folderInfoRepository, folderBodyRepository)
	folderUsecase := usecase.NewFolderUsecase(db, folderInfoRepository, folderBodyRepository, folderInfoService)
	folderHandler := handler.NewFolderHandler(folderUsecase)

	folders := r.Group("/folders")
	{
		folders.POST("/", folderHandler.Create)
		folders.GET("/*path", folderHandler.FindOne)
		folders.PUT("/:id", folderHandler.Update)
		folders.DELETE("/:id", folderHandler.Remove)
		folders.PUT("/:id/move", folderHandler.Move)
		folders.POST("/:id/copy", folderHandler.Copy)
	}

	fileInfoRepository := infrastructure.NewFileInfoInfrastructure()
	fileBodyRepository := infrastructure.NewFileBodyInfrastructure()
	fileInfoService := service.NewFileInfoService(fileInfoRepository)
	fileUsecase := usecase.NewFileUsecase(db, fileInfoRepository, fileBodyRepository, folderInfoRepository, fileInfoService)
	fileHandler := handler.NewFileHandler(fileUsecase)

	files := r.Group("/files")
	{
		files.POST("/", fileHandler.Create)
		files.PUT("/:id", fileHandler.Update)
		files.DELETE("/:id", fileHandler.Remove)
	}
}
