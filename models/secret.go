package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type Secret struct {
	ID               uuid.UUID `json:"hash"`
	ExpAfter         time.Time `json:"expireAfter"`
	ExpireAfterViews int       `json:"expireAfterViews"`
	Msg              string    `json:"message"`
}
