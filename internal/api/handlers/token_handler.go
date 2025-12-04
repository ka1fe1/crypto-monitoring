package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ka1fe1/crypto-monitoring/internal/service"
)

type TokenHandler struct {
	service service.TokenService
}

func NewTokenHandler(service service.TokenService) *TokenHandler {
	return &TokenHandler{
		service: service,
	}
}

// GetTokenPrice godoc
// @Summary      Get Token Price
// @Description  Get price of tokens by IDs
// @Tags         token
// @Accept       json
// @Produce      json
// @Param        ids  query     string  true  "Token IDs (comma separated)"
// @Success      200  {object}  map[string]utils.TokenInfo
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/v1/token/price [get]
func (h *TokenHandler) GetTokenPrice(c *gin.Context) {
	idsParam := c.Query("ids")
	if idsParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ids parameter is required"})
		return
	}

	ids := strings.Split(idsParam, ",")
	prices, err := h.service.GetTokenPrice(ids)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, prices)
}
