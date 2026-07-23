package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/n0en0o/marketplace/internal/catalog/applications/queries"
)

type CatalogItemsHandler struct {
	catalogItems *queries.CatalogItemsHandler
}

func NewCatalogItemsHandler(catalogItems *queries.CatalogItemsHandler) *CatalogItemsHandler {
	return &CatalogItemsHandler{
		catalogItems: catalogItems,
	}
}

func (h *CatalogItemsHandler) CatalogItems(c *gin.Context) {
	items, err := h.catalogItems.Handle(c.Request.Context())
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"catalog items": items})
}
