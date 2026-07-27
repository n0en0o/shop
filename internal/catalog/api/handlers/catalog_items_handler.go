package handlers

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/n0en0o/marketplace/internal/catalog/applications/commands"
	"github.com/n0en0o/marketplace/internal/catalog/applications/queries"
	"github.com/n0en0o/marketplace/internal/catalog/domain/repositories"
)

type CatalogItemsHandler struct {
	catalogItems        *queries.CatalogItemsHandler
	catalogItemById     *queries.CatalogItemByIdHandler
	catalogItemsByTitle *queries.CatalogItemsByTitleHandler
	catalogItemsByBrand *queries.CatalogItemsByBrandHandler

	createCatalogItem *commands.CreateCatalogItemHandler
	updateCatalogItem *commands.UpdateCatalogItemHandler
	deleteCatalogItem *commands.DeleteCatalogItemHandler
}

func NewCatalogItemsHandler(
	catalogItems *queries.CatalogItemsHandler,
	catalogItemById *queries.CatalogItemByIdHandler,
	catalogItemsByTitle *queries.CatalogItemsByTitleHandler,
	catalogItemsByBrand *queries.CatalogItemsByBrandHandler,

	createCatalogItem *commands.CreateCatalogItemHandler,
	updateCatalogItem *commands.UpdateCatalogItemHandler,
	deleteCatalogItem *commands.DeleteCatalogItemHandler,
) *CatalogItemsHandler {
	return &CatalogItemsHandler{
		catalogItems:        catalogItems,
		catalogItemById:     catalogItemById,
		catalogItemsByTitle: catalogItemsByTitle,
		catalogItemsByBrand: catalogItemsByBrand,

		createCatalogItem: createCatalogItem,
		updateCatalogItem: updateCatalogItem,
		deleteCatalogItem: deleteCatalogItem,
	}
}

func (h *CatalogItemsHandler) CatalogItems(c *gin.Context) {
	items, err := h.catalogItems.Handle(c.Request.Context())
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"catalog_items": items})
}

func (h *CatalogItemsHandler) CatalogItemById(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid UUID"})
		return
	}

	item, err := h.catalogItemById.Handle(c.Request.Context(), id)
	if errors.Is(err, repositories.ErrItemNotFound) {
		c.JSON(404, gin.H{"error": "catalog item not found"})
		return
	}
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"catalog_item": item})
}

func (h *CatalogItemsHandler) CatalogItemsByTitle(c *gin.Context) {
	title := c.Param("title")

	items, err := h.catalogItemsByTitle.Handle(c.Request.Context(), title)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if items == nil {
		c.JSON(404, gin.H{"error": "there is no catalog item with this title"})
		return
	}

	c.JSON(200, gin.H{"catalog_items": items})
}

func (h *CatalogItemsHandler) CreateCatalogItem(c *gin.Context) {
	var cmd commands.CreateCatalogItemCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	id, err := h.createCatalogItem.Handle(c.Request.Context(), cmd)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, gin.H{"created_catalog_item_id": id})

}

func (h *CatalogItemsHandler) UpdateCatalogItem(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid UUID"})
		return
	}

	var cmd commands.UpdateCatalogItemCommand

	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	cmd.ID = id

	success, err := h.updateCatalogItem.Handle(c.Request.Context(), cmd)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if !success {
		c.JSON(404, gin.H{"error": "catalog item not found"})
		return
	}

	c.JSON(200, gin.H{"message": "catalog item updated successfully"})
}

func (h *CatalogItemsHandler) DeleteCatalogItem(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid UUID"})
		return
	}

	success, err := h.deleteCatalogItem.Handle(
		c.Request.Context(), commands.DeleteCatalogItemCommand{ID: id})

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if !success {
		c.JSON(404, gin.H{"error": "catalog item not found"})
		return
	}

	c.JSON(200, gin.H{"message": "catalog item deleted successfully"})
}

func (h *CatalogItemsHandler) CatalogItemsByBrand(c *gin.Context) {
	brand := c.Param("brand_title")

	items, err := h.catalogItemsByBrand.Handle(c.Request.Context(), brand)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if items == nil {
		c.JSON(404, gin.H{"error": "there is no catalog item with this brand"})
		return
	}

	c.JSON(200, gin.H{"catalog_items": items})
}
