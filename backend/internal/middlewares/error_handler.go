package middlewares

import (
	"net/http"

	"github.com/ZAUakaAlexey/backend_go/internal/responses"
	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last()

			if c.Writer.Written() {
				return
			}

			switch err.Type {
			case gin.ErrorTypeBind:
				responses.ErrorResponse(c, http.StatusBadRequest, "Invalid request format", nil)
			case gin.ErrorTypePublic:
				responses.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
			default:
				responses.ErrorResponse(c, http.StatusInternalServerError, "Internal server error", nil)
			}
		}
	}
}

func NotFound() gin.HandlerFunc {
	return func(c *gin.Context) {
		errors := map[string][]string{
			"route": {"Route not found"},
		}
		responses.ErrorResponse(c, http.StatusNotFound, "Endpoint not found", errors)
	}
}
