package statuspage

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Index - render statuspage
func Index() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title": "Main website",
		})
	}
}
