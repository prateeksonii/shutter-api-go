package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/prateeksonii/shutter-api-go/app/handlers"
)

func UserRoutes(r fiber.Router) {
	r.Get("/", handlers.IsAuthenticated, handlers.SearchUsersByUsername)
}
