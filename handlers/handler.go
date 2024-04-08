package handlers

import (
	"OFFLINECOMICS/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ReciveComicUrl(c *gin.Context) {
	var URL string
	if err := c.ShouldBind(&URL); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad Request", "error": err.Error()})
		return
	}
	n := 200
	usecase.Fetch(URL, n)
}
