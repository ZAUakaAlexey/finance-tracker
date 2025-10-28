package responses

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Tag     string `json:"tag,omitempty"`
}

type NoPaginatedResponse struct {
	Data struct {
		Resource interface{} `json:"resource"`
	} `json:"data"`
	Message *string             `json:"message"`
	Errors  map[string][]string `json:"errors"`
}

type PaginatedResponse struct {
	Data struct {
		Resource struct {
			Items interface{}    `json:"items"`
			Meta  PaginationMeta `json:"meta"`
		} `json:"resource"`
	} `json:"data"`
	Message *string             `json:"message"`
	Errors  map[string][]string `json:"errors"`
}

type PaginationMeta struct {
	CurrentPage int `json:"current_page"`
	From        int `json:"from"`
	LastPage    int `json:"last_page"`
	PerPage     int `json:"per_page"`
	To          int `json:"to"`
	Total       int `json:"total"`
}

func SuccessResponse(c *gin.Context, statusCode int, resource interface{}, message string) {
	response := NoPaginatedResponse{
		Message: &message,
		Errors:  make(map[string][]string),
	}
	response.Data.Resource = resource

	c.JSON(statusCode, response)
}

func ErrorResponse(c *gin.Context, statusCode int, message string, errors map[string][]string) {
	if errors == nil {
		errors = make(map[string][]string)
	}

	response := NoPaginatedResponse{
		Message: &message,
		Errors:  errors,
	}
	response.Data.Resource = nil

	c.JSON(statusCode, response)
}

func ValidationErrorResponse(c *gin.Context, message string, validationErrors []ValidationError) {
	errors := make(map[string][]string)

	for _, err := range validationErrors {
		errors[err.Field] = append(errors[err.Field], err.Message)
	}

	ErrorResponse(c, http.StatusBadRequest, message, errors)
}

func PaginatedSuccessResponse(c *gin.Context, statusCode int, items interface{}, meta PaginationMeta, message string) {
	response := PaginatedResponse{
		Message: &message,
		Errors:  make(map[string][]string),
	}
	response.Data.Resource.Items = items
	response.Data.Resource.Meta = meta

	c.JSON(statusCode, response)
}
