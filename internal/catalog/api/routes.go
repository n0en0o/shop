package api

import (
	"github.com/gin-gonic/gin"
	"github.com/n0en0o/marketplace/internal/catalog/api/handlers"
)

func RegisterRoutes(
	r *gin.Engine,
	brands *handlers.BrandsHandler,
	categories *handlers.CategoriesHandler,
	items *handlers.CatalogItemsHandler,
) {
	v1 := r.Group("/api/v1")

	v1.GET("/brands", brands.Brands)             // http://localhost:9001/api/v1/brands
	v1.GET("/categories", categories.Categories) // http://localhost:9001/api/v1/categories
	v1.GET("/catalog-items", items.CatalogItems) // http://localhost:9001/api/v1/catalog-items
}
