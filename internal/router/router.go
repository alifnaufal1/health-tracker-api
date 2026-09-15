package router

import (
	"health-tracker-api/internal/controller"
	"health-tracker-api/pkg/middleware"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
)

func SetupRoutes(app *fiber.App, userController controller.UserController, authController controller.AuthController, deviceController controller.DeviceController, workoutDataController controller.WorkoutDataController) {
	// Middleware
	api := app.Group("/api", logger.New())

	// Auth
	auth := api.Group("/auth")
	auth.Post("/login", authController.Login)
	auth.Post("/register", authController.Register)

	// User
	user := api.Group("/user")
	user.Post("", userController.Create)
	user.Put("/:id", middleware.Protected(), userController.Update)
	user.Delete("/:id", middleware.Protected(), userController.Delete)
	user.Get("/:id", middleware.Protected(), userController.FindByID)
	user.Get("", middleware.Protected(), userController.FindAll)
	
	// Device
	user = api.Group("/device")
	user.Post("", middleware.Protected(), deviceController.Create)
	user.Put("/:id", middleware.Protected(), deviceController.Update)
	user.Delete("/:id", middleware.Protected(), deviceController.Delete)
	user.Get("/:id", middleware.Protected(), deviceController.FindByID)
	user.Get("", middleware.Protected(), deviceController.FindAll)

	// Workout Data
	user = api.Group("/workout-data")
	user.Post("/", middleware.Protected(), workoutDataController.Create)
}