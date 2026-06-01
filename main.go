package main

import (
	"log"
	"net/http"
	"os"

	"garagefy-api/config"
	"garagefy-api/controllers"
	"garagefy-api/models"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Println("⚠️ Aviso: Arquivo .env não encontrado. Usando variáveis do sistema.")
	}

	config.ConnectDatabase()
	errorMigrate := config.DB.AutoMigrate(&models.Vehicle{}, &models.LogbookEntry{})
	if errorMigrate != nil {
		log.Fatal("Erro ao rodar as migrações: ", errorMigrate)
	}

	r := gin.Default()

	r.Static("/uploads", "./uploads")

	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})
	api := r.Group("/api")
	{
		// Rotas do Módulo de Veículos (CRUD Completo)
		api.POST("/vehicles", controllers.CreateVehicle)
		api.GET("/vehicles", controllers.GetVehicles)
		api.GET("/vehicles/:id", controllers.GetVehicleByID) // <-- Aqui usa :id
		api.PUT("/vehicles/:id", controllers.UpdateVehicle)
		api.DELETE("/vehicles/:id", controllers.DeleteVehicle)

		// Rotas do Módulo de Logbook (Ajustadas para evitar o conflito do Gin)
		api.POST("/vehicles/:id/logbook", controllers.CreateLogbookEntry)              // <-- Alterado de :vehicleId para :id
		api.GET("/vehicles/:id/logbook", controllers.GetLogbookEntries)                // <-- Alterado de :vehicleId para :id
		api.GET("/vehicles/:id/logbook/:logbookId", controllers.GetLogbookEntryByID)   // <-- Alterado para :id e :logbookId
		api.PUT("/vehicles/:id/logbook/:logbookId", controllers.UpdateLogbookEntry)    // <-- Alterado para :id e :logbookId
		api.DELETE("/vehicles/:id/logbook/:logbookId", controllers.DeleteLogbookEntry) // <-- Alterado para :id e :logbookId
	}

	appPort := os.Getenv("PORT")
	if appPort == "" {
		appPort = "8080"
	}

	r.Run(":" + appPort)
}
