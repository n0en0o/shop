package handlers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/n0en0o/shop/internal/checkout/applications/queries"
	"github.com/n0en0o/shop/internal/checkout/domain"
	"github.com/n0en0o/shop/internal/shared"
)

type OrderHandler struct {
	orderByIDQueryHandler           *queries.OrderByIDQueryHandler
	ordersByAccountNameQueryHandler *queries.OrdersByAccountNameQueryHandler
}

func NewOrderHandler(
	orderByID *queries.OrderByIDQueryHandler,
	ordersByAccountName *queries.OrdersByAccountNameQueryHandler,
) *OrderHandler {
	return &OrderHandler{
		orderByIDQueryHandler:           orderByID,
		ordersByAccountNameQueryHandler: ordersByAccountName,
	}
}

func (h *OrderHandler) Order(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "неверный id заказа",
		})
		return
	}

	query := queries.OrderByIDQuery{ID: id}
	order, err := h.orderByIDQueryHandler.Handle(c.Request.Context(), query)
	if err != nil {
		var notFound *shared.NotFoundError
		if errors.As(err, &notFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "заказ не существует",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if order == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "заказ не существует",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": order,
	})
}

func (h *OrderHandler) OrdersByAccountName(c *gin.Context) {
	accountName := strings.TrimSpace(c.Param("name"))
	if accountName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "имя аккаунта обязательно",
		})
		return
	}

	query := queries.OrdersByAccountNameQuery{AccountName: accountName}
	orders, err := h.ordersByAccountNameQueryHandler.Handle(c.Request.Context(), query)
	if err != nil {
		var notFound *shared.NotFoundError
		if errors.As(err, &notFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "заказы аккаунта не найдены",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if orders == nil {
		orders = make([]*domain.Order, 0)
	}

	c.JSON(http.StatusOK, gin.H{
		"result": orders,
	})
}
