package controllers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func UploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Arquivo não enviado"})
		return
	}

	ext := filepath.Ext(file.Filename)
	filename := uuid.New().String() + ext
	savePath := filepath.Join("uploads", filename)

	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao abrir arquivo"})
		return
	}
	defer src.Close()

	dst, err := os.Create(savePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao salvar arquivo"})
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao escrever arquivo"})
		return
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	url := fmt.Sprintf("http://localhost:%s/uploads/%s", port, filename)

	c.JSON(http.StatusCreated, gin.H{
		"url":       url,
		"filename":  filename,
		"size":      file.Size,
		"createdAt": time.Now(),
	})
}
