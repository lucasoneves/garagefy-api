package controllers

import (
	"net/http"

	"garagefy-api/config"
	"garagefy-api/models"

	"github.com/gin-gonic/gin"
)

// POST /api/vehicles
func CreateVehicle(c *gin.Context) {
	var input struct {
		Brand      string `json:"brand" binding:"required"`
		Model      string `json:"model" binding:"required"`
		Year       int    `json:"year" binding:"required"`
		Plate      string `json:"plate" binding:"required"`
		CurrentOdo int    `json:"current_odo" binding:"required"`
		Color      string `json:"color"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	vehicle := models.Vehicle{
		Brand:      input.Brand,
		Model:      input.Model,
		Year:       input.Year,
		Plate:      input.Plate,
		CurrentOdo: input.CurrentOdo,
		Color:      input.Color,
	}

	if err := config.DB.Create(&vehicle).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao cadastrar veículo. Placa duplicada?"})
		return
	}

	c.JSON(http.StatusCreated, vehicle)
}

// GET /api/vehicles
func GetVehicles(c *gin.Context) {
	var vehicles []models.Vehicle

	if err := config.DB.Order("created_at desc").Find(&vehicles).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar veículos"})
		return
	}

	c.JSON(http.StatusOK, vehicles)
}

// GET /api/vehicles/:id
func GetVehicleByID(c *gin.Context) {
	id := c.Param("id")

	var vehicle models.Vehicle
	// Correção: Uso do Where para compatibilidade garantida com UUID e Preload opcional das linhas de manutenção
	if err := config.DB.Preload("LogbookLines").Where("id = ?", id).First(&vehicle).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Veículo não encontrado"})
		return
	}

	c.JSON(http.StatusOK, vehicle)
}

// PUT /api/vehicles/:id
func UpdateVehicle(c *gin.Context) {
	id := c.Param("id")

	var vehicle models.Vehicle
	// Correção: Uso do Where para UUID
	if err := config.DB.Where("id = ?", id).First(&vehicle).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Veículo não encontrado"})
		return
	}

	var input struct {
		Brand      string `json:"brand"`
		Model      string `json:"model"`
		Year       int    `json:"year"`
		Plate      string `json:"plate"`
		CurrentOdo int    `json:"current_odo"`
		Color      string `json:"color"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Correção: Blindagem contra sobrescrita de dados nulos/zerados
	if input.Brand != "" {
		vehicle.Brand = input.Brand
	}
	if input.Model != "" {
		vehicle.Model = input.Model
	}
	if input.Year != 0 {
		vehicle.Year = input.Year
	}
	if input.Plate != "" {
		vehicle.Plate = input.Plate
	}
	if input.CurrentOdo != 0 {
		vehicle.CurrentOdo = input.CurrentOdo
	}
	if input.Color != "" {
		vehicle.Color = input.Color
	}

	if err := config.DB.Save(&vehicle).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao atualizar veículo. Conflito de placa?"})
		return
	}

	c.JSON(http.StatusOK, vehicle)
}

// DELETE /api/vehicles/:id
func DeleteVehicle(c *gin.Context) {
	id := c.Param("id")

	var vehicle models.Vehicle
	// Correção: Uso do Where para UUID
	if err := config.DB.Where("id = ?", id).First(&vehicle).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Veículo não encontrado"})
		return
	}

	if err := config.DB.Delete(&vehicle).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao deletar veículo"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Veículo removido com sucesso"})
}
