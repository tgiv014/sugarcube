package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type statusResponse struct {
	SetupCompleted bool `json:"setupCompleted"`
}

func (a *App) status(c *gin.Context) {
	settings, err := a.settings.Get()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		SetupCompleted: settings.Completed(),
	})
}
