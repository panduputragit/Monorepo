package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"service": "employee-service"})
	})

	group := router.Group("/employees")
	group.GET("", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": []gin.H{}})
	})
	group.POST("", func(c *gin.Context) {
		c.JSON(http.StatusAccepted, gin.H{"message": "employee create endpoint"})
	})
}
