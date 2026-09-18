package router

import (
	"health-tracker-api/internal/auth"
	"health-tracker-api/internal/device"
	"health-tracker-api/internal/user"
	"health-tracker-api/internal/workout_data"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
)

func SetupRoutes(app *fiber.App, userHandler user.UserHandler, authHandler auth.AuthHandler, deviceHandler device.DeviceHandler, workoutDataHandler workout_data.WorkoutDataHandler) {
	// Middleware
	api := app.Group("/api", logger.New())

	auth.RegisterAuthRoutes(api, authHandler)
	user.RegisterUserRoutes(api, userHandler)
	device.RegisterDeviceRoutes(api, deviceHandler)
	workout_data.RegisterWorkoutDataRoutes(api, workoutDataHandler)
}