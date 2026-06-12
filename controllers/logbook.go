package controllers

import (
	"net/http"

	"garagefy-api/config"
	"garagefy-api/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// POST /api/vehicles/:id/logbook
func CreateLogbookEntry(c *gin.Context) {
	vehicleID := c.Param("id")

	userIDContext, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao identificar usuário"})
		return
	}
	userID := userIDContext.(uuid.UUID)

	var vehicle models.Vehicle
	if err := config.DB.Where("id = ? AND user_id = ?", vehicleID, userID).First(&vehicle).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Veículo não encontrado ou você não tem permissão"})
		return
	}

	var input struct {
		Title       string `json:"title" binding:"required"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newEntry := models.LogbookEntry{
		VehicleID:   vehicle.ID,
		Title:       input.Title,
		Description: input.Description,
	}

	if err := config.DB.Create(&newEntry).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao criar registro no logbook"})
		return
	}

	c.JSON(http.StatusCreated, newEntry)
}

// GET /api/vehicles/:id/logbook
func GetLogbookEntries(c *gin.Context) {
	vehicleID := c.Param("id")

	userIDContext, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao identificar usuário"})
		return
	}
	userID := userIDContext.(uuid.UUID)

	var vehicle models.Vehicle
	if err := config.DB.Where("id = ? AND user_id = ?", vehicleID, userID).First(&vehicle).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Veículo não encontrado ou você não tem permissão"})
		return
	}

	var entries []models.LogbookEntry
	if err := config.DB.Where("vehicle_id = ?", vehicle.ID).Find(&entries).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar registros"})
		return
	}

	c.JSON(http.StatusOK, entries)
}

// GET /api/vehicles/:id/logbook/:logbookId
func GetLogbookEntryByID(c *gin.Context) {
	vehicleID := c.Param("id")
	logbookID := c.Param("logbookId")

	userIDContext, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao identificar usuário"})
		return
	}
	userID := userIDContext.(uuid.UUID)

	var vehicle models.Vehicle
	if err := config.DB.Where("id = ? AND user_id = ?", vehicleID, userID).First(&vehicle).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Veículo não encontrado ou você não tem permissão"})
		return
	}

	var entry models.LogbookEntry
	if err := config.DB.Where("id = ? AND vehicle_id = ?", logbookID, vehicle.ID).First(&entry).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Registro de logbook não encontrado"})
		return
	}

	c.JSON(http.StatusOK, entry)
}

// PUT /api/vehicles/:id/logbook/:logbookId
func UpdateLogbookEntry(c *gin.Context) {
	vehicleID := c.Param("id")
	logbookID := c.Param("logbookId")

	userIDContext, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao identificar usuário"})
		return
	}
	userID := userIDContext.(uuid.UUID)

	var vehicle models.Vehicle
	if err := config.DB.Where("id = ? AND user_id = ?", vehicleID, userID).First(&vehicle).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Veículo não encontrado ou você não tem permissão"})
		return
	}

	var entry models.LogbookEntry
	if err := config.DB.Where("id = ? AND vehicle_id = ?", logbookID, vehicle.ID).First(&entry).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Registro não encontrado"})
		return
	}

	var input struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.Title != "" {
		entry.Title = input.Title
	}
	if input.Description != "" {
		entry.Description = input.Description
	}

	config.DB.Save(&entry)
	c.JSON(http.StatusOK, entry)
}

// DELETE /api/vehicles/:id/logbook/:logbookId
func DeleteLogbookEntry(c *gin.Context) {
	vehicleID := c.Param("id")
	logbookID := c.Param("logbookId")

	userIDContext, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao identificar usuário"})
		return
	}
	userID := userIDContext.(uuid.UUID)

	var vehicle models.Vehicle
	if err := config.DB.Where("id = ? AND user_id = ?", vehicleID, userID).First(&vehicle).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Veículo não encontrado ou você não tem permissão"})
		return
	}

	var entry models.LogbookEntry
	if err := config.DB.Where("id = ? AND vehicle_id = ?", logbookID, vehicle.ID).First(&entry).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Registro não encontrado"})
		return
	}

	config.DB.Delete(&entry)
	c.JSON(http.StatusOK, gin.H{"message": "Registro removido com sucesso"})
}
