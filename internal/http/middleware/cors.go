package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// CORSMiddleware is a layer security middleware that provides CORS security.
// This middleware defines allow headers, methods for CORS pre-flight OPTIONS requests.
//
// If the request method is OPTIONS, server returns http.StatusNoContent status code.
// In other case, continue the process.
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-Access-Token, X-Refresh-Token")
		c.Header("Access-Control-Allow-Methods", "POST,HEAD,PATCH, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
