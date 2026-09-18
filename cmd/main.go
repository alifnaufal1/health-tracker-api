package main

import (
	"fmt"
	"health-tracker-api/internal/auth"
	"health-tracker-api/internal/device"
	"health-tracker-api/internal/router"
	"health-tracker-api/internal/user"
	"health-tracker-api/internal/workout_data"
	"health-tracker-api/pkg/database"
	"health-tracker-api/pkg/middleware"
	"log"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
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
	app.Use(cors.New())

	database.ConnectDB()
	database.Migrate(&user.User{}, &device.Device{}, &workout_data.WorkoutData{})
	validate := validator.New()

	
	userRepository := user.NewUserRepository(logger)
	userService := user.NewUserService(userRepository, validate, logger)
	userController := user.NewUserHandler(userService, logger)
	
	authService := auth.NewAuthService(userService, userRepository, validate, logger)
	authController := auth.NewAuthHandler(authService, logger)
	
	deviceRepository := device.NewDeviceRepository(logger)
	deviceService := device.NewDeviceService(deviceRepository, validate, logger)
	deviceController := device.NewDeviceHandler(deviceService, logger)

	workoutDataRepository := workout_data.NewWorkoutDataRepository(logger)
	workoutDataService := workout_data.NewWorkoutDataService(workoutDataRepository, validate, logger)
	workoutDataController := workout_data.NewWorkoutDataHandler(workoutDataService)

	router.SetupRoutes(app, userController, authController, deviceController, workoutDataController)
	log.Fatal(app.Listen(":3108", fiber.ListenConfig{EnablePrefork: true}))
}