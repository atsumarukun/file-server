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
	fileInfoRepository := infrastructure.NewFileInfoInfrastructure()
	fileBodyRepository := infrastructure.NewFileBodyInfrastructure()
	credentialRepository := infrastructure.NewCredentialInfrastructure()
	folderInfoService := service.NewFolderInfoService(folderInfoRepository)
	fileInfoService := service.NewFileInfoService(fileInfoRepository)
	folderUsecase := usecase.NewFolderUsecase(db, folderInfoRepository, folderBodyRepository, folderInfoService)
	fileUsecase := usecase.NewFileUsecase(db, fileInfoRepository, fileBodyRepository, folderInfoRepository, fileInfoService)
	authUsecase := usecase.NewAuthUsecase(db, credentialRepository)
	folderHandler := handler.NewFolderHandler(folderUsecase)
	fileHandler := handler.NewFileHandler(fileUsecase)
	authHandler := handler.NewAuthHandler(authUsecase)

	folders := r.Group("/folders")
	{
		folders.POST("/", folderHandler.Create)
		folders.GET("/*path", folderHandler.FindOne)
		folders.PUT("/:id", folderHandler.Update)
		folders.DELETE("/:id", folderHandler.Remove)
		folders.PUT("/:id/move", folderHandler.Move)
		folders.POST("/:id/copy", folderHandler.Copy)
	}

	files := r.Group("/files")
	{
		files.POST("/", fileHandler.Create)

		files.Use(authMiddleware())

		files.PUT("/:id", fileHandler.Update)
		files.DELETE("/:id", fileHandler.Remove)
		files.PUT("/:id/move", fileHandler.Move)
		files.POST("/:id/copy", fileHandler.Copy)
	}

	auth := r.Group("/auth")
	{
		auth.POST("/signin", authHandler.Signin)
	}

	r.Static("/static", "./storage")
}
