package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tgiv014/sugarcube/settings"
)

func (a *App) getSettings(c *gin.Context) {
	settings, err := a.settings.Get()
	if err != nil {
		Error(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, settings)
}

type settingsUpdate struct {
	DexcomUsername *string `json:"dexcomUsername"`
	DexcomPassword *string `json:"dexcomPassword"`
}

func (s settingsUpdate) apply(settings *settings.Settings) {
	if s.DexcomUsername != nil {
		settings.DexcomUsername = *s.DexcomUsername
	}
	if s.DexcomPassword != nil {
		settings.DexcomPassword = *s.DexcomPassword
	}
}

func (a *App) updateSettings(c *gin.Context) {
	update := settingsUpdate{}

	err := c.ShouldBindJSON(&update)
	if err != nil {
		Error(c, http.StatusBadRequest, err)
		return
	}

	settings, err := a.settings.Get()
	if err != nil {
		Error(c, http.StatusInternalServerError, err)
		return
	}

	update.apply(settings)

	err = a.settings.Save(settings)
	if err != nil {
		Error(c, http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}
