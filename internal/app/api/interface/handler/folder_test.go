package handler

import (
	"bytes"
	"encoding/json"
	"file-server/internal/app/api/interface/requests"
	"file-server/internal/app/api/usecase/dto"
	mock_usecase "file-server/test/mock/usecase"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
)

func TestCreateFolder(t *testing.T) {
	gin.SetMode(gin.TestMode)

	input := requests.CreateFolderRequest{
		ParentFolderID: 1,
		Name:           "name",
		IsHide:         false,
	}

	body, err := json.Marshal(input)
	if err != nil {
		t.Error(err.Error())
	}

	req, err := http.NewRequest("POST", "/folders", bytes.NewBuffer(body))
	if err != nil {
		t.Error(err.Error())
	}

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dto := &dto.FolderDTO{
		ID:             1,
		ParentFolderID: nil,
		Name:           "name",
		Path:           "/path/name/",
		IsHide:         false,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	fu := mock_usecase.NewMockFolderUsecase(ctrl)
	fu.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).Return(dto, nil)

	fh := NewFolderHandler(fu)

	fh.Create(ctx)

	if w.Code != http.StatusOK {
		t.Error(w.Body.String())
	}
}

func TestUpdateFolder(t *testing.T) {
	gin.SetMode(gin.TestMode)

	input := requests.UpdateFolderRequest{
		Name:   "name",
		IsHide: false,
	}

	body, err := json.Marshal(input)
	if err != nil {
		t.Error(err.Error())
	}

	req, err := http.NewRequest("PUT", "/folders/1", bytes.NewBuffer(body))
	if err != nil {
		t.Error(err.Error())
	}

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req
	ctx.Params = append(ctx.Params, gin.Param{Key: "id", Value: strconv.Itoa(1)})

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dto := &dto.FolderDTO{
		ID:             1,
		ParentFolderID: nil,
		Name:           "name",
		Path:           "/path/name/",
		IsHide:         false,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	fu := mock_usecase.NewMockFolderUsecase(ctrl)
	fu.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(dto, nil)

	fh := NewFolderHandler(fu)

	fh.Update(ctx)

	if w.Code != http.StatusOK {
		t.Error(w.Body.String())
	}
}

func TestRemoveFolder(t *testing.T) {
	gin.SetMode(gin.TestMode)

	req, err := http.NewRequest("DELETE", "/folders/1", nil)
	if err != nil {
		t.Error(err.Error())
	}

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req
	ctx.Params = append(ctx.Params, gin.Param{Key: "id", Value: strconv.Itoa(1)})

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fu := mock_usecase.NewMockFolderUsecase(ctrl)
	fu.EXPECT().Remove(gomock.Any(), gomock.Any()).Return(nil)

	fh := NewFolderHandler(fu)

	fh.Remove(ctx)

	if w.Code != http.StatusOK {
		t.Error(w.Body.String())
	}
}

func TestMoveFolder(t *testing.T) {
	gin.SetMode(gin.TestMode)

	input := requests.MoveFolderRequest{
		ParentFolderID: 1,
	}

	body, err := json.Marshal(input)
	if err != nil {
		t.Error(err.Error())
	}

	req, err := http.NewRequest("PUT", "/folders/1/move", bytes.NewBuffer(body))
	if err != nil {
		t.Error(err.Error())
	}

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req
	ctx.Params = append(ctx.Params, gin.Param{Key: "id", Value: strconv.Itoa(1)})

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dto := &dto.FolderDTO{
		ID:             1,
		ParentFolderID: nil,
		Name:           "name",
		Path:           "/path/name/",
		IsHide:         false,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	fu := mock_usecase.NewMockFolderUsecase(ctrl)
	fu.EXPECT().Move(gomock.Any(), gomock.Any(), gomock.Any()).Return(dto, nil)

	fh := NewFolderHandler(fu)

	fh.Move(ctx)

	if w.Code != http.StatusOK {
		t.Error(w.Body.String())
	}
}

func TestCopyFolder(t *testing.T) {
	gin.SetMode(gin.TestMode)

	input := requests.CopyFolderRequest{
		ParentFolderID: 1,
	}

	body, err := json.Marshal(input)
	if err != nil {
		t.Error(err.Error())
	}

	req, err := http.NewRequest("PUT", "/folders/1/copy", bytes.NewBuffer(body))
	if err != nil {
		t.Error(err.Error())
	}

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req
	ctx.Params = append(ctx.Params, gin.Param{Key: "id", Value: strconv.Itoa(1)})

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dto := &dto.FolderDTO{
		ID:             1,
		ParentFolderID: nil,
		Name:           "name",
		Path:           "/path/name/",
		IsHide:         false,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	fu := mock_usecase.NewMockFolderUsecase(ctrl)
	fu.EXPECT().Copy(gomock.Any(), gomock.Any(), gomock.Any()).Return(dto, nil)

	fh := NewFolderHandler(fu)

	fh.Copy(ctx)

	if w.Code != http.StatusOK {
		t.Error(w.Body.String())
	}
}
