package main

import (
	"fmt"
	"health-tracker-api/internal/controller"
	"health-tracker-api/internal/database"
	"health-tracker-api/internal/repository"
	"health-tracker-api/internal/router"
	"health-tracker-api/internal/service"
	"log"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v3"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		fmt.Println("Warning: no .env file found, using environment variables")
	}

	app := fiber.New(fiber.Config{
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "Fiber",
		AppName:       "Health Tracker API",
	})
	// app.Use(cors.New())

	database.ConnectDB()
	validate := validator.New()
	userRepository := repository.NewUserRepository()
	userService := service.NewUserService(userRepository, validate)
	userController := controller.NewUserController(userService)

	router.SetupRoutes(app, userController)
	log.Fatal(app.Listen(":3108", fiber.ListenConfig{EnablePrefork: true}))
}