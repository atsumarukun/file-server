package usecase

import (
	"file-server/internal/app/api/domain/entity"
	"file-server/internal/pkg/types"
	"file-server/test/database"
	mock_repository "file-server/test/mock/domain/repository"
	mock_service "file-server/test/mock/domain/service"
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

	fileInfo, err := entity.NewFileInfo(1, "name", "/path/name", "mime/type", false)
	if err != nil {
		t.Error(err.Error())
	}

	fileBody := entity.NewFileBody("name", []byte("file"))

	folderInfo, err := entity.NewFolderInfo(nil, "name", "/path/name/", false)
	if err != nil {
		t.Error(err.Error())
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fileInfoRepository := mock_repository.NewMockFileInfoRepository(ctrl)
	fileInfoRepository.EXPECT().Creates(gomock.Any(), gomock.Any()).Return([]entity.FileInfo{*fileInfo}, nil)

	fileBodyRepository := mock_repository.NewMockFileBodyRepository(ctrl)
	fileBodyRepository.EXPECT().Create(gomock.Any()).Return(nil)

	folderInfoRepository := mock_repository.NewMockFolderInfoRepository(ctrl)
	folderInfoRepository.EXPECT().FindOneByID(gomock.Any(), gomock.Any()).Return(folderInfo, nil)

	fileInfoService := mock_service.NewMockFileInfoService(ctrl)
	fileInfoService.EXPECT().IsExists(gomock.Any(), gomock.Any()).Return(false, nil)

	fu := NewFileUsecase(db, fileInfoRepository, fileBodyRepository, folderInfoRepository, fileInfoService)

	result, err := fu.Create(fileInfo.FolderID, fileInfo.IsHide, []types.File{{Name: fileInfo.Name.Value, Body: fileBody.Body}})
	if err != nil {
		t.Error(err.Error())
	}

	if result == nil {
		t.Error("failed to create file")
	}
}

func TestUpdateFile(t *testing.T) {
	db, mock, err := database.Open()
	if err != nil {
		t.Error(err.Error())
	}
	mock.ExpectBegin()
	mock.ExpectCommit()

	fileInfo, err := entity.NewFileInfo(1, "name", "/path/name", "mime/type", false)
	if err != nil {
		t.Error(err.Error())
	}
	fileInfo.ID = 1

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fileInfoRepository := mock_repository.NewMockFileInfoRepository(ctrl)
	fileInfoRepository.EXPECT().FindOneByIDAndIsHide(gomock.Any(), gomock.Any(), gomock.Any()).Return(fileInfo, nil)
	fileInfoRepository.EXPECT().Update(gomock.Any(), gomock.Any()).Return(fileInfo, nil)

	fileBodyRepository := mock_repository.NewMockFileBodyRepository(ctrl)
	fileBodyRepository.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)

	folderInfoRepository := mock_repository.NewMockFolderInfoRepository(ctrl)

	fileInfoService := mock_service.NewMockFileInfoService(ctrl)
	fileInfoService.EXPECT().IsExists(gomock.Any(), gomock.Any()).Return(false, nil)

	fu := NewFileUsecase(db, fileInfoRepository, fileBodyRepository, folderInfoRepository, fileInfoService)

	result, err := fu.Update(fileInfo.ID, "update", true, false)
	if err != nil {
		t.Error(err.Error())
	}

	if result == nil {
		t.Error("failed to update file")
	}
}

func TestRemoveFile(t *testing.T) {
	db, mock, err := database.Open()
	if err != nil {
		t.Error(err.Error())
	}
	mock.ExpectBegin()
	mock.ExpectCommit()

	fileInfo, err := entity.NewFileInfo(1, "name", "/path/name", "mime/type", false)
	if err != nil {
		t.Error(err.Error())
	}
	fileInfo.ID = 1

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fileInfoRepository := mock_repository.NewMockFileInfoRepository(ctrl)
	fileInfoRepository.EXPECT().FindOneByIDAndIsHide(gomock.Any(), fileInfo.ID, fileInfo.IsHide).Return(fileInfo, nil)
	fileInfoRepository.EXPECT().Remove(gomock.Any(), gomock.Any()).Return(nil)

	fileBodyRepository := mock_repository.NewMockFileBodyRepository(ctrl)
	fileBodyRepository.EXPECT().Remove(gomock.Any()).Return(nil)

	folderInfoRepository := mock_repository.NewMockFolderInfoRepository(ctrl)

	fileInfoService := mock_service.NewMockFileInfoService(ctrl)

	fu := NewFileUsecase(db, fileInfoRepository, fileBodyRepository, folderInfoRepository, fileInfoService)

	err = fu.Remove(fileInfo.ID, false)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestMoveFile(t *testing.T) {
	db, mock, err := database.Open()
	if err != nil {
		t.Error(err.Error())
	}
	mock.ExpectBegin()
	mock.ExpectCommit()

	fileInfo, err := entity.NewFileInfo(1, "name", "/path/name", "mime/type", false)
	if err != nil {
		t.Error(err.Error())
	}
	fileInfo.ID = 1

	folderInfo, err := entity.NewFolderInfo(nil, "name", "/path/name/", false)
	if err != nil {
		t.Error(err.Error())
	}
	folderInfo.ID = 2

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fileInfoRepository := mock_repository.NewMockFileInfoRepository(ctrl)
	fileInfoRepository.EXPECT().FindOneByIDAndIsHide(gomock.Any(), gomock.Any(), gomock.Any()).Return(fileInfo, nil)
	fileInfoRepository.EXPECT().Update(gomock.Any(), gomock.Any()).Return(fileInfo, nil)

	fileBodyRepository := mock_repository.NewMockFileBodyRepository(ctrl)
	fileBodyRepository.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)

	folderInfoRepository := mock_repository.NewMockFolderInfoRepository(ctrl)
	folderInfoRepository.EXPECT().FindOneByID(gomock.Any(), gomock.Any()).Return(folderInfo, nil)

	fileInfoService := mock_service.NewMockFileInfoService(ctrl)
	fileInfoService.EXPECT().IsExists(gomock.Any(), gomock.Any()).Return(false, nil)

	fu := NewFileUsecase(db, fileInfoRepository, fileBodyRepository, folderInfoRepository, fileInfoService)

	result, err := fu.Move(fileInfo.ID, 2, false)
	if err != nil {
		t.Error(err.Error())
	}

	if result == nil {
		t.Error("failed to move file")
	}
}

