package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// CORSMiddleware returns a middleware for handling CORS
func CORSMiddleware() gin.HandlerFunc {
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "X-API-Key", "Authorization"}
	config.AllowCredentials = true
	config.ExposeHeaders = []string{"Content-Length", "X-API-Key"}
	
	return cors.New(config)
}
