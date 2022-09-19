package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prateeksonii/shutter-api-go/pkg/configs"
	"github.com/prateeksonii/shutter-api-go/pkg/models"
	"gorm.io/gorm"
)

func SendInvite(c *gin.Context) {
	inviteDto := models.SendInviteDto{}

	if err := c.ShouldBindJSON(&inviteDto); err != nil {
		c.Status(http.StatusBadRequest)
		panic(err)
	}

	sender := c.MustGet("user").(models.User)

	receiver := new(models.User)

	receiverUsername := c.Param("username")

	result := configs.Db.Where("username = ?", receiverUsername).First(&receiver)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		c.Status(http.StatusNotFound)
		panic(result.Error)
	}

	existingInvite := &models.Invite{}

	result = configs.Db.Where("sender_id = ? and receiver_id = ?", sender.ID, receiver.ID).First(existingInvite)

	if result.RowsAffected > 0 {
		c.Status(http.StatusConflict)
		panic(errors.New("invite already sent"))
	}

	invite := &models.Invite{
		SenderID:   sender.ID,
		ReceiverID: receiver.ID,
		Message:    inviteDto.Message,
	}

	configs.Db.Create(&invite)

	c.JSON(http.StatusOK, gin.H{
		"ok":     true,
		"invite": invite,
	})

}
