package main

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/meritlabs/overlord/overseer/controllers"
)

func main() {
	fmt.Printf("Overseer is.\n")

	r := gin.Default()

	r.GET("/ping", controllers.Ping)
	r.GET("/check", controllers.Status)

	r.Run()
}
