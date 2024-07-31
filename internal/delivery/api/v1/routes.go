package v1

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/zsandibe/messaggio-microservice/docs"
)

func (h *Handler) Routes() *gin.Engine {
	router := gin.Default()
	router.Use(gin.Recovery())
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	api := router.Group("/api/v1")
	{
		messages := api.Group("/messages")
		{
			messages.GET("", h.getMessagesList)
			messages.POST("", h.addMessage)
			messages.GET("/:id", h.getMessageById)
			messages.DELETE("/:id", h.deleteMessageById)
		}

		stats := api.Group("/stats")
		{
			stats.GET("", h.getStatsList)
			stats.GET("/:id", h.getStatById)
		}
	}
	return router
}
