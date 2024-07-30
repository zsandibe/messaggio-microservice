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
			messages.GET("/")
			messages.POST("/")
			messages.GET("/:id")
			messages.PUT("/:id")
			messages.DELETE("/:id")
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
