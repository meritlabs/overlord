package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/meritlabs/overlord/pkg/blockchain"
)

// Status - retruns daemon version info as HTTP response
func Version(mode string) gin.HandlerFunc {
	return func(c *gin.Context) {
		output, err := blockchain.GetInfo(mode)

		var data blockchain.VersionResponse
		if output != nil {
			data = blockchain.VersionResponse{Info: *output, Error: ""}
		} else {
			data = blockchain.VersionResponse{Info: blockchain.VersionInfo{}, Error: err.Error()}
		}

		c.JSON(http.StatusOK, data)
	}
}
