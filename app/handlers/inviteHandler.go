package handlers

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/prateeksonii/shutter-api-go/pkg/configs"
	"github.com/prateeksonii/shutter-api-go/pkg/models"
	"github.com/prateeksonii/shutter-api-go/pkg/utils"
	"gorm.io/gorm"
)

func SendInvite(c *fiber.Ctx) error {
	inviteDto := &models.SendInviteDto{}

	if err := c.BodyParser(inviteDto); err != nil {
		c.Status(fiber.StatusBadRequest)
		return err
	}

	validate := utils.NewValidator()

	if err := validate.Struct(inviteDto); err != nil {
		c.Status(fiber.StatusBadRequest)
		return err
	}

	sender := c.Locals("user").(models.User)

	receiver := new(models.User)

	receiverUsername := c.Params("username")

	result := configs.Db.Where("username = ?", receiverUsername).First(&receiver)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		c.Status(fiber.StatusNotFound)
		return result.Error
	}

	existingInvite := &models.Invite{}

	result = configs.Db.Where("sender_id = ? and receiver_id = ?", sender.ID, receiver.ID).First(existingInvite)

	if result.RowsAffected > 0 {
		c.Status(fiber.StatusConflict)
		return errors.New("invite already sent")
	}

	invite := &models.Invite{
		SenderID:   sender.ID,
		ReceiverID: receiver.ID,
		Message:    inviteDto.Message,
	}

	configs.Db.Create(&invite)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"ok":     true,
		"invite": invite,
	})

}
