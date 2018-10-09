package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/meritlabs/overlord/pkg/blockchain"
)

// Status - retruns blockchain status info as HTTP response
func Status(mode string) gin.HandlerFunc {
	return func(c *gin.Context) {
		output, err := blockchain.GetBlockchainInfo(mode)

		var data blockchain.CheckResponse
		if output != nil {
			data = blockchain.CheckResponse{*output, ""}
		} else {
			data = blockchain.CheckResponse{blockchain.BlockchainInfo{}, err.Error()}
		}

		c.JSON(http.StatusOK, data)
	}
}
