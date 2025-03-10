package handler

import (
	"log/slog"
	"millionaire-list/internal/logger"
	"millionaire-list/internal/models"
	"millionaire-list/internal/service"
	"net/http"
	"strconv"

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

// 📋 Получить всех миллионеров
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

// 🔎 Получить миллионера по ID
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

// ➕ Создать миллионера
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

// 🔄 Обновить миллионера
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

// ❌ Удалить миллионера
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

// 🔍 Поиск миллионеров по фильтрам
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
