package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/n0en0o/marketplace/internal/catalog/applications/queries"
)

type BrandsHandler struct {
	brands *queries.BrandsHandler
}

func NewBrandsHandler(brands *queries.BrandsHandler) *BrandsHandler {
	return &BrandsHandler{
		brands: brands,
	}
}

func (h *BrandsHandler) Brands(c *gin.Context) {
	brands, err := h.brands.Handle(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"brands": brands})
}
