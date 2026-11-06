package adapter

import (
	"health-tracker-api/pkg/context"
	"health-tracker-api/pkg/helper"
	"health-tracker-api/pkg/jwt"
	"time"

	"github.com/google/uuid"
)

type BcryptHasher struct{}

func (BcryptHasher) Hash(plain string) (string, error) { return helper.HashPassword(plain) }
func (BcryptHasher) Compare(plain, hashed string) bool {
	return helper.CheckPasswordHash(plain, hashed)
}

type AccessTokenIssuer struct{}

func (AccessTokenIssuer) Issue(userID uuid.UUID) (string, error) {
	t, err := jwt.GenerateJWT(&context.AuthenticatedUser{UserID: userID}) // package context & jwt milikmu
	if err != nil {
		return "", err
	}
	return *t, nil
}

type RefreshTokenCodec struct{}

func (RefreshTokenCodec) Generate() (string, string, error) {
	raw, err := helper.GenerateRandomString(32)
	if err != nil {
		return "", "", err
	}
	return raw, helper.HashToken(raw), nil
}

func (RefreshTokenCodec) Hash(raw string) string { return helper.HashToken(raw) }

type SystemClock struct{}

func (SystemClock) Now() time.Time { return time.Now() }
