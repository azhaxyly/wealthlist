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

// AddPhotoForMillionaire uploads a photo for a specific millionaire.
// @Summary Upload a photo for a millionaire
// @Description Allows uploading a photo file for an existing millionaire.
// @Tags millionaires
// @Accept multipart/form-data
// @Produce json
// @Param photo formData file true "Photo file to upload"
// @Param id path int true "Millionaire ID"
// @Success 200 {object} map[string]string "Photo uploaded successfully"
// @Failure 400 {object} map[string]string "Error receiving file"
// @Failure 500 {object} map[string]string "Error uploading or updating photo"
// @Router /api/photo/add/{millionaireId} [post]
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

// DeleteMillionairePhoto deletes a millionaire's photo.
// @Summary Delete a millionaire's photo
// @Description Removes the associated photo of a millionaire and clears its record in the database.
// @Tags millionaires
// @Produce json
// @Param id path int true "Millionaire ID"
// @Success 200 {object} map[string]string "Photo deleted successfully"
// @Failure 404 {object} map[string]string "No photo found for this millionaire"
// @Failure 500 {object} map[string]string "Error deleting or clearing photo path"
// @Router /api/photo/delete/{millionaireId} [delete]
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

// GetPhoto retrieves a millionaire's photo.
// @Summary Get a millionaire's photo
// @Description Serves an image file from the uploads/photos directory based on the provided image name.
// @Tags millionaires
// @Produce image/jpeg
// @Param imageName path string true "Image filename"
// @Success 200 "Returns the requested image file"
// @Failure 400 {object} map[string]string "Image name is required"
// @Failure 404 {object} map[string]string "Image not found"
// @Router /api/photo/{imageName} [get]
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
