package auth

import (
	"health-tracker-api/pkg/helper"
	"time"
)

const (
	refreshTokenBytes = 32
	refreshTokenTTL   = 30 * 24 * time.Hour
)

func newRefreshToken() (raw, hash string, err error) {
	raw, err = helper.GenerateRandomString(refreshTokenBytes)
	if err != nil {
		return "", "", err
	}
	return raw, helper.HashToken(raw), nil
}
