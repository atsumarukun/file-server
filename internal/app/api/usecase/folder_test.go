package usecase

import (
	"file-server/internal/app/api/domain/entity"
	"file-server/test/database"
	mock_repository "file-server/test/mock/repository"
	mock_service "file-server/test/mock/service"
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
