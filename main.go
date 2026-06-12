package main

import (
	"log"
	"net/http"
	"os"

	"garagefy-api/config"
	"garagefy-api/controllers"
	"garagefy-api/middlewares" // Importa o pacote onde criamos o middleware
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
	errorMigrate := config.DB.AutoMigrate(&models.Vehicle{}, &models.LogbookEntry{}, &models.Service{}, &models.User{}, &models.FuelLog{})
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
		// 1. ROTAS PÚBLICAS (NÃO exigem Token JWT)
		auth := api.Group("/auth")
		{
			auth.POST("/register", controllers.Register)
			auth.POST("/login", controllers.Login)
		}

		// 2. ROTAS PROTEGIDAS (Exigem o cabeçalho 'Authorization: Bearer <TOKEN>')
		protected := api.Group("/")
		protected.Use(middlewares.AuthMiddleware()) // Ativa a tranca JWT para o bloco abaixo
		{
			// Rotas do Módulo de Veículos (CRUD Completo)
			protected.POST("/vehicles", controllers.CreateVehicle)
			protected.GET("/vehicles", controllers.GetVehicles)
			protected.GET("/vehicles/:id", controllers.GetVehicleByID)
			protected.PUT("/vehicles/:id", controllers.UpdateVehicle)
			protected.DELETE("/vehicles/:id", controllers.DeleteVehicle)

			// Rotas do Módulo de Logbook
			protected.POST("/vehicles/:id/logbook", controllers.CreateLogbookEntry)
			protected.GET("/vehicles/:id/logbook", controllers.GetLogbookEntries)
			protected.GET("/vehicles/:id/logbook/:logbookId", controllers.GetLogbookEntryByID)
			protected.PUT("/vehicles/:id/logbook/:logbookId", controllers.UpdateLogbookEntry)
			protected.DELETE("/vehicles/:id/logbook/:logbookId", controllers.DeleteLogbookEntry)

			// Rotas do Módulo de Serviço
			protected.POST("/services", controllers.CreateService)
			protected.GET("/services", controllers.GetServicesByVehicle)
			protected.GET("/services/:id", controllers.GetServiceByID)
			protected.PUT("/services/:id", controllers.UpdateService)
			protected.DELETE("/services/:id", controllers.DeleteService)

			// Novas Rotas de Abastecimento
			protected.POST("/fuels", controllers.CreateFuelLog)
			protected.GET("/fuels", controllers.GetFuelLogsByVehicle)
			protected.GET("/fuels/:id", controllers.GetFuelLogByID)
			protected.PUT("/fuels/:id", controllers.UpdateFuelLog)
			protected.DELETE("/fuels/:id", controllers.DeleteFuelLog)
		}
	}

	appPort := os.Getenv("PORT")
	if appPort == "" {
		appPort = "8080"
	}

	r.Run(":" + appPort)
}
