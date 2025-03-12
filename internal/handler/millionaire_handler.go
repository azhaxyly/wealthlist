package handler

import (
	"log/slog"
	"net/http"
	"strconv"
	"wealthlist/internal/logger"
	"wealthlist/internal/models"
	"wealthlist/internal/service"

	"github.com/gin-gonic/gin"
)

type MillionaireHandler struct {
	service service.MillionaireServiceInterface
	log     *slog.Logger
}

func NewMillionaireHandler(service service.MillionaireServiceInterface, log *slog.Logger) *MillionaireHandler {
	return &MillionaireHandler{
		service: service,
		log:     log,
	}
}

// GetAll retrieves a paginated list of millionaires.
// @Summary Get all millionaires
// @Description Fetches a paginated list of millionaires from the database.
// @Tags millionaires
// @Produce json
// @Param pageNum query int false "Page number (default: 1)"
// @Param pageSize query int false "Page size (default: 10)"
// @Success 200 {object} models.PaginationMillionaireDto "List of millionaires retrieved successfully"
// @Failure 500 {object} map[string]string "Error retrieving data"
// @Router /api/millionaires [get]
func (mh *MillionaireHandler) GetAll(c *gin.Context) {
	pageNum, _ := strconv.Atoi(c.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	result, err := mh.service.GetAllMillionaires(pageNum, pageSize)
	if err != nil {
		mh.log.Error("Error recieving data", logger.Err(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error recieving data"})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetByID retrieves a millionaire by ID.
// @Summary Get millionaire by ID
// @Description Fetches a millionaire's details using their unique ID.
// @Tags millionaires
// @Produce json
// @Param id path int true "Millionaire ID"
// @Success 200 {object} models.Millionaire "Millionaire retrieved successfully"
// @Failure 400 {object} map[string]string "Incorrect ID format"
// @Failure 404 {object} map[string]string "Millionaire not found"
// @Router /api/millionaires/{id} [get]
func (mh *MillionaireHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		mh.log.Error("Incorrect ID", logger.Err(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect ID"})
		return
	}

	millionaire, err := mh.service.GetMillionaireByID(id)
	if err != nil {
		mh.log.Error("Millionaire not found", logger.Err(err))
		c.JSON(http.StatusNotFound, gin.H{"error": "Millionaire not found"})
		return
	}

	c.JSON(http.StatusOK, millionaire)
}

// Create adds a new millionaire.
// @Summary Create a new millionaire
// @Description Adds a new millionaire to the database.
// @Tags millionaires
// @Accept json
// @Produce json
// @Param millionaire body models.Millionaire true "Millionaire data"
// @Success 201 {object} map[string]string "Millionaire created"
// @Failure 400 {object} map[string]string "Incorrect JSON format"
// @Failure 500 {object} map[string]string "Error creating millionaire"
// @Router /api/millionaires [post]
func (mh *MillionaireHandler) Create(c *gin.Context) {
	var millionaire models.Millionaire
	if err := c.ShouldBindJSON(&millionaire); err != nil {
		mh.log.Error("Incorrect JSON", logger.Err(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect JSON"})
		return
	}

	err := mh.service.CreateMillionaire(&millionaire)
	if err != nil {
		mh.log.Error("Error creating millionaire", logger.Err(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating millionaire"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Millionaire created"})
}

// Update modifies an existing millionaire.
// @Summary Update a millionaire
// @Description Updates millionaire details based on the provided ID.
// @Tags millionaires
// @Accept json
// @Produce json
// @Param id path int true "Millionaire ID"
// @Param millionaire body models.Millionaire true "Updated millionaire data"
// @Success 200 {object} map[string]string "Millionaire updated"
// @Failure 400 {object} map[string]string "Incorrect ID or JSON format"
// @Failure 500 {object} map[string]string "Error updating millionaire"
// @Router /api/millionaires/{id} [put]
func (mh *MillionaireHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		mh.log.Error("Incorrect ID", logger.Err(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect ID"})
		return
	}

	var millionaire models.Millionaire
	if err := c.ShouldBindJSON(&millionaire); err != nil {
		mh.log.Error("Incorrect JSON", logger.Err(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect JSON"})
		return
	}

	millionaire.ID = id
	err = mh.service.UpdateMillionaire(&millionaire)
	if err != nil {
		mh.log.Error("Error updating millionaire", logger.Err(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating millionaire"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Millionaire updated"})
}

// Delete removes a millionaire from the system.
// @Summary Delete a millionaire
// @Description Deletes a millionaire based on the provided ID.
// @Tags millionaires
// @Produce json
// @Param id path int true "Millionaire ID"
// @Success 200 {object} map[string]string "Millionaire deleted"
// @Failure 400 {object} map[string]string "Incorrect ID format"
// @Failure 500 {object} map[string]string "Error deleting millionaire"
// @Router /millionaires/{id} [delete]
func (mh *MillionaireHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		mh.log.Error("Incorrect ID", logger.Err(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect ID"})
		return
	}

	err = mh.service.DeleteMillionaire(id)
	if err != nil {
		mh.log.Error("Error deleting millionaire", logger.Err(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting millionaire"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Millionaire deleted"})
}

// Search finds millionaires based on given query parameters.
// @Summary Search for millionaires
// @Description Searches for millionaires using optional filters such as name and country.
// @Tags millionaires
// @Produce json
// @Param lastName query string false "Last name of the millionaire"
// @Param firstName query string false "First name of the millionaire"
// @Param middleName query string false "Middle name of the millionaire"
// @Param country query string false "Country of the millionaire"
// @Param page query int false "Page number" default(1)
// @Param pageSize query int false "Number of records per page" default(10)
// @Success 200 {array} models.Millionaire "List of matching millionaires"
// @Failure 500 {object} map[string]string "Error searching millionaire"
// @Router /millionaires/search [get]
func (mh *MillionaireHandler) Search(c *gin.Context) {
	lastName := c.Query("lastName")
	firstName := c.Query("firstName")
	middleName := c.Query("middleName")
	country := c.Query("country")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	result, err := mh.service.SearchMillionaire(lastName, firstName, middleName, country, page, pageSize)
	if err != nil {
		mh.log.Error("Error searching millionaire", logger.Err(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error searching millionaire"})
		return
	}

	c.JSON(http.StatusOK, result)
}
