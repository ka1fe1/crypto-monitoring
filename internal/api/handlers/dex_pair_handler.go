package handlers

import (
	"net/http"
	"strings"

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
// @Description  Get information about DEX pairs (supports comma-separated addresses)
// @Tags         dex
// @Accept       json
// @Produce      json
// @Param        contract_address  query     string  true  "Contract Addresses (comma-separated)"
// @Param        network_slug      query     string  false "Network Slug"
// @Param        network_id        query     string  false "Network ID"
// @Success      200  {array}   service.DexPairInfo
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/v1/dex/pair [get]
func (h *DexPairHandler) GetDexPair(c *gin.Context) {
	contractAddressQuery := c.Query("contract_address")
	networkSlug := c.Query("network_slug")
	networkId := c.Query("network_id")

	if contractAddressQuery == "" || (networkSlug == "" && networkId == "") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "contract_address and (network_slug or network_id) are required"})
		return
	}

	contractAddresses := strings.Split(contractAddressQuery, ",")
	for i := range contractAddresses {
		contractAddresses[i] = strings.TrimSpace(contractAddresses[i])
	}

	infos, err := h.service.GetDexPairInfo(contractAddresses, networkSlug, networkId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, infos)
}
