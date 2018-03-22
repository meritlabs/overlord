package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/meritlabs/overlord/blockchain"
)

func Status(c *gin.Context) {
	output, err := blockchain.GetBlockchainInfo()

	var data blockchain.CheckResponse
	if output != nil {
		data = blockchain.CheckResponse{*output, err}
	} else {
		data = blockchain.CheckResponse{blockchain.BlockchainInfo{}, err}
	}

	c.JSON(http.StatusOK, data)
}
