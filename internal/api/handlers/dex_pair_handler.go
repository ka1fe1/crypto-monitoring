package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ka1fe1/crypto-monitoring/internal/service"
)

type DexPairHandler struct {
	service service.DexPairService
}

func NewDexPairHandler(service service.DexPairService) *DexPairHandler {
	return &DexPairHandler{
		service: service,
	}
}

// GetDexPair godoc
// @Summary      Get Dex Pair Info
// @Description  Get information about a DEX pair
// @Tags         dex
// @Accept       json
// @Produce      json
// @Param        contract_address  query     string  true  "Contract Address"
// @Param        network_slug      query     string  true  "Network Slug"
// @Success      200  {object}  service.DexPairInfo
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/v1/dex/pair [get]
func (h *DexPairHandler) GetDexPair(c *gin.Context) {
	contractAddress := c.Query("contract_address")
	networkSlug := c.Query("network_slug")

	if contractAddress == "" || networkSlug == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "contract_address and network_slug are required"})
		return
	}

	info, err := h.service.GetDexPairInfo(contractAddress, networkSlug)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, info)
}
