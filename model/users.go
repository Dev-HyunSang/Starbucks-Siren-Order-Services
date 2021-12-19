package model

import (
	"time"

	"github.com/twinj/uuid"
)

type Users struct {
	UUID      uuid.UUID  `json:"uuid"`
	Name      string     `json:"name"`
	NickName  string     `json:"nickname"`
	Birthday  *time.Time `json:"birthday"`
	Email     string     `json:"email"`
	Password  string     `json:"password"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}
