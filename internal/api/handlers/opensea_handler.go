package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ka1fe1/crypto-monitoring/internal/service"
)

type OpenSeaHandler struct {
	service            service.OpenSeaService
	defaultCollections []string
}

func NewOpenSeaHandler(service service.OpenSeaService, defaultCollections []string) *OpenSeaHandler {
	return &OpenSeaHandler{
		service:            service,
		defaultCollections: defaultCollections,
	}
}

// GetNFTFloorPrice godoc
// @Summary      Get NFT Floor Price
// @Description  Get floor price of NFT collections
// @Tags         nft
// @Accept       json
// @Produce      json
// @Param        slugs           query     string  false  "Collection Slugs (comma separated)"
// @Param        convert_to_usd  query     bool    false  "Convert to USD"
// @Success      200             {array}   service.NFTFloorPriceInfo
// @Failure      400             {object}  map[string]string
// @Failure      500             {object}  map[string]string
// @Router       /api/v1/nft/floor_price [get]
func (h *OpenSeaHandler) GetNFTFloorPrice(c *gin.Context) {
	slugsParam := c.Query("slugs")
	var slugs []string
	if slugsParam == "" {
		slugs = h.defaultCollections
	} else {
		slugs = strings.Split(slugsParam, ",")
	}

	// Trim spaces and filter empty strings
	var validSlugs []string
	for _, s := range slugs {
		trimmed := strings.TrimSpace(s)
		if trimmed != "" {
			validSlugs = append(validSlugs, trimmed)
		}
	}

	if len(validSlugs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no collections specified and no defaults configured"})
		return
	}

	convertToUsd := c.Query("convert_to_usd") == "true"

	prices, err := h.service.GetNFTFloorPrices(validSlugs, convertToUsd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, prices)
}
