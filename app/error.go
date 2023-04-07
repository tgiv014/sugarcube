package app

import (
	"github.com/gin-gonic/gin"
)

func Error(c *gin.Context, statusCode int, err error) {
	c.JSON(statusCode, gin.H{
		"error": err.Error(),
	})
}
