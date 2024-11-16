package api

import (
	"blog-backend/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes() {
	router := gin.Default()

	router.MaxMultipartMemory = 8 << 20 // 8 MiB

	router.Use(corsMiddleware())

	router.POST("/users", handlers.HandleAddUsers)
	router.POST("/posts", handlers.HandleAddPosts)
	router.GET("/posts", handlers.HandleFetchPosts)

	router.Run()
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
