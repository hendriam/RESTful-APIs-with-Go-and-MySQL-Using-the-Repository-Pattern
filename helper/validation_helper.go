package helper

import (
	"RESTful-APIs-with-Go-and-MySQL-Using-the-Repository-Pattern/models"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// HandleValidationError process validation errors and send appropriate responses.
func HandleValidationError(c *gin.Context, err error) {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		validationErrors := make([]models.ValidationErrorDetail, len(ve))
		for i, fieldError := range ve {
			validationErrors[i] = models.ValidationErrorDetail{
				Field:   fieldError.Field(),
				Message: msgForTag(fieldError.Tag()),
			}
		}
		SendValidationErrorResponse(c, validationErrors)
	} else {
		SendErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
	}
}

func msgForTag(tag string) string {
	switch tag {
	case "required":
		return "This field is required"
	}
	return ""
}
