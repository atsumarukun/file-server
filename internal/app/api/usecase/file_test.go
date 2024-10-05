package usecase

import (
	"file-server/internal/app/api/domain/entity"
	"file-server/internal/pkg/types"
	"file-server/test/database"
	mock_repository "file-server/test/mock/repository"
	mock_service "file-server/test/mock/service"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestCreateFile(t *testing.T) {
	db, mock, err := database.Open()
	if err != nil {
		t.Error(err.Error())
	}
	mock.ExpectBegin()
	mock.ExpectCommit()

	fileInfo, err := entity.NewFileInfo(1, "name", "/path/", "mime/type", false)
	if err != nil {
		t.Error(err.Error())
	}

	fileBody := entity.NewFileBody("name", []byte("file"))

	folderInfo, err := entity.NewFolderInfo(nil, "name", "/path/", false)
	if err != nil {
		t.Error(err.Error())
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fileInfoRepository := mock_repository.NewMockFileInfoRepository(ctrl)
	fileInfoRepository.EXPECT().Creates(gomock.Any(), gomock.Any()).Return([]entity.FileInfo{*fileInfo}, nil)

	fileBodyRepository := mock_repository.NewMockFileBodyRepository(ctrl)
	fileBodyRepository.EXPECT().Create(gomock.Any()).Return(nil)

	folderInforepository := mock_repository.NewMockFolderInfoRepository(ctrl)
	folderInforepository.EXPECT().FindOneByID(gomock.Any(), fileInfo.FolderID).Return(folderInfo, nil)

	fileInfoService := mock_service.NewMockFileInfoService(ctrl)
	fileInfoService.EXPECT().IsExists(gomock.Any(), gomock.Any()).Return(false, nil)

	fu := NewFileUsecase(db, fileInfoRepository, fileBodyRepository, folderInforepository, fileInfoService)

	result, err := fu.Create(fileInfo.FolderID, fileInfo.IsHide, []types.File{{Name: fileInfo.Name.Value, Body: fileBody.Body}})
	if err != nil {
		t.Error(err.Error())
	}

	if result == nil {
		t.Error("failed to create file")
	}
}
