package handler

import (
	"fmt"
	"log/slog"
	"millionaire-list/internal/logger"
	"millionaire-list/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HomeHandler struct {
	service *service.HomeService
	log     *slog.Logger
}

func NewHomeHandler(service *service.HomeService, log *slog.Logger) *HomeHandler {
	return &HomeHandler{service: service, log: log}
}

// GetHomePage retrieves homepage data.
// @Summary Get homepage data
// @Description Fetches homepage-related data and returns it as JSON.
// @Tags home
// @Produce json
// @Success 200 {object} interface{} "Homepage data successfully retrieved"
// @Failure 500 {object} map[string]string "Failed to get homepage data"
// @Router /home [get]
func (h *HomeHandler) GetHomePage(c *gin.Context) {
	h.log.Info("Received request for homepage data")

	baseURL := fmt.Sprintf("%s://%s", c.Request.URL.Scheme, c.Request.Host)

	h.log.Info("Constructed base URL", slog.String("baseURL", baseURL))

	data, err := h.service.GetHomePageData(baseURL)
	if err != nil {
		h.log.Error("Failed to get homepage data", logger.Err(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get homepage data"})
		return
	}

	h.log.Info("Successfully retrieved homepage data")

	c.JSON(http.StatusOK, data)
}
