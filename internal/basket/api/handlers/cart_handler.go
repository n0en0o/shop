package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/n0en0o/shop/internal/basket/applications/commands"
	"github.com/n0en0o/shop/internal/basket/applications/queries"
	"github.com/n0en0o/shop/internal/basket/domain"
)

type CartHandler struct {
	saveCart   *commands.SaveCartHandler
	getCart    *queries.GetCartHandler
	removeCart *commands.RemoveCartHandler
}

func NewCartHandler(
	saveCart *commands.SaveCartHandler,
	getCart *queries.GetCartHandler,
	removeCart *commands.RemoveCartHandler,
) *CartHandler {
	return &CartHandler{
		saveCart:   saveCart,
		getCart:    getCart,
		removeCart: removeCart,
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
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"account_name": accountName,
		"location":     fmt.Sprintf("%s/%s", c.FullPath(), accountName),
	})
}

func (h *CartHandler) GetCart(c *gin.Context) {
	accountName := c.Param("accountName")
	cart, err := h.getCart.Handle(c.Request.Context(), accountName)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"result": cart,
	})
}

func (h *CartHandler) RemoveCart(c *gin.Context) {
	accountName := c.Param("accountName")
	success, err := h.removeCart.Handle(
		c.Request.Context(),
		accountName,
	)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": success,
	})
}