func TestCopyFile(t *testing.T) {
	db, mock, err := database.Open()
	if err != nil {
		t.Error(err.Error())
	}
	mock.ExpectBegin()
	mock.ExpectCommit()

	fileInfo, err := entity.NewFileInfo(1, "name", "/path/name", "mime/type", false)
	if err != nil {
		t.Error(err.Error())
	}
	fileInfo.ID = 1

	fileBody := entity.NewFileBody("name", []byte("file"))

	folderInfo, err := entity.NewFolderInfo(nil, "name", "/path/name/", false)
	if err != nil {
		t.Error(err.Error())
	}
	folderInfo.ID = 2

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fileInfoRepository := mock_repository.NewMockFileInfoRepository(ctrl)
	fileInfoRepository.EXPECT().FindOneByIDAndIsHide(gomock.Any(), gomock.Any(), gomock.Any()).Return(fileInfo, nil)
	fileInfoRepository.EXPECT().Create(gomock.Any(), gomock.Any()).Return(fileInfo, nil)

	fileBodyRepository := mock_repository.NewMockFileBodyRepository(ctrl)
	fileBodyRepository.EXPECT().Read(gomock.Any()).Return(fileBody, nil)
	fileBodyRepository.EXPECT().Create(gomock.Any()).Return(nil)

	folderInfoRepository := mock_repository.NewMockFolderInfoRepository(ctrl)
	folderInfoRepository.EXPECT().FindOneByID(gomock.Any(), gomock.Any()).Return(folderInfo, nil)

	fileInfoService := mock_service.NewMockFileInfoService(ctrl)
	fileInfoService.EXPECT().IsExists(gomock.Any(), gomock.Any()).Return(false, nil)

	fu := NewFileUsecase(db, fileInfoRepository, fileBodyRepository, folderInfoRepository, fileInfoService)

	result, err := fu.Copy(fileInfo.ID, 2, false)
	if err != nil {
		t.Error(err.Error())
	}

	if result == nil {
		t.Error("failed to copy file")
	}
}

func CreateReadFile(t *testing.T) {
	db, mock, err := database.Open()
	if err != nil {
		t.Error(err.Error())
	}
	mock.ExpectBegin()
	mock.ExpectCommit()

	fileInfo, err := entity.NewFileInfo(1, "name", "/path/name", "mime/type", false)
	if err != nil {
		t.Error(err.Error())
	}
	fileInfo.ID = 1

	fileBody := entity.NewFileBody("name", []byte("file"))

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fileInfoRepository := mock_repository.NewMockFileInfoRepository(ctrl)
	fileInfoRepository.EXPECT().FindOneByIDAndIsHide(gomock.Any(), gomock.Any(), gomock.Any()).Return(fileInfo, nil)

	fileBodyRepository := mock_repository.NewMockFileBodyRepository(ctrl)
	fileBodyRepository.EXPECT().Read(gomock.Any()).Return(fileBody, nil)

	folderInfoRepository := mock_repository.NewMockFolderInfoRepository(ctrl)

	fileInfoService := mock_service.NewMockFileInfoService(ctrl)

	fu := NewFileUsecase(db, fileInfoRepository, fileBodyRepository, folderInfoRepository, fileInfoService)

	result, err := fu.Read(fileInfo.ID, false)
	if err != nil {
		t.Error(err.Error())
	}

	if result == nil {
		t.Error("failed to read file")
	}
}
