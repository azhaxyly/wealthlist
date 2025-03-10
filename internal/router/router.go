package router

import (
	"millionaire-list/internal/handler"
	"millionaire-list/internal/logger"
	"time"

	"github.com/gin-gonic/gin"
)

func SetupRouter(millionaireHandler *handler.MillionaireHandler, photoHandler *handler.PhotoHandler) *gin.Engine {
	router := gin.Default()

	log := logger.SetupLogger("dev")

	router.Use(func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start)

		log.Info("HTTP-запрос",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", c.Writer.Status(),
			"duration", duration,
		)
	})

	millionaireGroup := router.Group("/api/millionaires")
	{
		millionaireGroup.GET("/", millionaireHandler.GetAll)
		millionaireGroup.GET("/:id", millionaireHandler.GetByID)
		millionaireGroup.POST("/", millionaireHandler.Create)
		millionaireGroup.PUT("/:id", millionaireHandler.Update)
		millionaireGroup.DELETE("/:id", millionaireHandler.Delete)
		millionaireGroup.GET("/search", millionaireHandler.Search)
	}

	// Группа маршрутов для работы с фото
	photoGroup := router.Group("/api/photo")
	{
		photoGroup.POST("/add/:millionaireId", photoHandler.AddPhotoForMillionaire)
		// photoGroup.PUT("/update/:millionaireId", photoHandler.UpdateMillionairePhoto)
		// photoGroup.DELETE("/delete/:millionaireId", photoHandler.DeleteMillionairePhoto)
		// photoGroup.PUT("/update-paths", photoHandler.UpdatePhotoPathForMillionaires)
		// photoGroup.GET("/:imageName", photoHandler.GetPhoto)
	}

	return router
}
