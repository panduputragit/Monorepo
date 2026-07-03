package http

import (
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, serviceURLs map[string]string) error {
	api := router.Group("/api/v1")
	for service, rawURL := range serviceURLs {
		target, err := url.Parse(rawURL)
		if err != nil {
			return err
		}

		proxy := httputil.NewSingleHostReverseProxy(target)
		group := api.Group("/" + service)
		group.Any("", proxyHandler(service, target, proxy))
		group.Any("/*path", proxyHandler(service, target, proxy))
	}

	router.POST("/login", forward("/auth/login", serviceURLs["auth"]))
	router.GET("/validate", forward("/auth/validate", serviceURLs["auth"]))
	router.POST("/refresh", forward("/auth/refresh", serviceURLs["auth"]))

	return nil
}

func forward(path, rawURL string) gin.HandlerFunc {
	target, err := url.Parse(rawURL)
	if err != nil {
		return func(c *gin.Context) {
			c.JSON(500, gin.H{"error": "invalid service URL"})
		}
	}

	proxy := httputil.NewSingleHostReverseProxy(target)
	return func(c *gin.Context) {
		c.Request.URL.Path = path
		c.Request.Host = target.Host
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}

func proxyHandler(service string, target *url.URL, proxy *httputil.ReverseProxy) gin.HandlerFunc {
	prefix := "/api/v1/" + service
	return func(c *gin.Context) {
		c.Request.URL.Path = strings.TrimPrefix(c.Request.URL.Path, prefix)
		if c.Request.URL.Path == "" {
			c.Request.URL.Path = "/"
		}
		c.Request.Host = target.Host
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}
