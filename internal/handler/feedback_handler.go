package handler

import (
	"log/slog"
	"millionaire-list/internal/logger"
	"millionaire-list/internal/models"
	"millionaire-list/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type FeedbackHandler struct {
	service  *service.FeedbackService
	log      *slog.Logger
	validate *validator.Validate
}

func NewFeedbackHandler(service *service.FeedbackService, log *slog.Logger) *FeedbackHandler {
	return &FeedbackHandler{
		service:  service,
		log:      log,
		validate: validator.New(),
	}
}

func (h *FeedbackHandler) SendFeedback(c *gin.Context) {
	h.log.Info("Received feedback submission request")

	var feedback models.FeedbackDto
	if err := c.ShouldBindJSON(&feedback); err != nil {
		h.log.Error("Failed to parse JSON", logger.Err(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid data format",
		})
		return
	}

	if err := h.validate.Struct(feedback); err != nil {
		h.log.Error("Validation error", logger.Err(err))

		var validationErrors []string
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, err.Field()+" does not meet the requirements")
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation error",
			"details": validationErrors,
		})
		return
	}

	err := h.service.SendFeedbackEmail(feedback)
	if err != nil {
		h.log.Error("Failed to send feedback via email", logger.Err(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error while sending feedback",
		})
		return
	}

	h.log.Info("Feedback successfully sent via email",
		slog.String("name", feedback.Name),
		slog.String("email", feedback.Email))

	c.JSON(http.StatusOK, gin.H{
		"message": "Feedback successfully sent",
	})
}
