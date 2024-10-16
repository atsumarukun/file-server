package usecase

import (
	"file-server/internal/app/api/domain/entity"
	"file-server/test/database"
	mock_repository "file-server/test/mock/domain/repository"
	mock_service "file-server/test/mock/domain/service"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestCreateFolder(t *testing.T) {
	db, mock, err := database.Open()
	if err != nil {
		t.Error(err.Error())
	}
	mock.ExpectBegin()
	mock.ExpectCommit()

	folderInfo, err := entity.NewFolderInfo(nil, "name", "/path/name/", false)
	if err != nil {
		t.Error(err.Error())
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	folderInfoRepository := mock_repository.NewMockFolderInfoRepository(ctrl)
	folderInfoRepository.EXPECT().FindOneByID(gomock.Any(), gomock.Any()).Return(folderInfo, nil)
	folderInfoRepository.EXPECT().Create(gomock.Any(), gomock.Any()).Return(folderInfo, nil)

	folderBodyRepository := mock_repository.NewMockFolderBodyRepository(ctrl)
	folderBodyRepository.EXPECT().Create(gomock.Any()).Return(nil)

	folderInfoService := mock_service.NewMockFolderInfoService(ctrl)
	folderInfoService.EXPECT().IsExists(gomock.Any(), gomock.Any()).Return(false, nil)

	fu := NewFolderUsecase(db, folderInfoRepository, folderBodyRepository, folderInfoService)

	result, err := fu.Create(1, folderInfo.Name.Value, folderInfo.IsHide)
	if err != nil {
		t.Error(err.Error())
	}

	if result == nil {
		t.Error("failed to create folder")
	}
}

func TestUpdateFolder(t *testing.T) {
	db, mock, err := database.Open()
	if err != nil {
		t.Error(err.Error())
	}
	mock.ExpectBegin()
	mock.ExpectCommit()

	folderInfo, err := entity.NewFolderInfo(nil, "name", "/path/name/", false)
	if err != nil {
		t.Error(err.Error())
	}
	folderInfo.ID = 2
	var parentFolderID uint64 = 1
	folderInfo.ParentFolderID = &parentFolderID

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	folderInfoRepository := mock_repository.NewMockFolderInfoRepository(ctrl)
	folderInfoRepository.EXPECT().FindOneByIDAndIsHideWithLower(gomock.Any(), gomock.Any(), gomock.Any()).Return(folderInfo, nil)
	folderInfoRepository.EXPECT().Update(gomock.Any(), gomock.Any()).Return(folderInfo, nil)

	folderBodyRepository := mock_repository.NewMockFolderBodyRepository(ctrl)
	folderBodyRepository.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)

	folderInfoService := mock_service.NewMockFolderInfoService(ctrl)
	folderInfoService.EXPECT().IsExists(gomock.Any(), gomock.Any()).Return(false, nil)

	fu := NewFolderUsecase(db, folderInfoRepository, folderBodyRepository, folderInfoService)

	result, err := fu.Update(folderInfo.ID, "update", false, false)
	if err != nil {
		t.Error(err.Error())
	}

	if result == nil {
		t.Error("failed to update folder")
	}
}

func TestRemoveFolder(t *testing.T) {
	db, mock, err := database.Open()
	if err != nil {
		t.Error(err.Error())
	}
	mock.ExpectBegin()
	mock.ExpectCommit()

	folderInfo, err := entity.NewFolderInfo(nil, "name", "/path/name/", false)
	if err != nil {
		t.Error(err.Error())
	}
	folderInfo.ID = 2
	var parentFolderID uint64 = 1
	folderInfo.ParentFolderID = &parentFolderID

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	folderInfoRepository := mock_repository.NewMockFolderInfoRepository(ctrl)
	folderInfoRepository.EXPECT().FindOneByIDAndIsHideWithLower(gomock.Any(), gomock.Any(), gomock.Any()).Return(folderInfo, nil)
	folderInfoRepository.EXPECT().Remove(gomock.Any(), gomock.Any()).Return(nil)

	folderBodyRepository := mock_repository.NewMockFolderBodyRepository(ctrl)
	folderBodyRepository.EXPECT().Remove(gomock.Any()).Return(nil)

	folderInfoService := mock_service.NewMockFolderInfoService(ctrl)

	fu := NewFolderUsecase(db, folderInfoRepository, folderBodyRepository, folderInfoService)

	err = fu.Remove(folderInfo.ID, false)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestMoveFolder(t *testing.T) {
	db, mock, err := database.Open()
	if err != nil {
		t.Error(err.Error())
	}
	mock.ExpectBegin()
	mock.ExpectCommit()

	folderInfo, err := entity.NewFolderInfo(nil, "name", "/path/name/", false)
	if err != nil {
		t.Error(err.Error())
	}
	folderInfo.ID = 2
	var parentFolderID uint64 = 1
	folderInfo.ParentFolderID = &parentFolderID

	parentFolderInfo, err := entity.NewFolderInfo(nil, "root", "/", false)
	if err != nil {
		t.Error(err.Error())
	}
	folderInfo.ID = 1

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	folderInfoRepository := mock_repository.NewMockFolderInfoRepository(ctrl)
	folderInfoRepository.EXPECT().FindOneByID(gomock.Any(), gomock.Any()).Return(parentFolderInfo, nil)
	folderInfoRepository.EXPECT().FindOneByIDAndIsHideWithLower(gomock.Any(), gomock.Any(), gomock.Any()).Return(folderInfo, nil)
	folderInfoRepository.EXPECT().Update(gomock.Any(), gomock.Any()).Return(folderInfo, nil)

	folderBodyRepository := mock_repository.NewMockFolderBodyRepository(ctrl)
	folderBodyRepository.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)

	folderInfoService := mock_service.NewMockFolderInfoService(ctrl)
	folderInfoService.EXPECT().IsExists(gomock.Any(), gomock.Any()).Return(false, nil)

	fu := NewFolderUsecase(db, folderInfoRepository, folderBodyRepository, folderInfoService)

	result, err := fu.Move(folderInfo.ID, 1, false)
	if err != nil {
		t.Error(err.Error())
	}

	if result == nil {
		t.Error("failed to move folder")
	}
}

func TestCopyFolder(t *testing.T) {
	db, mock, err := database.Open()
	if err != nil {
		t.Error(err.Error())
	}
	mock.ExpectBegin()
	mock.ExpectCommit()

	folderInfo, err := entity.NewFolderInfo(nil, "name", "/path/name/", false)
	if err != nil {
		t.Error(err.Error())
	}
	folderInfo.ID = 2
	var parentFolderID uint64 = 1
	folderInfo.ParentFolderID = &parentFolderID

	parentFolderInfo, err := entity.NewFolderInfo(nil, "root", "/", false)
	if err != nil {
		t.Error(err.Error())
	}
	folderInfo.ID = 1

	folderBody := entity.NewFolderBody("/path/name/")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	folderInfoRepository := mock_repository.NewMockFolderInfoRepository(ctrl)
	folderInfoRepository.EXPECT().FindOneByID(gomock.Any(), gomock.Any()).Return(parentFolderInfo, nil)
	folderInfoRepository.EXPECT().FindOneByIDAndIsHideWithLower(gomock.Any(), gomock.Any(), gomock.Any()).Return(folderInfo, nil)
	folderInfoRepository.EXPECT().Create(gomock.Any(), gomock.Any()).Return(folderInfo, nil)

	folderBodyRepository := mock_repository.NewMockFolderBodyRepository(ctrl)
	folderBodyRepository.EXPECT().Read(gomock.Any()).Return(folderBody, nil)
	folderBodyRepository.EXPECT().Create(gomock.Any()).Return(nil)

	folderInfoService := mock_service.NewMockFolderInfoService(ctrl)
	folderInfoService.EXPECT().IsExists(gomock.Any(), gomock.Any()).Return(false, nil)

	fu := NewFolderUsecase(db, folderInfoRepository, folderBodyRepository, folderInfoService)

	result, err := fu.Copy(folderInfo.ID, 1, false)
	if err != nil {
		t.Error(err.Error())
	}

	if result == nil {
		t.Error("failed to copy folder")
	}
}

func TestFindOneFolder(t *testing.T) {
	db, _, err := database.Open()
	if err != nil {
		t.Error(err.Error())
	}

	folderInfo, err := entity.NewFolderInfo(nil, "name", "/path/name/", false)
	if err != nil {
		t.Error(err.Error())
	}
	folderInfo.ID = 1

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	folderInfoRepository := mock_repository.NewMockFolderInfoRepository(ctrl)
	folderInfoRepository.EXPECT().FindOneByPathAndIsHideWithChildren(gomock.Any(), gomock.Any(), gomock.Any()).Return(folderInfo, nil)

	folderBodyRepository := mock_repository.NewMockFolderBodyRepository(ctrl)

	folderInfoService := mock_service.NewMockFolderInfoService(ctrl)

	fu := NewFolderUsecase(db, folderInfoRepository, folderBodyRepository, folderInfoService)

	result, err := fu.FindOne(folderInfo.Path.Value, false)
	if err != nil {
		t.Error(err.Error())
	}

	if result == nil {
		t.Error("failed to find one folder")
	}
}

func TestReadFolder(t *testing.T) {
	db, _, err := database.Open()
	if err != nil {
		t.Error(err.Error())
	}

	folderInfo, err := entity.NewFolderInfo(nil, "name", "/path/name/", false)
	if err != nil {
		t.Error(err.Error())
	}
	folderInfo.ID = 1

	folderBody := entity.NewFolderBody("/path/name/")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	folderInfoRepository := mock_repository.NewMockFolderInfoRepository(ctrl)
	folderInfoRepository.EXPECT().FindOneByIDAndIsHideWithLower(gomock.Any(), gomock.Any(), gomock.Any()).Return(folderInfo, nil)

	folderBodyRepository := mock_repository.NewMockFolderBodyRepository(ctrl)
	folderBodyRepository.EXPECT().Read(gomock.Any()).Return(folderBody, nil)

	folderInfoService := mock_service.NewMockFolderInfoService(ctrl)

	fu := NewFolderUsecase(db, folderInfoRepository, folderBodyRepository, folderInfoService)

	result, err := fu.Read(folderInfo.ID, false)
	if err != nil {
		t.Error(err.Error())
	}

	if result == nil {
		t.Error("failed to read folder")
	}
}
