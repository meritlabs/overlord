package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/meritlabs/overlord/blockchain"
)

func Status(c *gin.Context, mode string) {
	output, err := blockchain.GetBlockchainInfo(mode)

	var data blockchain.CheckResponse
	if output != nil {
		data = blockchain.CheckResponse{*output, ""}
	} else {
		data = blockchain.CheckResponse{blockchain.BlockchainInfo{}, err.Error()}
	}

	c.JSON(http.StatusOK, data)
}
