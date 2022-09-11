package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"github.com/prateeksonii/shutter-api-go/pkg/configs"
	"github.com/prateeksonii/shutter-api-go/pkg/routes"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	configs.Connect()

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			status := c.Response().StatusCode()

			if status == fiber.StatusOK {
				status = fiber.StatusInternalServerError
			}

			return c.Status(status).JSON(fiber.Map{
				"error":   true,
				"message": err.Error(),
			})
		},
	})

	app.Use(cors.New())

	router := app.Group("/api/v1")

	router.Route("/auth", routes.AuthRoutes)

	port, hasPortEnv := os.LookupEnv("PORT")

	if !hasPortEnv {
		port = "4000"
	}

	app.Listen(":" + port)
}
