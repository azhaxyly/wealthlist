package router

import (
	"millionaire-list/internal/handler"

	"github.com/gin-gonic/gin"
)

// SetupRouter инициализирует маршруты Gin
func SetupRouter(
	millionaireHandler *handler.MillionaireHandler,
	// feedbackHandler *handler.FeedbackHandler,
	// photoHandler *handler.PhotoHandler,
) *gin.Engine {
	router := gin.Default()

	// Эндпоинты для миллионеров
	millionaireGroup := router.Group("/api/millionaires")
	{
		millionaireGroup.GET("/", millionaireHandler.GetAll)       // Получить всех миллионеров
		millionaireGroup.GET("/:id", millionaireHandler.GetByID)   // Получить миллионера по ID
		millionaireGroup.POST("/", millionaireHandler.Create)      // Создать миллионера
		millionaireGroup.PUT("/:id", millionaireHandler.Update)    // Обновить миллионера
		millionaireGroup.DELETE("/:id", millionaireHandler.Delete) // Удалить миллионера
	}

	// // Эндпоинты для отзывов
	// feedbackGroup := router.Group("/api/feedback")
	// {
	// 	feedbackGroup.POST("/", feedbackHandler.SendFeedback) // Оставить отзыв
	// }

	// // Эндпоинты для фотографий
	// photoGroup := router.Group("/api/photos")
	// {
	// 	photoGroup.POST("/upload", photoHandler.Upload)      // Загрузить фото
	// 	photoGroup.GET("/:filename", photoHandler.Get)       // Получить фото
	// 	photoGroup.DELETE("/:filename", photoHandler.Delete) // Удалить фото
	// }

	return router
}
