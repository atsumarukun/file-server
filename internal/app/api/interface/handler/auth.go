package handler

import (
	"errors"
	"file-server/internal/app/api/interface/requests"
	"file-server/internal/app/api/interface/responses"
	"file-server/internal/app/api/usecase"
	"file-server/internal/app/api/usecase/dto"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthHandler interface {
	Signin(*gin.Context)
}

type authHandler struct {
	usecase usecase.AuthUsecase
}

func NewAuthHandler(usecase usecase.AuthUsecase) AuthHandler {
	return &authHandler{
		usecase: usecase,
	}
}

func (ah *authHandler) Signin(c *gin.Context) {
	var request requests.SigninRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	dto, err := ah.usecase.Signin(request.Password)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.String(http.StatusNotFound, err.Error())
		} else {
			c.String(http.StatusInternalServerError, err.Error())
		}
		return
	}

	c.JSON(http.StatusOK, ah.dtoToResponse(dto))
}

func (ah *authHandler) dtoToResponse(auth *dto.AuthDTO) *responses.AuthResponse {
	return &responses.AuthResponse{
		Token: auth.Token,
	}
}
