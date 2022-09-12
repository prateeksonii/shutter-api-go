package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/prateeksonii/shutter-api-go/app/handlers"
)

func InviteRoutes(r fiber.Router) {
	r.Post("/:username", handlers.IsAuthenticated, handlers.SendInvite)
}
