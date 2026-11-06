package jwt

import (
	"health-tracker-api/config"
	"health-tracker-api/pkg/context"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(user *context.AuthenticatedUser) (*string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = user.Username
	claims["user_id"] = user.UserID
	claims["exp"] = time.Now().Add(time.Minute * 15).Unix()

	t, err := token.SignedString([]byte(config.Config("SECRET")))
	if err != nil {
		return nil, err
	}
	return &t, nil
}
