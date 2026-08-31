package router

import (
	"health-tracker-api/internal/controller"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
)

func SetupRoutes(app *fiber.App, userController controller.UserController) {
	// // Middleware
	api := app.Group("/api", logger.New())
	// api.Get("/", handler.Hello)

	// // Auth
	// auth := api.Group("/auth")
	// auth.Post("/login", handler.Login)

	// User
	user := api.Group("/user")
	// user.Get("/:id", controller.GetUser)
	user.Post("/", userController.Create)
	// user.Patch("/:id", middleware.Protected(), controller.UpdateUser)
	// user.Delete("/:id", middleware.Protected(), controller.DeleteUser)

	// // Product
	// product := api.Group("/product")
	// product.Get("/", handler.GetAllProducts)
	// product.Get("/:id", handler.GetProduct)
	// product.Post("/", middleware.Protected(), handler.CreateProduct)
	// product.Delete("/:id", middleware.Protected(), handler.DeleteProduct)
}