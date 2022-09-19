package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/prateeksonii/shutter-api-go/app/handlers"
)

func AuthRoutes(r gin.RouterGroup) {
	r.POST("/signup", handlers.SignUp)
	r.POST("/signin", handlers.SignIn)
	r.GET("/me", handlers.IsAuthenticated, handlers.GetAuthenticatedUser)
}
