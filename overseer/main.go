package main

import (
	"flag"
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/meritlabs/overlord/overseer/controllers"
)

func main() {
	fmt.Printf("Overseer is deployed.\n")

	var mode string

	flag.StringVar(&mode, "mode", "mainnet", "Daemon mode. Can be mainnet, testnet or regtest")
	flag.Parse()

	fmt.Printf("Daemon mode: %s \n\n", mode)

	r := gin.Default()

	r.GET("/ping", controllers.Ping)
	r.GET("/check", func(c *gin.Context) {
		controllers.Status(c, mode)
	})
	r.GET("/version", func(c *gin.Context) {
		controllers.Version(c, mode)
	})

	r.Run()
}
