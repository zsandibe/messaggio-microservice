package v1

import "github.com/gin-gonic/gin"

func (h *Handler) Routes() *gin.Engine {
	router := gin.Default()
	router.Use(gin.Recovery())
	// router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	api := router.Group("/api/v1")
	{
		messages := api.Group("/messages")
		{
			messages.GET("/", h.getMessagesList)
			messages.POST("/", h.addMessage)
			messages.GET("/:id", h.getMessageById)
			messages.PUT("/:id")
			messages.DELETE("/:id", h.deleteMessageById)
		}

		stats := api.Group("/stats")
		{
			stats.POST("/")
			stats.PUT("/:id")
			stats.DELETE("/:id")
			stats.GET("/user/:id")
		}
	}
	return router
}
