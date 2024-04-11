package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	envValueMbUsedToCreatePizza = "MB_USED_TO_CREATE_PIZZA"
)

func main() {
	router := gin.Default()

	sizeMbStr := os.Getenv(envValueMbUsedToCreatePizza)
	sizeMb, err := strconv.Atoi(sizeMbStr)

	if err != nil {
		// Use default value (10 mo) if environment variable is not set or error parsing
		fmt.Printf("error parsing %s => %s", envValueMbUsedToCreatePizza, err)
		sizeMb = 10
	}

	router.POST("/order/pizza", func(c *gin.Context) {
		simulateHeavyMemoryWork(sizeMb)

		c.JSON(http.StatusCreated, gin.H{
			"message": "Pizza order successful",
		})
	})

	router.Run(":8080")
}

func simulateHeavyMemoryWork(sizeMB int) {
	// Calculate the size in bytes
	sizeBytes := sizeMB * 1024 * 1024

	memorySlice := make([]byte, sizeBytes)

	// Perform some operation on the memory slice
	for i := 0; i < len(memorySlice); i++ {
		memorySlice[i] = byte(i % 256)
	}

	// Simulate some heavy computation with the memory
	for i := 0; i < sizeBytes; i++ {
		memorySlice[i] = memorySlice[i] + 1
	}
}
