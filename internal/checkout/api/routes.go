package api

import (
	"github.com/gin-gonic/gin"
	"github.com/n0en0o/shop/internal/checkout/api/handlers"
)

func RegisterRoutes(r *gin.Engine, orderHandler *handlers.OrderHandler) {
	v1 := r.Group("/api/v1")
	{
		orders := v1.Group("/orders")
		{
			orders.GET("/account/:name", orderHandler.OrdersByAccountName)
			orders.GET("/:id", orderHandler.Order)
		}
	}
}
