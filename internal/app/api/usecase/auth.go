package usecase

import (
	"file-server/internal/app/api/domain/repository"
	"file-server/internal/app/api/usecase/dto"
	"file-server/internal/pkg/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthUsecase interface {
	Signin(string) (*dto.AuthDTO, error)
}

type authUsecase struct {
	db                   *gorm.DB
	credentialRepository repository.CredentialRepository
}

func NewAuthUsecase(db *gorm.DB, credentialRepository repository.CredentialRepository) AuthUsecase {
	return &authUsecase{
		db:                   db,
		credentialRepository: credentialRepository,
	}
}

func (au authUsecase) Signin(password string) (*dto.AuthDTO, error) {
	credential, err := au.credentialRepository.FindOne(au.db)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(credential.GetPassword()), []byte(password)); err != nil {
		return nil, err
	}

	claims := jwt.MapClaims{
		"exp": time.Now().Add(time.Hour * 1).Unix(),
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(config.JWT_SECRET_KEY))
	if err != nil {
		return nil, err
	}

	return dto.NewAuthDTO(token), nil
}
