package handler

import (
	"bytes"
	"encoding/json"
	"file-server/internal/app/api/interface/requests"
	"file-server/internal/app/api/usecase/dto"
	mock_usecase "file-server/test/mock/usecase"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
)

func TestCreateFile(t *testing.T) {
	gin.SetMode(gin.TestMode)

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	folderID, err := writer.CreateFormField("folder_id")
	if err != nil {
		t.Error(err.Error())
	}
	if _, err := folderID.Write([]byte(strconv.Itoa(1))); err != nil {
		t.Error(err.Error())
	}
	isHide, err := writer.CreateFormField("is_hide")
	if err != nil {
		t.Error(err.Error())
	}
	if _, err := isHide.Write([]byte(strconv.FormatBool(false))); err != nil {
		t.Error(err.Error())
	}
	files, err := writer.CreateFormFile("files[]", "file")
	if err != nil {
		t.Error(err.Error())
	}
	if _, err := files.Write([]byte("file")); err != nil {
		t.Error(err.Error())
	}

	writer.Close()

	req, err := http.NewRequest("POST", "/files", body)
	if err != nil {
		t.Error(err.Error())
	}
	req.Header.Add("Content-Type", writer.FormDataContentType())

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dtos := []dto.FileInfoDTO{*dto.NewFileInfoDTO(1, 1, "name", "path/name", "mime/type", false, time.Now(), time.Now())}

	fu := mock_usecase.NewMockFileUsecase(ctrl)
	fu.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).Return(dtos, nil)

	fh := NewFileHandler(fu)

	fh.Create(ctx)

	if w.Code != http.StatusOK {
		t.Error(w.Body.String())
	}
}

func TestUpdateFile(t *testing.T) {
	gin.SetMode(gin.TestMode)

	input := requests.UpdateFileRequest{
		Name:   "name",
		IsHide: false,
	}

	body, err := json.Marshal(input)
	if err != nil {
		t.Error(err.Error())
	}

	req, err := http.NewRequest("PUT", "/files/1", bytes.NewBuffer(body))
	if err != nil {
		t.Error(err.Error())
	}

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req
	ctx.Params = append(ctx.Params, gin.Param{Key: "id", Value: strconv.Itoa(1)})

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dto := dto.NewFileInfoDTO(1, 1, "name", "path/name", "mime/type", false, time.Now(), time.Now())

	fu := mock_usecase.NewMockFileUsecase(ctrl)
	fu.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(dto, nil)

	fh := NewFileHandler(fu)

	fh.Update(ctx)

	if w.Code != http.StatusOK {
		t.Error(w.Body.String())
	}
}

func TestRemoveFile(t *testing.T) {
	gin.SetMode(gin.TestMode)

	req, err := http.NewRequest("DELETE", "/files/1", nil)
	if err != nil {
		t.Error(err.Error())
	}

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req
	ctx.Params = append(ctx.Params, gin.Param{Key: "id", Value: strconv.Itoa(1)})

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fu := mock_usecase.NewMockFileUsecase(ctrl)
	fu.EXPECT().Remove(gomock.Any(), gomock.Any()).Return(nil)

	fh := NewFileHandler(fu)

	fh.Remove(ctx)

	if w.Code != http.StatusOK {
		t.Error(w.Body.String())
	}
}

func TestMoveFile(t *testing.T) {
	gin.SetMode(gin.TestMode)

	input := requests.MoveFileRequest{
		FolderID: 1,
	}

	body, err := json.Marshal(input)
	if err != nil {
		t.Error(err.Error())
	}

	req, err := http.NewRequest("PUT", "/files/1/move", bytes.NewBuffer(body))
	if err != nil {
		t.Error(err.Error())
	}

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req
	ctx.Params = append(ctx.Params, gin.Param{Key: "id", Value: strconv.Itoa(1)})

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dto := dto.NewFileInfoDTO(1, 1, "name", "path/name", "mime/type", false, time.Now(), time.Now())

	fu := mock_usecase.NewMockFileUsecase(ctrl)
	fu.EXPECT().Move(gomock.Any(), gomock.Any(), gomock.Any()).Return(dto, nil)

	fh := NewFileHandler(fu)

	fh.Move(ctx)

	if w.Code != http.StatusOK {
		t.Error(w.Body.String())
	}
}

func TestCopyFile(t *testing.T) {
	gin.SetMode(gin.TestMode)

	input := requests.CopyFileRequest{
		FolderID: 1,
	}

	body, err := json.Marshal(input)
	if err != nil {
		t.Error(err.Error())
	}

	req, err := http.NewRequest("PUT", "/files/1/copy", bytes.NewBuffer(body))
	if err != nil {
		t.Error(err.Error())
	}

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req
	ctx.Params = append(ctx.Params, gin.Param{Key: "id", Value: strconv.Itoa(1)})

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dto := dto.NewFileInfoDTO(1, 1, "name", "path/name", "mime/type", false, time.Now(), time.Now())

	fu := mock_usecase.NewMockFileUsecase(ctrl)
	fu.EXPECT().Copy(gomock.Any(), gomock.Any(), gomock.Any()).Return(dto, nil)

	fh := NewFileHandler(fu)

	fh.Copy(ctx)

	if w.Code != http.StatusOK {
		t.Error(w.Body.String())
	}
}

func TestReadFile(t *testing.T) {
	gin.SetMode(gin.TestMode)

	req, err := http.NewRequest("PUT", "/files/1/body", nil)
	if err != nil {
		t.Error(err.Error())
	}

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req
	ctx.Params = append(ctx.Params, gin.Param{Key: "id", Value: strconv.Itoa(1)})

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dto := dto.NewFileBodyDTO("mime/type", []byte("file"))

	fu := mock_usecase.NewMockFileUsecase(ctrl)
	fu.EXPECT().Read(gomock.Any(), gomock.Any()).Return(dto, nil)

	fh := NewFileHandler(fu)

	fh.Read(ctx)

	if w.Code != http.StatusOK {
		t.Error(w.Body.String())
	}
}
