package main

import (
	"fmt"
	"log"

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
		AppName:       "App Name",
	})
	// app.Use(cors.New())

	// database.ConnectDB()

	// router.SetupRoutes(app)
	log.Fatal(app.Listen(":3000", fiber.ListenConfig{EnablePrefork: true}))
}