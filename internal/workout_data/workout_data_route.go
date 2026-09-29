package workout_data

import (
	"health-tracker-api/pkg/middleware"

	"github.com/gofiber/fiber/v3"
)

func RegisterWorkoutDataRoutes(router fiber.Router, handler WorkoutDataHandler) {
	router.Get("/devices/:deviceId/workout-data", middleware.JwtProtected(), handler.GetByDeviceID)
	router.Post("/devices/:deviceId/workout-data", middleware.JwtProtected(), handler.Create)

	workout_data := router.Group("/workout-data")

	workout_data.Get("/:id", middleware.JwtProtected(), handler.GetByID)
	workout_data.Delete("/:id", middleware.JwtProtected())
}
