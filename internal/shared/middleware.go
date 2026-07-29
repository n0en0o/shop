package shared

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 || c.Writer.Written() {
			return
		}

		err := c.Errors.Last().Err

		var validationErrs validator.ValidationErrors
		var notFound *NotFoundError

		switch {
		case errors.As(err, &validationErrs):
			fields := make([]gin.H, 0, len(validationErrs))
			for _, fe := range validationErrs {
				fields = append(fields, gin.H{
					"field":   fe.Field(),
					"message": fe.Error(),
				})
			}

			c.JSON(http.StatusBadRequest, gin.H{
				"title":            "Ошибка валидации",
				"details":          "Ошибка валидации запроса",
				"status":           http.StatusBadRequest,
				"validationErrors": fields,
			})

		case errors.As(err, &notFound):
			c.JSON(http.StatusNotFound, gin.H{
				"title":   "Ресурс не найден",
				"details": notFound.Error(),
				"status":  http.StatusNotFound,
			})

		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"title":   "Внутренняя ошибка сервера",
				"details": "Произошла непредвиденная ошибка",
				"status":  http.StatusInternalServerError,
			})
		}
	}
}
