package httputil

import (
	"github.com/gin-gonic/gin"
)

type ResponseError struct {
	Status int           `json:"status"`
	Error  []interface{} `json:"error"`
}

func NewResponseError(c *gin.Context, statusCode int, err interface{}) {
	response := &ResponseError{
		Status: statusCode,
		Error: []interface{}{
			err,
		},
	}

	c.JSON(statusCode, response)
}

func NewResponseAbort(c *gin.Context, statusCode int, err interface{}) {
	response := &ResponseError{
		Status: statusCode,
		Error: []interface{}{
			err,
		},
	}

	c.AbortWithStatusJSON(statusCode, response)
}

type ResponseSuccess struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func NewResponseSuccess(c *gin.Context, statusCode int, message string) {
	response := &ResponseSuccess{
		Status:  statusCode,
		Message: message,
	}

	c.JSON(statusCode, response)
}
