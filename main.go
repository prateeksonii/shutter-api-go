package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/prateeksonii/shutter-api-go/pkg/routes"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	app := fiber.New()

	router := app.Group("/api/v1")

	router.Route("/users", routes.UserRoutes)

	port, hasPortEnv := os.LookupEnv("PORT")

	if !hasPortEnv {
		port = "4000"
	}

	app.Listen(":" + port)
}
