package workout_data

import (
	"health-tracker-api/pkg/middleware"

	"github.com/gofiber/fiber/v3"
)

func RegisterWorkoutDataRoutes(router fiber.Router, handler WorkoutDataHandler) {
	router.Get("/devices/:deviceId/workout-data", middleware.JwtProtected(), handler.GetByDeviceID)
	router.Post("/workout-data", middleware.JwtProtected(), handler.Create)
	// workoutData.Put("/:id", middleware.Protected(), handler.Update)
	// workoutData.Delete("/:id", middleware.Protected(), handler.Delete)
	// workoutData.Get("/:id", middleware.Protected(), handler.FindByID)
	// workoutData.Get("", middleware.Protected(), handler.FindAll)
}