package model

import "time"

type CredentialModel struct {
	ID        uint64
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (cm *CredentialModel) TableName() string {
	return "credentials"
}
