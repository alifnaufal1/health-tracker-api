package middleware

import (
	"health-tracker-api/config"
	"health-tracker-api/pkg/context"
	"health-tracker-api/pkg/response"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func JwtProtected() fiber.Handler {
	return func(c fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return response.Error(c, fiber.StatusUnauthorized, "missing or invalid authorization header")
		}
		
		accessTokenString := strings.TrimPrefix(authHeader, "Bearer ")
		accessToken, err := jwt.Parse(accessTokenString, func(t *jwt.Token) (interface{}, error) {
			return []byte(config.Config("SECRET")), nil
		})
		if err != nil || !accessToken.Valid {
			return response.Error(c, fiber.StatusUnauthorized, "invalid or expired token")
		}
		
		claims := accessToken.Claims.(jwt.MapClaims)
		
		userID, err := uuid.Parse(claims["user_id"].(string))
		if err != nil {
			return response.Error(c, fiber.StatusUnauthorized, "invalid token payload")
		}

		context.SetAuthUser(c, &context.AuthenticatedUser{
			UserID:   userID,
			Username: claims["username"].(string),
		})

		return c.Next()
	}
}
