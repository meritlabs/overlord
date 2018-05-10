package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/meritlabs/overlord/blockchain"
)

func Version(c *gin.Context, mode string) {
	output, err := blockchain.GetInfo(mode)

	var data blockchain.VersionResponse
	if output != nil {
		data = blockchain.VersionResponse{Info: *output, Error: ""}
	} else {
		data = blockchain.VersionResponse{Info: blockchain.VersionInfo{}, Error: err.Error()}
	}

	c.JSON(http.StatusOK, data)
}
