package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/meritlabs/overlord/overseer/blockchain"
)

func Status(c *gin.Context) {
	output, err := blockchain.GetBlockchainInfo()

	c.JSON(http.StatusOK,
		gin.H{
			"status": output,
			"error":  err,
		},
	)
}
