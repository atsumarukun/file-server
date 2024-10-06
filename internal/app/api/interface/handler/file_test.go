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
	files.Write([]byte("file"))

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

	dto := []dto.FileDTO{{
		ID:        1,
		FolderID:  1,
		Name:      "name",
		Path:      "/path/name",
		MimeType:  "mime/type",
		IsHide:    false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}}

	fu := mock_usecase.NewMockFileUsecase(ctrl)
	fu.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).Return(dto, nil)

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

	dto := &dto.FileDTO{
		ID:        1,
		FolderID:  1,
		Name:      "name",
		Path:      "/path/name",
		MimeType:  "mime/type",
		IsHide:    false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

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
