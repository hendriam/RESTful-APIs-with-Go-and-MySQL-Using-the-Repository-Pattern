package helper

import (
	"RESTful-APIs-with-Go-and-MySQL-Using-the-Repository-Pattern/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// SendSuccessResponse sends a successful response.
func SendSuccessResponse(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, models.ResponseSuccess{
		Code:    statusCode,
		Message: message,
		Data:    data,
	})
}

// SendErrorResponse sends an error response with an error message.
func SendErrorResponse(c *gin.Context, statusCode int, message string, errors interface{}) {
	c.JSON(statusCode, models.ResponseError{
		Code:    statusCode,
		Message: message,
		Errors:  errors,
	})
}

// SendValidationError Response sends validation error response.
func SendValidationErrorResponse(c *gin.Context, errors []models.ValidationErrorDetail) {
	c.JSON(http.StatusUnprocessableEntity, models.ResponseError{
		Code:    http.StatusUnprocessableEntity,
		Message: "validation error",
		Errors:  errors,
	})
}
