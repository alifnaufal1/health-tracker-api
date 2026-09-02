package main

import (
	"fmt"
	"health-tracker-api/internal/controller"
	"health-tracker-api/internal/repository"
	"health-tracker-api/internal/router"
	"health-tracker-api/internal/service"
	"health-tracker-api/pkg/database"
	"health-tracker-api/pkg/middleware"
	"log"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/requestid"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		fmt.Println("Warning: no .env file found, using environment variables")
	}
	logger := logrus.New()

	app := fiber.New(fiber.Config{
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "Fiber",
		AppName:       "Health Tracker API",
		ErrorHandler: middleware.GlobalErrorHandler,
	})
	app.Use(requestid.New())

	database.ConnectDB()
	validate := validator.New()
	userRepository := repository.NewUserRepository(logger)
	userService := service.NewUserService(userRepository, validate, logger)
	userController := controller.NewUserController(userService, logger)

	router.SetupRoutes(app, userController)
	log.Fatal(app.Listen(":3108", fiber.ListenConfig{EnablePrefork: true}))
}