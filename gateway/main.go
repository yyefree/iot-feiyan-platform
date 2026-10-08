package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {
	backends := map[string]*url.URL{
		"/api/v1/devices": mustParse("http://iot-device-service:8081"),
		"/api/v1/products": mustParse("http://iot-product-service:8082"),
		"/api/v1/rules":    mustParse("http://iot-rule-service:8083"),
		"/api/v1/data":     mustParse("http://iot-data-service:8084"),
		"/api/v1/ops":      mustParse("http://iot-ops-service:8085"),
		"/api/v1/push":     mustParse("http://iot-ops-service:8085"),
	}

	r := gin.Default()

	// Register health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"code": 0, "message": "ok", "service": "api-gateway"})
	})

	// Register catch-all handler for all routes
	r.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path
		
		// Find matching backend
		var backend *url.URL
		for prefix, url := range backends {
			if strings.HasPrefix(path, prefix) {
				backend = url
				break
			}
		}
		
		if backend == nil {
			c.JSON(404, gin.H{"code": 404, "message": "Not Found"})
			return
		}
		
		// Strip prefix and proxy
		newPath := strings.TrimPrefix(path, backend.String())
		proxy := httputil.NewSingleHostReverseProxy(backend)
		originalDirector := proxy.Director
		proxy.Director = func(req *http.Request) {
			originalDirector(req)
			req.URL.Path = newPath
			req.Host = backend.Host
		}
		proxy.ServeHTTP(c.Writer, c.Request)
	})

	log.Printf("api-gateway starting on :8080")
	r.Run(":8080")
}

func mustParse(s string) *url.URL {
	u, err := url.Parse(s)
	if err != nil {
		log.Fatalf("invalid backend URL %s: %v", s, err)
	}
	return u
}
