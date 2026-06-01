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

// POST /api/vehicles/:id/logbook
func CreateLogbookEntry(c *gin.Context) {
	vehicleID := c.Param("id") // Correto: mapeia o :id do veículo

	vid, err := uuid.Parse(vehicleID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID do veículo inválido"})
		return
	}

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
		VehicleID:     vid,
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

// GET /api/vehicles/:id/logbook
func GetLogbookEntries(c *gin.Context) {
	vehicleID := c.Param("id") // Correto: mapeia o :id do veículo
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

// GET /api/vehicles/:id/logbook/:logbookId
func GetLogbookEntryByID(c *gin.Context) {
	vehicleID := c.Param("id")        // Ajustado de "vehicleId" para "id"
	logbookID := c.Param("logbookId") // Ajustado de "id" para "logbookId"

	var entry models.LogbookEntry
	if err := config.DB.Where("id = ? AND vehicle_id = ?", logbookID, vehicleID).First(&entry).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Registro não encontrado"})
		return
	}

	c.JSON(http.StatusOK, entry)
}

// PUT /api/vehicles/:id/logbook/:logbookId
func UpdateLogbookEntry(c *gin.Context) {
	vehicleID := c.Param("id")        // Mapeia o :id da URL
	logbookID := c.Param("logbookId") // Mapeia o :logbookId da URL

	// Validação de segurança para garantir que os parâmetros não vieram vazios
	if vehicleID == "" || logbookID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "IDs ausentes na requisição"})
		return
	}

	var entry models.LogbookEntry
	// Busca o registro original garantindo que as chaves batem perfeitamente
	if err := config.DB.Where("id = ? AND vehicle_id = ?", logbookID, vehicleID).First(&entry).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Registro de manutenção não encontrado"})
		return
	}

	// Captura os dados textuais enviados via form-data para atualização
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

	// Como a struct 'entry' agora foi populada corretamente pelo First, o Save executará um UPDATE real
	if err := config.DB.Save(&entry).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao atualizar o registro no banco"})
		return
	}

	c.JSON(http.StatusOK, entry)
}

// DELETE /api/vehicles/:id/logbook/:logbookId
func DeleteLogbookEntry(c *gin.Context) {
	vehicleID := c.Param("id")        // Mapeia o :id da URL
	logbookID := c.Param("logbookId") // Mapeia o :logbookId da URL

	if vehicleID == "" || logbookID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "IDs ausentes na requisição"})
		return
	}

	var entry models.LogbookEntry
	// Localiza o registro antes para ter certeza de que as amarrações de ID existem
	if err := config.DB.Where("id = ? AND vehicle_id = ?", logbookID, vehicleID).First(&entry).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Registro de manutenção não encontrado"})
		return
	}

	// Deleta passando a struct populada com a chave primária real obtida no passo anterior
	if err := config.DB.Delete(&entry).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao deletar o registro"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Registro removido com sucesso"})
}
