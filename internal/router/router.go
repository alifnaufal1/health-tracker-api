package router

import (
	authhandler "health-tracker-api/internal/auth/handler"
	userhandler "health-tracker-api/internal/user/handler"
	"health-tracker-api/internal/workout_data"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
)

func SetupRoutes(
	app *fiber.App,
	userHandler *userhandler.UserHandler,
	authHandler *authhandler.AuthHandler,
	workoutDataHandler workout_data.WorkoutDataHandler,
) {
	api := app.Group("/api", logger.New())

	RegisterAuthRoutes(api, authHandler)
	RegisterUserRoutes(api, userHandler)
	workout_data.RegisterWorkoutDataRoutes(api, workoutDataHandler)
}
