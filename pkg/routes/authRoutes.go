package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/prateeksonii/shutter-api-go/app/handlers"
)

func AuthRoutes(r *gin.RouterGroup) {
	authRouter := r.Group("/auth")
	authRouter.POST("/signup", handlers.SignUp)
	authRouter.POST("/signin", handlers.SignIn)
	authRouter.GET("/me", handlers.IsAuthenticated, handlers.GetAuthenticatedUser)
}
