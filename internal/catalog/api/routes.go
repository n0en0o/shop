package api

import (
	"github.com/gin-gonic/gin"
	"github.com/n0en0o/marketplace/internal/catalog/api/handlers"
)

func RegisterRoutes(
	r *gin.Engine,
	brands *handlers.BrandsHandler,
) {
	v1 := r.Group("/api/v1")

	v1.GET("/brands", brands.Brands) // http://localhost:9001/api/v1/brands
}
