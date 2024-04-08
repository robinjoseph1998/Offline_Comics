package routes

import (
	"OFFLINECOMICS/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	r.POST("/fetch", handlers.ReciveComicUrl)
}
