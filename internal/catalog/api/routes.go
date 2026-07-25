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

	itemsV2 *handlers.CatalogItemsHandlerV2,
) {
	v1 := r.Group("/api/v1")

	v1.GET("/brands", brands.Brands)                                       // http://localhost:9001/api/v1/brands
	v1.GET("/categories", categories.Categories)                           // http://localhost:9001/api/v1/categories
	v1.GET("/catalog-items", items.CatalogItems)                           // http://localhost:9001/api/v1/catalog-items
	v1.GET("/catalog-items/:id", items.CatalogItemById)                    // http://localhost:9001/api/v1/catalog-items/:id
	v1.GET("/catalog-items/title/:title", items.CatalogItemsByTitle)       // http://localhost:9001/api/v1/catalog-items/title/:title
	v1.GET("/catalog-items/brand/:brand_title", items.CatalogItemsByBrand) // http://localhost:9001/api/v1/catalog-items/brand/:brand_title

	v1.POST("/catalog-items", items.CreateCatalogItem)       // http://localhost:9001/api/v1/catalog-items
	v1.PUT("/catalog-items/:id", items.UpdateCatalogItem)    // http://localhost:9001/api/v1/catalog-items/:id
	v1.DELETE("/catalog-items/:id", items.DeleteCatalogItem) // http://localhost:9001/api/v1/catalog-items/:id

	v2 := r.Group("/api/v2")

	v2.GET("/catalog-items", itemsV2.CatalogItems) // http://localhost:9001/api/v2/catalog-items

}
