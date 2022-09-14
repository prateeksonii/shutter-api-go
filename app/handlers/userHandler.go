package handlers

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/prateeksonii/shutter-api-go/pkg/configs"
	"github.com/prateeksonii/shutter-api-go/pkg/models"
)

func SearchUsersByUsername(c *fiber.Ctx) error {
	username := c.Query("username")

	if len(username) == 0 {
		c.Status(fiber.StatusBadRequest)
		return errors.New("no username provided")
	}

	users := &[]models.User{}

	result := configs.Db.Where("username like ?", "%"+username+"%").Find(&users)

	if result.RowsAffected == 0 {
		c.Status(fiber.StatusNotFound)
		return errors.New("no records found")
	}

	return c.JSON(fiber.Map{
		"ok":    true,
		"users": users,
	})
}
