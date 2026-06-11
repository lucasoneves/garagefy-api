package controllers

import (
	"garagefy-api/config"
	"garagefy-api/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// GET /api/services
func GetServicesByVehicle(c *gin.Context) {
	vehicleID := c.Query("vehicle_id")

	if vehicleID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "O parâmetro vehicle_id é obrigatório"})
		return
	}

	var services []models.Service
	// Busca os serviços do carro ordenados pelos mais recentes
	if err := config.DB.Where("vehicle_id = ?", vehicleID).Order("service_date DESC").Find(&services).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao buscar serviços"})
		return
	}

	c.JSON(http.StatusOK, services)
}

// GET /api/services/:id
func GetServiceByID(c *gin.Context) {
	id := c.Param("id")
	var service models.Service

	if err := config.DB.First(&service, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Serviço não encontrado"})
		return
	}

	c.JSON(http.StatusOK, service)
}

// POST /api/services
func CreateService(c *gin.Context) {
	var input struct {
		VehicleID   string  `json:"vehicle_id" binding:"required"`
		Title       string  `json:"title" binding:"required"`
		Description string  `json:"description"`
		ShopName    string  `json:"shop_name"`
		CurrentOdo  int     `json:"current_odo" binding:""`
		Cost        float64 `json:"cost" binding:"required"`
		ServiceDate string  `json:"service_date" binding:"required"`
	}

	// Valida o JSON estrito vindo do Front-end
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Faz o parse manual da string de data ISO que o Next enviou
	parsedDate, err := time.Parse(time.RFC3339, input.ServiceDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Formato de data inválido. Use o padrão ISO."})
		return
	}

	service := models.Service{
		VehicleID:   input.VehicleID,
		Title:       input.Title,
		Description: input.Description,
		ShopName:    input.ShopName,
		CurrentOdo:  input.CurrentOdo,
		Cost:        input.Cost,
		ServiceDate: parsedDate,
	}

	// Persiste o registro na nova tabela independente
	if err := config.DB.Create(&service).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao salvar o serviço no banco de dados"})
		return
	}

	c.JSON(http.StatusCreated, service)
}

// PUT /api/services/:id
func UpdateService(c *gin.Context) {
	id := c.Param("id")
	var service models.Service

	// Verifica se o serviço realmente existe antes de atualizar
	if err := config.DB.First(&service, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Serviço não encontrado"})
		return
	}

	// Struct interna para validação do payload de atualização
	var input struct {
		Title       string  `json:"title" binding:"required"`
		Description string  `json:"description"`
		ShopName    string  `json:"shop_name"`
		CurrentOdo  int     `json:"current_odo" binding:"required"`
		Cost        float64 `json:"cost" binding:"required"`
		ServiceDate string  `json:"service_date" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse da data ISO enviada pelo Next
	parsedDate, err := time.Parse(time.RFC3339, input.ServiceDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Formato de data inválido. Use o padrão ISO."})
		return
	}

	// Atualiza os campos do modelo mapeado
	service.Title = input.Title
	service.Description = input.Description
	service.ShopName = input.ShopName
	service.CurrentOdo = input.CurrentOdo
	service.Cost = input.Cost
	service.ServiceDate = parsedDate

	if err := config.DB.Save(&service).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao atualizar o serviço no banco"})
		return
	}

	c.JSON(http.StatusOK, service)
}

// DELETE /api/services/:id
func DeleteService(c *gin.Context) {
	id := c.Param("id")
	var service models.Service

	if err := config.DB.First(&service, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Serviço não encontrado"})
		return
	}

	if err := config.DB.Delete(&service).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao eliminar o serviço do banco de dados"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Serviço eliminado com sucesso"})
}
