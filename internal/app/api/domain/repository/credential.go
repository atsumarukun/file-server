package repository

import (
	"file-server/internal/app/api/domain/entity"

	"gorm.io/gorm"
)

type CredentialRepository interface {
	FindOne(*gorm.DB) (*entity.Credential, error)
}
