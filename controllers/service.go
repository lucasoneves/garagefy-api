package controllers

import (
	"net/http"
	"time"

	"garagefy-api/config"
	"garagefy-api/models"
	"garagefy-api/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// POST /api/services
// POST /api/services
func CreateService(c *gin.Context) {
	userIDContext, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao identificar usuário"})
		return
	}
	userID := userIDContext.(uuid.UUID)

	var input struct {
		VehicleID   string    `json:"vehicle_id" binding:"required"`
		Title       string    `json:"title" binding:"required"`
		Description string    `json:"description"`
		ShopName    string    `json:"shop_name"`
		CurrentOdo  int       `json:"current_odo"`
		Cost        float64   `json:"cost" binding:"required"`
		ServiceDate time.Time `json:"service_date" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	var vehicle models.Vehicle
	if err := config.DB.Where("id = ? AND user_id = ?", input.VehicleID, userID).First(&vehicle).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Veículo não encontrado ou você não tem permissão"})
		return
	}

	newService := models.Service{
		VehicleID:   vehicle.ID,
		Title:       input.Title,
		Description: input.Description,
		ShopName:    input.ShopName,
		CurrentOdo:  input.CurrentOdo,
		Cost:        input.Cost,
		ServiceDate: input.ServiceDate,
	}

	if err := config.DB.Create(&newService).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao registrar serviço"})
		return
	}

	c.JSON(http.StatusCreated, newService)
}

// GET /api/services
func GetServicesByVehicle(c *gin.Context) {
	vehicleID := c.Query("vehicle_id")
	if vehicleID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "O parâmetro vehicle_id é obrigatório"})
		return
	}

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

	var services []models.Service
	if err := config.DB.Where("vehicle_id = ?", vehicle.ID).Find(&services).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar serviços"})
		return
	}

	c.JSON(http.StatusOK, services)
}

// GET /api/services/:id
func GetServiceByID(c *gin.Context) {
	serviceID := c.Param("id")

	userIDContext, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao identificar usuário"})
		return
	}
	userID := userIDContext.(uuid.UUID)

	var service models.Service
	// Segurança por Joins: Busca o serviço garantindo que o veículo associado pertence ao usuário logado
	if err := config.DB.Joins("JOIN vehicles ON vehicles.id = services.vehicle_id").
		Where("services.id = ? AND vehicles.user_id = ?", serviceID, userID).
		First(&service).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Serviço não encontrado ou você não tem permissão"})
		return
	}

	c.JSON(http.StatusOK, service)
}

// PUT /api/services/:id
func UpdateService(c *gin.Context) {
	serviceID := c.Param("id")

	userIDContext, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao identificar usuário"})
		return
	}
	userID := userIDContext.(uuid.UUID)

	var service models.Service
	if err := config.DB.Joins("JOIN vehicles ON vehicles.id = services.vehicle_id").
		Where("services.id = ? AND vehicles.user_id = ?", serviceID, userID).
		First(&service).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Serviço não encontrado ou você não tem permissão"})
		return
	}

	var input struct {
		Title       string    `json:"title"`
		Description string    `json:"description"`
		ShopName    string    `json:"shop_name"`
		CurrentOdo  int       `json:"current_odo"`
		Cost        float64   `json:"cost"`
		ServiceDate time.Time `json:"service_date"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	if input.Title != "" {
		service.Title = input.Title
	}
	if input.Description != "" {
		service.Description = input.Description
	}
	if input.ShopName != "" {
		service.ShopName = input.ShopName
	}
	if input.CurrentOdo != 0 {
		service.CurrentOdo = input.CurrentOdo
	}
	if input.Cost != 0 {
		service.Cost = input.Cost
	}
	if !input.ServiceDate.IsZero() {
		service.ServiceDate = input.ServiceDate
	}

	config.DB.Save(&service)
	c.JSON(http.StatusOK, service)
}

// DELETE /api/services/:id
func DeleteService(c *gin.Context) {
	serviceID := c.Param("id")

	userIDContext, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao identificar usuário"})
		return
	}
	userID := userIDContext.(uuid.UUID)

	var service models.Service
	if err := config.DB.Joins("JOIN vehicles ON vehicles.id = services.vehicle_id").
		Where("services.id = ? AND vehicles.user_id = ?", serviceID, userID).
		First(&service).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Serviço não encontrado ou você não tem permissão"})
		return
	}

	config.DB.Delete(&service)
	c.JSON(http.StatusOK, gin.H{"message": "Serviço removido com sucesso"})
}
