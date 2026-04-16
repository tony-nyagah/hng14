package main

import (
	"github.com/gin-gonic/gin"
)

type UserInfo struct {
	Name   string `json:"name"`
	Email  string `json:"email"`
	Github string `json:"github"`
}

var myInfo = UserInfo{
	Name:   "Antony Nyagah",
	Email:  "tony.m.nyagah@gmail.com",
	Github: "https://github.com/tony-nyagah",
}

func indexHandler(c *gin.Context) {
	c.JSON(200, gin.H{"message": "API is running"})
}

func healthHandler(c *gin.Context) {
	c.JSON(200, gin.H{"message": "healthy"})
}

func infoHandler(c *gin.Context) {
	c.JSON(200, myInfo)
}

func main() {
	r := gin.Default()
	r.GET("/", indexHandler)
	r.GET("/info", infoHandler)
	r.GET("/health", healthHandler)
	r.Run()
}
