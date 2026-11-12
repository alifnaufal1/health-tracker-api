package router

import (
	authhandler "health-tracker-api/internal/auth/handler"
	userhandler "health-tracker-api/internal/user/handler"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
)

func SetupRoutes(
	app *fiber.App,
	userHandler *userhandler.UserHandler,
	authHandler *authhandler.AuthHandler,
	workoutDataHandler *workoutHandler.WorkoutDataHandler,
) {
	api := app.Group("/api", logger.New())

	RegisterAuthRoutes(api, authHandler)
	RegisterUserRoutes(api, userHandler)
	RegisterWorkoutDataRoutes(api, workoutDataHandler)
}
