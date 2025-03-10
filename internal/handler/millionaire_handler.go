package handler

import (
	"fmt"
	"log"
	"millionaire-list/internal/models"
	"millionaire-list/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type MillionaireHandler struct {
	service *service.MillionaireService
}

func NewMillionaireHandler(service *service.MillionaireService) *MillionaireHandler {
	return &MillionaireHandler{service: service}
}

// 📋 Получить всех миллионеров
func (mh *MillionaireHandler) GetAll(c *gin.Context) {
	pageNum, _ := strconv.Atoi(c.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	result, err := mh.service.GetAllMillionaires(pageNum, pageSize)
	if err != nil {
		fmt.Println("Ошибка в GetAllMillionaires:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении данных"})
		return
	}

	c.JSON(http.StatusOK, result)
}

// 🔎 Получить миллионера по ID
func (mh *MillionaireHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный ID"})
		return
	}

	millionaire, err := mh.service.GetMillionaireByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Миллионер не найден"})
		return
	}

	c.JSON(http.StatusOK, millionaire)
}

// ➕ Создать миллионера
func (mh *MillionaireHandler) Create(c *gin.Context) {
	var millionaire models.Millionaire
	if err := c.ShouldBindJSON(&millionaire); err != nil {
		log.Println("Ошибка разбора JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный JSON"})
		return
	}

	err := mh.service.CreateMillionaire(&millionaire)
	if err != nil {
		log.Println("Ошибка создания миллионера:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при создании"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Миллионер создан"})
}

// 🔄 Обновить миллионера
func (mh *MillionaireHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный ID"})
		return
	}

	var millionaire models.Millionaire
	if err := c.ShouldBindJSON(&millionaire); err != nil {
		log.Println("Ошибка разбора JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный JSON"})
		return
	}

	millionaire.ID = id
	err = mh.service.UpdateMillionaire(&millionaire)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при обновлении"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Миллионер обновлен"})
}

// ❌ Удалить миллионера
func (mh *MillionaireHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный ID"})
		return
	}

	err = mh.service.DeleteMillionaire(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при удалении"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Миллионер удален"})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при поиске миллионеров"})
		return
	}

	c.JSON(http.StatusOK, result)
}
