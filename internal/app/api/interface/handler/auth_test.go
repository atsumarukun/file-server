package handler

import (
	"bytes"
	"encoding/json"
	"file-server/internal/app/api/interface/requests"
	"file-server/internal/app/api/usecase/dto"
	mock_usecase "file-server/test/mock/usecase"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
)

func TestSignin(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dto := dto.NewAuthDTO("token")

	au := mock_usecase.NewMockAuthUsecase(ctrl)
	au.EXPECT().Signin(gomock.Any()).Return(dto, nil)

	input := requests.SigninRequest{
		Password: "password",
	}

	body, err := json.Marshal(input)
	if err != nil {
		t.Error(err.Error())
	}

	req, err := http.NewRequest("POST", "/auth/signin", bytes.NewBuffer(body))
	if err != nil {
		t.Error(err.Error())
	}

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	ah := NewAuthHandler(au)

	ah.Signin(ctx)

	if w.Code != http.StatusOK {
		t.Error("failed to signin")
	}
}
