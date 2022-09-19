package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
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

	r := gin.Default()

	r.Use(gin.Logger())

	r.Use(gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(error); ok {
			if c.Writer.Status() == http.StatusOK {
				c.Status(http.StatusInternalServerError)
			}
			c.JSON(c.Writer.Status(), gin.H{
				"error":   true,
				"message": err.Error(),
			})
			c.AbortWithStatus(http.StatusInternalServerError)
		}
	}))

	router := r.Group("/api/v1")

	routes.AuthRoutes(router)
	routes.InviteRoutes(router)
	routes.UserRoutes(router)

	port, hasPortEnv := os.LookupEnv("PORT")

	if !hasPortEnv {
		port = "4000"
	}

	err = r.Run(":" + port)
	log.Println(err.Error())
}
