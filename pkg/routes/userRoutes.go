package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/prateeksonii/shutter-api-go/pkg/utils"
)

type UserDto struct {
	Name     string `validate:"required"`
	Username string `validate:"required,min=4"`
	Password string `validate:"required,min=4"`
}

func UserRoutes(r fiber.Router) {
	r.Post("/", func(c *fiber.Ctx) error {

		userDto := &UserDto{}

		if err := c.BodyParser(userDto); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   true,
				"message": err.Error(),
			})
		}

		validate := utils.NewValidator()

		if err := validate.Struct(userDto); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   true,
				"message": err.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"ok": true,
		})
	})
}
