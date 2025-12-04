package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// PingHandler responds with a pong message
// PingHandler godoc
// @Summary      Ping
// @Description  Check server health
// @Tags         health
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]string
// @Router       /ping [get]
func PingHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
