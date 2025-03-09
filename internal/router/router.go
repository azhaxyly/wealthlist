package router

import (
	"millionaire-list/internal/handler"

	"github.com/gin-gonic/gin"
)

func SetupRouter(millionaireHandler *handler.MillionaireHandler) *gin.Engine {
	router := gin.Default()

	millionaireGroup := router.Group("/api/millionaires")
	{
		millionaireGroup.GET("/", millionaireHandler.GetAll)
		millionaireGroup.GET("/:id", millionaireHandler.GetByID)
		millionaireGroup.POST("/", millionaireHandler.Create)
		millionaireGroup.PUT("/:id", millionaireHandler.Update)
		millionaireGroup.DELETE("/:id", millionaireHandler.Delete)
	}

	return router
}
