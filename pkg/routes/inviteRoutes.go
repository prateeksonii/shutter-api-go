package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/prateeksonii/shutter-api-go/app/handlers"
)

func InviteRoutes(r *gin.RouterGroup) {
	inviteRouter := r.Group("/auth")
	inviteRouter.POST("/:username", handlers.IsAuthenticated, handlers.SendInvite)
}
