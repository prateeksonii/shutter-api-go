package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prateeksonii/shutter-api-go/pkg/configs"
	"github.com/prateeksonii/shutter-api-go/pkg/models"
)

func SearchUsersByUsername(c *gin.Context) {
	username := c.Query("username")

	if len(username) == 0 {
		c.Status(http.StatusBadRequest)
		panic(errors.New("no username provided"))
	}

	users := &[]models.User{}

	result := configs.Db.Where("username like ?", "%"+username+"%").Find(&users)

	if result.RowsAffected == 0 {
		c.Status(http.StatusNotFound)
		panic(errors.New("no records found"))
	}

	c.JSON(http.StatusOK, gin.H{
		"ok":    true,
		"users": users,
	})
}
