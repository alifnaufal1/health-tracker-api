package jwt

import (
	"health-tracker-api/config"
	"health-tracker-api/internal/user"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var AccessTokenValidityPeriod int64 = time.Now().Add(time.Minute * 5).Unix()
var RefreshTokenValidityPeriod int64 = time.Now().Add(time.Minute * 5).Unix()

func GenerateAccessToken(user *user.User) (*string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	return writeToken(*token, user.ID, AccessTokenValidityPeriod)
}

func GenerateRefreshToken(user *user.User) (*string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	return writeToken(*token, user.ID, RefreshTokenValidityPeriod)
}

func writeToken(token jwt.Token, userID uuid.UUID, exp int64) (*string, error) {
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = userID
	claims["exp"] = exp

	t, err := token.SignedString([]byte(config.Config("SECRET")))
	if err != nil {
		return nil, err
	}
	return &t, nil
}
