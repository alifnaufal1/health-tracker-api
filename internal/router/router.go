package router

import (
	"health-tracker-api/internal/controller"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
)

func SetupRoutes(app *fiber.App, userController controller.UserController, deviceController controller.DeviceController) {
	// Middleware
	api := app.Group("/api", logger.New())
	// api.Get("/", handler.Hello)

	// // Auth
	// auth := api.Group("/auth")
	// auth.Post("/login", handler.Login)

	// User
	user := api.Group("/user")
	user.Post("/", userController.Create)
	user.Put("/:id", userController.Update)
	// user.Patch("/:id", middleware.Protected(), controller.UpdateUser)
	user.Delete("/:id", userController.Delete)
	// user.Delete("/:id", middleware.Protected(), controller.DeleteUser)
	user.Get("/:id", userController.FindByID)
	user.Get("/", userController.FindAll)
	
	// Device
	user = api.Group("/device")
	user.Post("/", deviceController.Create)
	user.Put("/:id", deviceController.Update)
	user.Delete("/:id", deviceController.Delete)
	user.Get("/:id", deviceController.FindByID)
	user.Get("/", deviceController.FindAll)
}