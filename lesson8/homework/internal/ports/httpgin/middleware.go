package httpgin

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		latency := time.Since(start)
		status := c.Writer.Status()
		log.Println("Latency", latency, "Request method:", c.Request.Method, "URL path:", c.Request.URL.Path, "Response status:", status, http.StatusText(status))
	}
}
