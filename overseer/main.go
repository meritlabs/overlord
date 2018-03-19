package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Printf("Overseer is.\n")

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.GET("/check", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "check succeeded",
		})
	})
	r.Run()
}
