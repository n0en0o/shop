package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/n0en0o/shop/internal/catalog/applications/commands"
	"github.com/n0en0o/shop/internal/catalog/applications/queries"
	"github.com/n0en0o/shop/internal/catalog/domain/repositories"
)

type CatalogItemsHandler struct {
	catalogItems        *queries.CatalogItemsHandler
	catalogItemByID     *queries.CatalogItemByIDHandler
	catalogItemsByTitle *queries.CatalogItemsByTitleHandler
	catalogItemsByBrand *queries.CatalogItemsByBrandHandler

	createCatalogItem *commands.CreateCatalogItemHandler
	updateCatalogItem *commands.UpdateCatalogItemHandler
	deleteCatalogItem *commands.DeleteCatalogItemHandler
}

func NewCatalogItemsHandler(
	catalogItems *queries.CatalogItemsHandler,
	catalogItemByID *queries.CatalogItemByIDHandler,
	catalogItemsByTitle *queries.CatalogItemsByTitleHandler,
	catalogItemsByBrand *queries.CatalogItemsByBrandHandler,

	createCatalogItem *commands.CreateCatalogItemHandler,
	updateCatalogItem *commands.UpdateCatalogItemHandler,
	deleteCatalogItem *commands.DeleteCatalogItemHandler,
) *CatalogItemsHandler {
	return &CatalogItemsHandler{
		catalogItems:        catalogItems,
		catalogItemByID:     catalogItemByID,
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (h *CatalogItemsHandler) CatalogItemByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid UUID"})
		return
	}

	item, err := h.catalogItemByID.Handle(c.Request.Context(), id)
	if errors.Is(err, repositories.ErrItemNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "catalog item not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"item": item})
}

func (h *CatalogItemsHandler) CatalogItemsByTitle(c *gin.Context) {
	title := c.Param("title")

	items, err := h.catalogItemsByTitle.Handle(c.Request.Context(), title)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (h *CatalogItemsHandler) CreateCatalogItem(c *gin.Context) {
	var cmd commands.CreateCatalogItemCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := h.createCatalogItem.Handle(c.Request.Context(), cmd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"created_catalog_item_id": id})

}

func (h *CatalogItemsHandler) UpdateCatalogItem(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid UUID"})
		return
	}

	var cmd commands.UpdateCatalogItemCommand

	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cmd.ID = id

	success, err := h.updateCatalogItem.Handle(c.Request.Context(), cmd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	if !success {
		c.JSON(http.StatusNotFound, gin.H{"error": "catalog item not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "catalog item updated successfully"})
}

func (h *CatalogItemsHandler) DeleteCatalogItem(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid UUID"})
		return
	}

	success, err := h.deleteCatalogItem.Handle(
		c.Request.Context(), commands.DeleteCatalogItemCommand{ID: id})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	if !success {
		c.JSON(http.StatusNotFound, gin.H{"error": "catalog item not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "catalog item deleted successfully"})
}

func (h *CatalogItemsHandler) CatalogItemsByBrand(c *gin.Context) {
	brand := c.Param("brand_title")

	items, err := h.catalogItemsByBrand.Handle(c.Request.Context(), brand)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"items": items})
}
