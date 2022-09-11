package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/prateeksonii/shutter-api-go/app/handlers"
)

func AuthRoutes(r fiber.Router) {
	r.Post("/signup", handlers.SignUp)
	r.Post("/signin", handlers.SignIn)
	r.Get("/me", handlers.IsAuthenticated, handlers.GetAuthenticatedUser)
}
