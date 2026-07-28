package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/n0en0o/marketplace/internal/basket/applications/commands"
	"github.com/n0en0o/marketplace/internal/basket/domain"
)

type CartHandler struct {
	saveCart *commands.SaveCartHandler
}

func NewCartHandler(
	saveCart *commands.SaveCartHandler,
) *CartHandler {
	return &CartHandler{
		saveCart: saveCart,
	}
}

func (h *CartHandler) SaveCart(c *gin.Context) {
	var req struct {
		Cart domain.ShoppingCart `json:"cart"`
	}

	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	accountName, err := h.saveCart.Handle(
		c.Request.Context(),
		commands.SaveCartCommand{Cart: req.Cart},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"account_name": accountName,
		"location":     fmt.Sprintf("%s/%s", c.FullPath(), accountName),
	})
}
