package usecase

import (
	"file-server/internal/app/api/domain/entity"
	"file-server/test/database"
	mock_repository "file-server/test/mock/domain/repository"
	"testing"

	"github.com/golang/mock/gomock"
	"golang.org/x/crypto/bcrypt"
)

func TestSignin(t *testing.T) {
	db, _, err := database.Open()
	if err != nil {
		t.Error(err.Error())
	}

	hash, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	if err != nil {
		t.Error(err.Error())
	}
	credential := entity.NewCredential(string(hash))

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mock_repository.NewMockCredentialRepository(ctrl)
	repo.EXPECT().FindOne(db).Return(credential, err)

	au := NewAuthUsecase(db, repo)
	result, err := au.Signin("password")
	if err != nil {
		t.Error(err.Error())
	}

	if result == nil {
		t.Error("failed to signin")
	}
}
