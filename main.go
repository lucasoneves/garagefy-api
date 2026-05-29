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
	config.DB.AutoMigrate(&models.LogbookEntry{})

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
		api.POST("/vehicles/:vehicleId/logbook", controllers.CreateLogbookEntry)
		api.GET("/vehicles/:vehicleId/logbook", controllers.GetLogbookEntries)
		api.GET("/vehicles/:vehicleId/logbook/:id", controllers.GetLogbookEntryByID) // <-- Nova rota de busca única
		api.PUT("/vehicles/:vehicleId/logbook/:id", controllers.UpdateLogbookEntry)
		api.DELETE("/vehicles/:vehicleId/logbook/:id", controllers.DeleteLogbookEntry)
	}

	appPort := os.Getenv("PORT")
	if appPort == "" {
		appPort = "8080"
	}

	r.Run(":" + appPort)
}
