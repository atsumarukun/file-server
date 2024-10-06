package handler

import (
	"bytes"
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
		t.Error("failed to create file")
	}
}
