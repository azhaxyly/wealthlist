package handler

import (
	"log/slog"
	"net/http"
	"strconv"

	"millionaire-list/internal/logger"
	"millionaire-list/internal/service"

	"github.com/gin-gonic/gin"
)

type PhotoHandler struct {
	photoService *service.PhotoService
	log          *slog.Logger
}

func NewPhotoHandler(photoService *service.PhotoService, log *slog.Logger) *PhotoHandler {
	return &PhotoHandler{
		photoService: photoService,
		log:          log,
	}
}

// Загрузка фото для миллионера
func (h *PhotoHandler) AddPhotoForMillionaire(c *gin.Context) {
	// Получаем ID миллионера из URL
	millionaireID, err := strconv.Atoi(c.Param("millionaireId"))
	if err != nil {
		h.log.Error("Incorrect ID", logger.Err(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect millionaire ID"})
		return
	}

	// Получаем файл из запроса
	file, err := c.FormFile("photo")
	if err != nil {
		h.log.Error("Error receiving file", logger.Err(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error receiving file"})
		return
	}

	// Загружаем файл и получаем путь
	filePath, err := h.photoService.UploadPhoto(millionaireID, file)
	if err != nil {
		h.log.Error("Error uploading file", logger.Err(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Обновляем путь к фото в БД
	err = h.photoService.UpdatePhoto(millionaireID, filePath)
	if err != nil {
		h.log.Error("Error updating millionaire photo path", logger.Err(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update millionaire photo"})
		return
	}

	// Отправляем успешный ответ
	c.JSON(http.StatusOK, gin.H{
		"message":   "Photo uploaded successfully",
		"photoPath": filePath,
	})
}
