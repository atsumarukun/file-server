package entity

import "time"

type Credential struct {
	id        uint64
	password  string
	createdAt time.Time
	updatedAt time.Time
}

func NewCredential(password string) *Credential {
	credential := &Credential{}
	credential.SetPassword(password)
	return credential
}

func (c *Credential) GetID() uint64 {
	return c.id
}

func (c *Credential) SetID(id uint64) {
	c.id = id
}

func (c *Credential) GetPassword() string {
	return c.password
}

func (c *Credential) SetPassword(password string) {
	c.password = password
}

func (c *Credential) GetCreatedAt() time.Time {
	return c.createdAt
}

func (c *Credential) SetCreatedAt(createdAt time.Time) {
	c.createdAt = createdAt
}

func (c *Credential) GetUpdatedAt() time.Time {
	return c.updatedAt
}

func (c *Credential) SetUpdatedAt(updatedAt time.Time) {
	c.updatedAt = updatedAt
}
