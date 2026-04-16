package main

import (
	"github.com/gin-gonic/gin"
)

func indexHandler(c *gin.Context) {
	c.JSON(200, gin.H{"message": "API is running"})
}

func healthHandler(c *gin.Context) {
	c.JSON(200, gin.H{"message": "healthy"})
}

func infoHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"name":   "Antony Nyagah",
		"email":  "tony.m.nyagah@gmail.com",
		"github": "https://github.com/tony-nyagah",
	})
}

func main() {
	r := gin.Default()
	r.GET("/", indexHandler)
	r.GET("/info", infoHandler)
	r.GET("/health", healthHandler)
	r.Run()
}
