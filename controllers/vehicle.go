package controllers

import (
	"net/http"

	"garagefy-api/config"
	"garagefy-api/models"
	"garagefy-api/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// POST /api/vehicles
func CreateVehicle(c *gin.Context) {
	// 1. Captura o userID injetado pelo AuthMiddleware
	userIDContext, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao identificar usuário"})
		return
	}
	userID := userIDContext.(uuid.UUID) // Converte a interface para o tipo uuid.UUID

	var input struct {
		Brand      string `json:"brand" binding:"required"`
		Model      string `json:"model" binding:"required"`
		Year       int    `json:"year" binding:"required"`
		Plate      string `json:"plate"`
		CurrentOdo int    `json:"current_odo"`
		Color      string `json:"color"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	// 2. Instancia o veículo já injetando o UserID do dono
	newVehicle := models.Vehicle{
		UserID:     userID,
		Brand:      input.Brand,
		Model:      input.Model,
		Year:       input.Year,
		Plate:      input.Plate,
		CurrentOdo: input.CurrentOdo,
		Color:      input.Color,
	}

	if err := config.DB.Create(&newVehicle).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao cadastrar veículo"})
		return
	}

	c.JSON(http.StatusCreated, newVehicle)
}

// GET /api/vehicles
func GetVehicles(c *gin.Context) {
	userIDContext, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao identificar usuário"})
		return
	}
	userID := userIDContext.(uuid.UUID)

	var vehicles []models.Vehicle
	// <-- Filtro Crítico: Traz apenas onde user_id bate com o ID do token
	if err := config.DB.Where("user_id = ?", userID).Find(&vehicles).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar veículos"})
		return
	}

	c.JSON(http.StatusOK, vehicles)
}

// GET /api/vehicles/:id
func GetVehicleByID(c *gin.Context) {
	vehicleID := c.Param("id")
	userIDContext, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao identificar usuário"})
		return
	}
	userID := userIDContext.(uuid.UUID)

	var vehicle models.Vehicle
	// Busca validando os dois escopos ao mesmo tempo
	if err := config.DB.Where("id = ? AND user_id = ?", vehicleID, userID).First(&vehicle).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Veículo não encontrado ou você não tem permissão"})
		return
	}

	c.JSON(http.StatusOK, vehicle)
}

// PUT /api/vehicles/:id
func UpdateVehicle(c *gin.Context) {
	id := c.Param("id")

	// 1. Captura o userID injetado pelo AuthMiddleware
	userIDContext, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao identificar usuário"})
		return
	}
	userID := userIDContext.(uuid.UUID)

	var vehicle models.Vehicle
	// Trava de segurança: Garante que o veículo pertence ao usuário que está tentando atualizar
	if err := config.DB.Where("id = ? AND user_id = ?", id, userID).First(&vehicle).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Veículo não encontrado ou você não tem permissão"})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	// Blindagem contra sobrescrita de dados nulos/zerados
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

	// 1. Captura o userID injetado pelo AuthMiddleware
	userIDContext, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao identificar usuário"})
		return
	}
	userID := userIDContext.(uuid.UUID)

	var vehicle models.Vehicle
	// Trava de segurança: Garante que só o dono consegue deletar o veículo
	if err := config.DB.Where("id = ? AND user_id = ?", id, userID).First(&vehicle).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Veículo não encontrado ou você não tem permissão"})
		return
	}

	if err := config.DB.Delete(&vehicle).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao deletar veículo"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Veículo removido com sucesso"})
}
