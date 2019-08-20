package auth

import (
	"github.com/gin-gonic/gin"
)

// ApplyRoutes auth
func ApplyRoutes(r *gin.RouterGroup) {
	auth := r.Group("/auth")
	{
		auth.POST("", login)
		auth.POST("/create", create)
		auth.POST("/create-password", passwordCreation)
		auth.POST("/reset-password", passwordReset)
		auth.GET("", list)
		auth.GET("/:id", read)
		auth.DELETE("/:id", remove)
		auth.PATCH("/:id", update)
	}
}
