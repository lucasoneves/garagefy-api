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

// POST /api/fuels
func CreateFuelLog(c *gin.Context) {
	userIDContext, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao identificar usuário"})
		return
	}
	userID := userIDContext.(uuid.UUID)

	var input struct {
		VehicleID    string    `json:"vehicle_id" binding:"required"`
		Date         time.Time `json:"date" binding:"required"`
		Odometer     int       `json:"odometer"`
		Liters       float64   `json:"liters" binding:"required"`
		PricePerLite float64   `json:"price_per_liter" binding:"required"`
		IsFullTank   bool      `json:"is_full_tank"`
		GasStation   string    `json:"gas_station"`
		FuelType     string    `json:"fuel_type"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	// 1. Validar se o veículo pertence ao utilizador logado
	var vehicle models.Vehicle
	if err := config.DB.Where("id = ? AND user_id = ?", input.VehicleID, userID).First(&vehicle).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Veículo não encontrado ou você não tem permissão"})
		return
	}

	// 2. Criar a instância do novo registo
	newFuelLog := models.FuelLog{
		VehicleID:   vehicle.ID,
		Date:        input.Date,
		Odometer:    input.Odometer,
		Liters:      input.Liters,
		PricePerLit: input.PricePerLite,
		IsFullTank:  input.IsFullTank,
		GasStation:  input.GasStation,
		FuelType:    input.FuelType,
	}

	// 3. Lógica do Cálculo de KM/L (Consumo)
	// Só faz sentido calcular se o tanque atual estiver cheio
	if newFuelLog.IsFullTank {
		var previousLog models.FuelLog
		// Procura o último abastecimento de tanque cheio deste veículo
		err := config.DB.Where("vehicle_id = ? AND is_full_tank = ?", vehicle.ID, true).
			Order("odometer desc").
			First(&previousLog).Error

		if err == nil {
			// Se encontrou um registo anterior, calcula a média
			kmPercorridos := newFuelLog.Odometer - previousLog.Odometer
			if kmPercorridos > 0 && newFuelLog.Liters > 0 {
				newFuelLog.KmLiter = float64(kmPercorridos) / newFuelLog.Liters
			}
		}
	}

	// 4. Salvar no Banco
	if err := config.DB.Create(&newFuelLog).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao registar abastecimento"})
		return
	}

	// 5. Atualizar preventivamente o odómetro atual do veículo se o novo odómetro for maior
	if newFuelLog.Odometer > vehicle.CurrentOdo {
		config.DB.Model(&vehicle).Update("current_odo", newFuelLog.Odometer)
	}

	c.JSON(http.StatusCreated, newFuelLog)
}

// GET /api/fuels?vehicle_id=UUID
func GetFuelLogsByVehicle(c *gin.Context) {
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

	// Validar posse do veículo
	var vehicle models.Vehicle
	if err := config.DB.Where("id = ? AND user_id = ?", vehicleID, userID).First(&vehicle).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Veículo não encontrado ou você não tem permissão"})
		return
	}

	var logs []models.FuelLog
	if err := config.DB.Where("vehicle_id = ?", vehicle.ID).Order("date desc").Find(&logs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar registos de combustível"})
		return
	}

	c.JSON(http.StatusOK, logs)
}

func GetFuelLogByID(c *gin.Context) {
	fuelID := c.Param("id")

	userIDContext, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao identificar usuário"})
		return
	}
	userID := userIDContext.(uuid.UUID)

	var fuelLog models.FuelLog
	if err := config.DB.Joins("JOIN vehicles ON vehicles.id = fuel_logs.vehicle_id").
		Where("fuel_logs.id = ? AND vehicles.user_id = ?", fuelID, userID).
		First(&fuelLog).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Registro de abastecimento não encontrado ou você não tem permissão"})
		return
	}

	c.JSON(http.StatusOK, fuelLog)
}

// PUT /api/fuels/:id
func UpdateFuelLog(c *gin.Context) {
	fuelID := c.Param("id")

	userIDContext, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao identificar usuário"})
		return
	}
	userID := userIDContext.(uuid.UUID)

	var fuelLog models.FuelLog
	// Segurança por Joins: Garante que o abastecimento pertence a um veículo do usuário logado
	if err := config.DB.Joins("JOIN vehicles ON vehicles.id = fuel_logs.vehicle_id").
		Where("fuel_logs.id = ? AND vehicles.user_id = ?", fuelID, userID).
		First(&fuelLog).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Registro de abastecimento não encontrado ou você não tem permissão"})
		return
	}

	var input struct {
		Date         time.Time `json:"date"`
		Odometer     int       `json:"odometer"`
		Liters       float64   `json:"liters"`
		PricePerLite float64   `json:"price_per_liter"`
		IsFullTank   *bool     `json:"is_full_tank"` // Usando ponteiro para detectar falso booleano enviado explicitamente
		GasStation   string    `json:"gas_station"`
		FuelType     string    `json:"fuel_type"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	// Atualizações parciais
	if !input.Date.IsZero() {
		fuelLog.Date = input.Date
	}
	if input.Odometer != 0 {
		fuelLog.Odometer = input.Odometer
	}
	if input.Liters != 0 {
		fuelLog.Liters = input.Liters
	}
	if input.PricePerLite != 0 {
		fuelLog.PricePerLit = input.PricePerLite
	}
	if input.IsFullTank != nil {
		fuelLog.IsFullTank = *input.IsFullTank
	}
	if input.GasStation != "" {
		fuelLog.GasStation = input.GasStation
	}
	if input.FuelType != "" {
		fuelLog.FuelType = input.FuelType
	}

	// Recalcula o custo total preventivamente se houver alterações
	fuelLog.TotalCost = fuelLog.Liters * fuelLog.PricePerLit

	// Reexecuta a lógica de consumo (KM/L) se o tanque continuar cheio
	if fuelLog.IsFullTank {
		var previousLog models.FuelLog
		err := config.DB.Where("vehicle_id = ? AND is_full_tank = ? AND odometer < ?", fuelLog.VehicleID, true, fuelLog.Odometer).
			Order("odometer desc").
			First(&previousLog).Error

		if err == nil {
			kmPercorridos := fuelLog.Odometer - previousLog.Odometer
			if kmPercorridos > 0 && fuelLog.Liters > 0 {
				fuelLog.KmLiter = float64(kmPercorridos) / fuelLog.Liters
			}
		} else {
			fuelLog.KmLiter = 0 // Zera se não houver anterior válido após a edição
		}
	} else {
		fuelLog.KmLiter = 0
	}

	if err := config.DB.Save(&fuelLog).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao atualizar abastecimento"})
		return
	}

	c.JSON(http.StatusOK, fuelLog)
}

// DELETE /api/fuels/:id
func DeleteFuelLog(c *gin.Context) {
	fuelID := c.Param("id")

	userIDContext, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao identificar usuário"})
		return
	}
	userID := userIDContext.(uuid.UUID)

	var fuelLog models.FuelLog
	if err := config.DB.Joins("JOIN vehicles ON vehicles.id = fuel_logs.vehicle_id").
		Where("fuel_logs.id = ? AND vehicles.user_id = ?", fuelID, userID).
		First(&fuelLog).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Registro de abastecimento não encontrado ou você não tem permissão"})
		return
	}

	if err := config.DB.Delete(&fuelLog).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao remover registro"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Registro de abastecimento removido com sucesso"})
}
