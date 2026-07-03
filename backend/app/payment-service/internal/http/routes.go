package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"service": "payment-service"})
	})

	group := router.Group("/payments")
	group.GET("", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": []gin.H{}})
	})
	group.POST("", func(c *gin.Context) {
		c.JSON(http.StatusAccepted, gin.H{"message": "payment create endpoint"})
	})
}
