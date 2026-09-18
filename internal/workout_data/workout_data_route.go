package workout_data

import (
	"github.com/gofiber/fiber/v3"
)

func RegisterWorkoutDataRoutes(api fiber.Router, handler WorkoutDataHandler) {
	workoutData := api.Group("/workout_data")
	workoutData.Post("", handler.Create)
	// workoutData.Put("/:id", middleware.Protected(), handler.Update)
	// workoutData.Delete("/:id", middleware.Protected(), handler.Delete)
	// workoutData.Get("/:id", middleware.Protected(), handler.FindByID)
	// workoutData.Get("", middleware.Protected(), handler.FindAll)
}