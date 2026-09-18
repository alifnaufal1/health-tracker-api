package jwt

import (
	"health-tracker-api/config"
	"health-tracker-api/internal/user"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(user *user.User) (*string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = user.Username
	claims["user_id"] = user.ID
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	t, err := token.SignedString([]byte(config.Config("SECRET")))
	if err != nil {
		return nil, err
	}
	return &t, nil
}