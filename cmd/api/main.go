package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/victor-lima-142/oak-bank/pkg/config"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("warning: no .env file loaded: %v", err)
	}

	_, err := config.OpenDB()
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = os.Getenv("APP_PORT")
	}
	if port == "" {
		port = "8080"
	}

	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "oak-bank API"})
	})

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	log.Printf("starting server on :%s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
