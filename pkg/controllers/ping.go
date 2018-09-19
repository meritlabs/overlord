package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Ping - healthcheck or heartbeat
func Ping() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	}
}
