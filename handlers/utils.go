package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gitlab.com/amcop-saas-platform/vcs/vcs/billing-microservice/db"
	"net/http"
	"strings"
)

// HandleApiErrorResponse sends a formatted error response for validation errors
func HandleApiErrorResponse(c *gin.Context, err error) {
	var validationErrors validator.ValidationErrors
	var errorMessage string
	switch {
	case errors.Is(err, db.ErrRecordNotFound):
		c.JSON(http.StatusNotFound, gin.H{
			"error": "record not found",
		})
	case errors.Is(err, db.ErrDuplicatedRecord):
		c.JSON(http.StatusConflict, gin.H{
			"error": "record not found",
		})
	case errors.Is(err, db.ErrInternalServer):
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "unable to get record from database",
		})
	case errors.As(err, &validationErrors):
		var errorsList []string
		for _, fieldErr := range validationErrors {
			tag := fieldErr.ActualTag()
			var msg string
			if tag == "oneof" {
				msg = "invalid value for " + fieldErr.Field() + ". Expected values : " + fieldErr.Param()
			} else {
				msg = fieldErr.Field() + " is " + fieldErr.Tag()
			}
			errorsList = append(errorsList, msg)
		}
		errorMessage = strings.Join(errorsList, ", ")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errorMessage,
		})
	default:
		if err.Error() == "EOF" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "missing required fields",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}

	}
}
