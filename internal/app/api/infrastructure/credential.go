package infrastructure

import (
	"file-server/internal/app/api/domain/entity"
	"file-server/internal/app/api/domain/repository"
	"file-server/internal/app/api/infrastructure/model"

	"gorm.io/gorm"
)

type credentialInfrastructure struct{}

func NewCredentialInfrastructure() repository.CredentialRepository {
	return &credentialInfrastructure{}
}

func (ci *credentialInfrastructure) FindOne(db *gorm.DB) (*entity.Credential, error) {
	var credentialModel model.CredentialModel
	if err := db.First(&credentialModel).Error; err != nil {
		return nil, err
	}
	return ci.modelToEntity(&credentialModel), nil
}

func (ci *credentialInfrastructure) modelToEntity(credential *model.CredentialModel) *entity.Credential {
	credentialEntity := &entity.Credential{}
	credentialEntity.SetID(credential.ID)
	credentialEntity.SetPassword(credential.Password)
	credentialEntity.SetCreatedAt(credential.CreatedAt)
	credentialEntity.SetUpdatedAt(credential.UpdatedAt)
	return credentialEntity
}
