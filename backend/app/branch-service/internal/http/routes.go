package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"service": "branch-service"})
	})

	group := router.Group("/branches")
	group.GET("", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": []gin.H{}})
	})
	group.POST("", func(c *gin.Context) {
		c.JSON(http.StatusAccepted, gin.H{"message": "branch create endpoint"})
	})
}
