package router

import (
	_ "millionaire-list/docs"
	"millionaire-list/internal/handler"
	"millionaire-list/internal/logger"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter(millionaireHandler *handler.MillionaireHandler, photoHandler *handler.PhotoHandler, homeHandler *handler.HomeHandler, feedbackHandler *handler.FeedbackHandler) *gin.Engine {
	router := gin.Default()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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

	photoGroup := router.Group("/api/photo")
	{
		photoGroup.POST("/add/:millionaireId", photoHandler.AddPhotoForMillionaire)
		photoGroup.DELETE("/delete/:millionaireId", photoHandler.DeleteMillionairePhoto)
		photoGroup.GET("/:imageName", photoHandler.GetPhoto)
	}

	homeGroup := router.Group("/home")
	{
		homeGroup.GET("/", homeHandler.GetHomePage)
	}

	feedbackGroup := router.Group("/api/feedback")
	{
		feedbackGroup.POST("/", feedbackHandler.SendFeedback)
	}

	return router
}
