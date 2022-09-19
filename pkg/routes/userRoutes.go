package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/prateeksonii/shutter-api-go/app/handlers"
)

func UserRoutes(r *gin.RouterGroup) {
	userRouter := r.Group("/auth")
	userRouter.GET("/", handlers.IsAuthenticated, handlers.SearchUsersByUsername)
}
