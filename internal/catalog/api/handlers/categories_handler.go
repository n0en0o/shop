package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/n0en0o/marketplace/internal/catalog/applications/queries"
)

type CategoriesHandler struct {
	categories *queries.CategoriesHandler
}

func NewCategoriesHandler(categories *queries.CategoriesHandler) *CategoriesHandler {
	return &CategoriesHandler{
		categories: categories,
	}
}

func (h *CategoriesHandler) Categories(c *gin.Context) {
	categories, err := h.categories.Handle(c.Request.Context())
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"categories": categories})
}
