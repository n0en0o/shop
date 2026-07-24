package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/n0en0o/marketplace/internal/catalog/applications/queries"
)

type CatalogItemsHandler struct {
	catalogItems    *queries.CatalogItemsHandler
	catalogItemById *queries.CatalogItemByIdHandler
}

func NewCatalogItemsHandler(
	catalogItems *queries.CatalogItemsHandler,
	catalogItemById *queries.CatalogItemByIdHandler,
) *CatalogItemsHandler {
	return &CatalogItemsHandler{
		catalogItems:    catalogItems,
		catalogItemById: catalogItemById,
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

func (h *CatalogItemsHandler) CatalogItemById(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid UUID"})
		return
	}

	item, err := h.catalogItemById.Handle(c.Request.Context(), id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if item == nil {
		c.JSON(404, gin.H{"error": "catalog item not found"})
		return
	}

	c.JSON(200, gin.H{"catalog item": item})
}
