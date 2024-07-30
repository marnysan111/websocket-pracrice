package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/marnysan111/websocket-pracrice/internal/websocket"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// APIエンドポイントに/apiをつける
	api := r.Group("/api")

	api.GET("", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"result":  "success",
			"message": []string{"hoge", "puge", "nuge"},
		})
	})

	r.GET("/ws/:roomID", func(c *gin.Context) {
		roomID := c.Param("roomID")
		websocket.ConnHandler(c, roomID)
	})
	return r
}
