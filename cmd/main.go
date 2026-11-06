package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"health-tracker-api/internal/auth"
	authhandler "health-tracker-api/internal/auth/handler"
	"health-tracker-api/internal/device"
	refreshtokengorm "health-tracker-api/internal/refresh_token/gormrepo"
	"health-tracker-api/internal/router"
	"health-tracker-api/internal/user"
	usergorm "health-tracker-api/internal/user/gormrepo"
	userhandler "health-tracker-api/internal/user/handler"
	"health-tracker-api/internal/workout_data"
	"health-tracker-api/pkg/adapter"
	"health-tracker-api/pkg/database"
	"health-tracker-api/pkg/database/gormtx"
	"health-tracker-api/pkg/logger"
	"health-tracker-api/pkg/logging/logrusx"
	"health-tracker-api/pkg/middleware"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/requestid"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		fmt.Println("Warning: no .env file found, using environment variables")
	}
	appLog := logger.NewLogger()

	app := fiber.New(fiber.Config{
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "Fiber",
		AppName:       "Health Tracker API",
		ErrorHandler:  middleware.GlobalErrorHandler,
	})
	app.Use(requestid.New())
	app.Use(cors.New())

	database.ConnectDB()
	if !fiber.IsChild() {
		database.Migrate(
			&usergorm.Model{},
			&device.Device{},
			&workout_data.WorkoutData{},
			&refreshtokengorm.Model{},
		)
	}

	validate := validator.New()
	txm := gormtx.NewManager(database.DB)
	lg := logrusx.New(appLog)

	// --- user ---
	userRepo := usergorm.New(txm)
	userService := user.NewUserService(userRepo, txm, adapter.BcryptHasher{}, lg)
	userHandler := userhandler.NewUserHandler(userService, validate)

	// --- auth ---
	authService := auth.NewAuthService(auth.Dependencies{
		Users:         userService,
		RefreshTokens: refreshtokengorm.New(txm),
		Tx:            txm,
		AccessTokens:  adapter.AccessTokenIssuer{},
		RefreshCodec:  adapter.RefreshTokenCodec{},
		Clock:         adapter.SystemClock{},
		Log:           lg,
	}, auth.Config{RefreshTokenTTL: envDuration("REFRESH_TOKEN_TTL", 30*24*time.Hour)})
	authHandler := authhandler.NewAuthHandler(authService, validate)

	// --- device & workout data ---
	deviceRepository := device.NewDeviceRepository()
	deviceService := device.NewDeviceService(deviceRepository, validate, appLog)

	workoutDataRepository := workout_data.NewWorkoutDataRepository()
	workoutDataService := workout_data.NewWorkoutDataService(workoutDataRepository, validate, appLog, deviceService)
	workoutDataHandler := workout_data.NewWorkoutDataHandler(workoutDataService)

	router.SetupRoutes(app, userHandler, authHandler, workoutDataHandler)

	log.Fatal(app.Listen(":3108", fiber.ListenConfig{
		EnablePrefork: os.Getenv("PREFORK") == "true",
	}))
}

func envDuration(key string, fallback time.Duration) time.Duration {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		log.Printf("invalid %s=%q, using default %s", key, v, fallback)
		return fallback
	}
	return d
}
