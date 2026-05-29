package controllers

import (
	"fmt"
	"net/http"
	"path/filepath"

	"garagefy-api/config"
	"garagefy-api/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// POST /api/vehicles/:vehicleId/logbook
func CreateLogbookEntry(c *gin.Context) {
	vehicleID := c.Param("vehicleId")

	title := c.PostForm("title")
	description := c.PostForm("description")
	category := c.PostForm("category")

	if title == "" || description == "" || category == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Campos de texto obrigatórios estão ausentes"})
		return
	}

	var attachmentURL *string

	file, err := c.FormFile("file")
	if err == nil {
		fileExt := filepath.Ext(file.Filename)
		newFileName := fmt.Sprintf("%s%s", uuid.New().String(), fileExt)

		uploadPath := filepath.Join("uploads", newFileName)
		if err := c.SaveUploadedFile(file, uploadPath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao salvar o anexo"})
			return
		}

		url := fmt.Sprintf("/uploads/%s", newFileName)
		attachmentURL = &url
	}

	entry := models.LogbookEntry{
		VehicleID:     vehicleID,
		Category:      models.LogbookCategory(category),
		Title:         title,
		Description:   description,
		AttachmentURL: attachmentURL,
	}

	if err := config.DB.Create(&entry).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao gravar no banco de dados"})
		return
	}

	c.JSON(http.StatusCreated, entry)
}

// GET /api/vehicles/:vehicleId/logbook
func GetLogbookEntries(c *gin.Context) {
	vehicleID := c.Param("vehicleId")
	categoryFilter := c.Query("category")

	var entries []models.LogbookEntry
	query := config.DB.Where("vehicle_id = ?", vehicleID)

	if categoryFilter != "" && categoryFilter != "all" {
		query = query.Where("category = ?", categoryFilter)
	}

	if err := query.Order("created_at desc").Find(&entries).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao buscar registros do Logbook"})
		return
	}

	c.JSON(http.StatusOK, entries)
}

func UpdateLogbookEntry(c *gin.Context) {
	id := c.Param("id")
	vehicleID := c.Param("vehicleId")

	var entry models.LogbookEntry
	// Busca o registro original garantindo que pertence ao veículo correto
	if err := config.DB.Where("id = ? AND vehicle_id = ?", id, vehicleID).First(&entry).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Registro não encontrado"})
		return
	}

	// Captura os dados textuais enviados para atualização
	title := c.PostForm("title")
	description := c.PostForm("description")
	category := c.PostForm("category")

	if title != "" {
		entry.Title = title
	}
	if description != "" {
		entry.Description = description
	}
	if category != "" {
		entry.Category = models.LogbookCategory(category)
	}

	// Verifica se foi enviado um novo arquivo para substituir o anexo
	file, err := c.FormFile("file")
	if err == nil {
		fileExt := filepath.Ext(file.Filename)
		newFileName := fmt.Sprintf("%s%s", uuid.New().String(), fileExt)
		uploadPath := filepath.Join("uploads", newFileName)

		if err := c.SaveUploadedFile(file, uploadPath); err == nil {
			url := fmt.Sprintf("/uploads/%s", newFileName)
			entry.AttachmentURL = &url
		}
	}

	if err := config.DB.Save(&entry).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao atualizar o registro"})
		return
	}

	c.JSON(http.StatusOK, entry)
}

// DELETE /api/vehicles/:vehicleId/logbook/:id
func DeleteLogbookEntry(c *gin.Context) {
	id := c.Param("id")
	vehicleID := c.Param("vehicleId")

	var entry models.LogbookEntry
	if err := config.DB.Where("id = ? AND vehicle_id = ?", id, vehicleID).First(&entry).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Registro não encontrado"})
		return
	}

	if err := config.DB.Delete(&entry).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao deletar o registro"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Registro removido com sucesso"})
}

// GET /api/vehicles/:vehicleId/logbook/:id
func GetLogbookEntryByID(c *gin.Context) {
	id := c.Param("id")
	vehicleID := c.Param("vehicleId")

	var entry models.LogbookEntry
	if err := config.DB.Where("id = ? AND vehicle_id = ?", id, vehicleID).First(&entry).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Registro não encontrado"})
		return
	}

	c.JSON(http.StatusOK, entry)
}
