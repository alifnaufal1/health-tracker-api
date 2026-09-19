package context

import "github.com/google/uuid"

type AuthenticatedUser struct {
	UserID   uuid.UUID
	Username string
}