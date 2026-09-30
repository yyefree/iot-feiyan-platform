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

	// Register each backend with a catch-all handler
	for prefix, backend := range backends {
		proxy := httputil.NewSingleHostReverseProxy(backend)
		originalDirector := proxy.Director
		proxy.Director = func(req *http.Request) {
			originalDirector(req)
			// Strip the prefix from the path
			newPath := strings.TrimPrefix(req.URL.Path, prefix)
			if newPath == "" {
				newPath = "/"
			}
			req.URL.Path = newPath
			req.Host = backend.Host
			req.Header.Set("X-Forwarded-Prefix", prefix)
			req.Header.Set("X-Real-IP", req.RemoteAddr)
		}
		proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
			log.Printf("proxy error for %s: %v", prefix, err)
			http.Error(w, "Bad Gateway", http.StatusBadGateway)
		}
		// Register all HTTP methods with catch-all
		r.Any(prefix+"/*path", gin.WrapH(proxy))
	}

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "api-gateway"})
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
