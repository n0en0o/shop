package handlers

import (
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
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"brands": brands})
}
