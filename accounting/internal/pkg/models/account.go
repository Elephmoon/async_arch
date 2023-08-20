package models

import "github.com/google/uuid"

type Account struct {
	Number       string
	PublicID     uuid.UUID
	UserID       uuid.UUID
	UserPublicID uuid.UUID
	Balance      int64
}
