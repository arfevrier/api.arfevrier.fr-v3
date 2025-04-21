package main

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

// bitcoinPrice godoc
// @Summary Get Bitcoin price
// @Schemes
// @Description Fetch the current Bitcoin price from an blockchain.info API
// @Tags bitcoin
// @Accept json
// @Produce json
// @Success 200 {string} string "API result"
// @Failure 500 {string} string "Failed to fetch Bitcoin price"
// @Router /bitcoin/price [get]
func (Api API) bitcoinPrice(c *gin.Context) {
	resp, err := http.Get("https://blockchain.info/ticker")
	if err != nil {
		c.JSON(500, gin.H{"msg": "Failed to fetch Bitcoin price"})
		return
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		c.JSON(500, gin.H{"msg": "Failed to parse response"})
		return
	}

	c.JSON(http.StatusOK, result)
}
