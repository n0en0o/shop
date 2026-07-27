package handlers

import (
	"net/http"

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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"categories": categories})
}
