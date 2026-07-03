package httpserver

import (
	"time"

	"github.com/gin-gonic/gin"
)

type Health struct {
	Service string `json:"service"`
	Status  string `json:"status"`
	Time    string `json:"time"`
}

func NewRouter(serviceName, mode string) *gin.Engine {
	if mode != "" {
		gin.SetMode(mode)
	}

	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, Health{
			Service: serviceName,
			Status:  "ok",
			Time:    time.Now().UTC().Format(time.RFC3339),
		})
	})

	return router
}
