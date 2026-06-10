package controllers

import (
	"garagefy-api/config"
	"garagefy-api/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

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
