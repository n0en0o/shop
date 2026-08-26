package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
	"github.com/n0en0o/shop/internal/basket/applications/commands"
	"github.com/n0en0o/shop/internal/basket/applications/queries"
	"github.com/n0en0o/shop/internal/basket/domain"
)

type CartHandler struct {
	saveCart     *commands.SaveCartHandler
	getCart      *queries.GetCartHandler
	removeCart   *commands.RemoveCartHandler
	checkoutCart *commands.CheckoutCartHandler
}

func NewCartHandler(
	saveCart *commands.SaveCartHandler,
	getCart *queries.GetCartHandler,
	removeCart *commands.RemoveCartHandler,
	checkoutCart *commands.CheckoutCartHandler,
) *CartHandler {
	return &CartHandler{
		saveCart:     saveCart,
		getCart:      getCart,
		removeCart:   removeCart,
		checkoutCart: checkoutCart,
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

func (h *CartHandler) CheckoutCart(c *gin.Context) {
	var req commands.CheckoutCartRequest

	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	correlationID := uuid.New().String()

	cmd := commands.CheckoutCartCommand{
		CheckoutCartRequest: req,
		CorrelationID:       correlationID,
	}

	result, err := h.checkoutCart.Handle(c.Request.Context(), cmd)

	if err != nil {
		if isValidationError(err) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		if strings.Contains(err.Error(), "не найдена") ||
			strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	response := commands.CheckoutCartResponse{
		OrderID:       result.OrderID,
		CorrelationID: result.CorrelationID,
		Message:       "Корзина успешно оформлена",
	}

	c.JSON(http.StatusOK, gin.H{
		"result": response,
	})

}

func isValidationError(err error) bool {
	errStr := err.Error()
	return strings.Contains(errStr, "validate") ||
		strings.Contains(errStr, "required") ||
		strings.Contains(errStr, "обязателен")
}
