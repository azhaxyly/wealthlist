package handler

import (
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
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

func (h *PhotoHandler) getMillionaireID(c *gin.Context) (int, bool) {
	millionaireID, err := strconv.Atoi(c.Param("millionaireId"))
	if err != nil {
		h.log.Error("Incorrect ID", logger.Err(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect millionaire ID"})
		return 0, false
	}
	return millionaireID, true
}

func (h *PhotoHandler) AddPhotoForMillionaire(c *gin.Context) {
	millionaireID, ok := h.getMillionaireID(c)
	if !ok {
		return
	}

	file, err := c.FormFile("photo")
	if err != nil {
		h.log.Error("Error receiving file", logger.Err(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error receiving file"})
		return
	}

	filePath, err := h.photoService.UploadPhoto(millionaireID, file)
	if err != nil {
		h.log.Error("Error uploading file", logger.Err(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = h.photoService.UpdatePhoto(millionaireID, filePath)
	if err != nil {
		h.log.Error("Error updating millionaire photo path", logger.Err(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update millionaire photo"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Photo uploaded successfully",
		"photoPath": filePath,
	})
}

func (h *PhotoHandler) DeleteMillionairePhoto(c *gin.Context) {
	millionaireID, ok := h.getMillionaireID(c)
	if !ok {
		return
	}

	photoPath, err := h.photoService.GetPhotoPath(millionaireID)
	if err != nil {
		h.log.Error("Error retrieving photo path", logger.Err(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve photo path"})
		return
	}

	if photoPath == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "No photo found for this millionaire"})
		return
	}

	if err := os.Remove(photoPath); err != nil && !os.IsNotExist(err) {
		h.log.Error("Error deleting photo file", logger.Err(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete photo file"})
		return
	}

	err = h.photoService.ClearPhotoPath(millionaireID)
	if err != nil {
		h.log.Error("Error clearing photo path", logger.Err(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear photo path"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Photo deleted successfully"})
}

func (h *PhotoHandler) GetPhoto(c *gin.Context) {
	imageName := c.Param("imageName")
	if imageName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Image name is required"})
		return
	}

	imagePath := filepath.Join("uploads/photos", imageName)

	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Image not found"})
		return
	}

	c.File(imagePath)
}
