package api

import (
	"github.com/gin-gonic/gin"
	"github.com/n0en0o/marketplace/internal/basket/api/handlers"
)

func RegisterRoutes(
	r *gin.Engine,
	cartHandler *handlers.CartHandler,
) {
	v1 := r.Group("/api/v1")
	v1.POST("/cart", cartHandler.SaveCart) // http://localhost:9002/api/v1/cart
}
