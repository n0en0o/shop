package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/n0en0o/marketplace/internal/catalog/applications/queries"
	"github.com/n0en0o/marketplace/internal/catalog/domain/spec"
)

type CatalogItemsHandlerV2 struct {
	catalogItems *queries.CatalogItemsV2Handler
}

func NewCatalogItemsHandlerV2(
	catalogItems *queries.CatalogItemsV2Handler,
) *CatalogItemsHandlerV2 {
	return &CatalogItemsHandlerV2{
		catalogItems: catalogItems,
	}
}

func (h *CatalogItemsHandlerV2) CatalogItems(c *gin.Context) {
	var args spec.QueryArgs
	if err := c.ShouldBindQuery(&args); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	result, err := h.catalogItems.Handle(c.Request.Context(), args)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"catalog_items": result,
	})
}
