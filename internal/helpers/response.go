package helpers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ResponseSuccess untuk response sukses dengan data
func ResponseSuccess(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": message,
		"data":    data,
	})
}

// ResponseCreated untuk response data berhasil dibuat
func ResponseCreated(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": message,
		"data":    data,
	})
}

// ResponseError untuk response error
func ResponseError(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, gin.H{
		"success": false,
		"error":   message,
	})
}

// ResponseBadRequest untuk response validation error
func ResponseBadRequest(c *gin.Context, message string) {
	ResponseError(c, http.StatusBadRequest, message)
}

// ResponseUnauthorized untuk response unauthorized
func ResponseUnauthorized(c *gin.Context, message string) {
	ResponseError(c, http.StatusUnauthorized, message)
}

// ResponseNotFound untuk response data tidak ditemukan
func ResponseNotFound(c *gin.Context, message string) {
	ResponseError(c, http.StatusNotFound, message)
}

// ResponseInternalServerError untuk response server error
func ResponseInternalServerError(c *gin.Context, message string) {
	ResponseError(c, http.StatusInternalServerError, message)
}
