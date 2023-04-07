package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type statusResponse struct {
	SetupCompleted bool `json:"setupCompleted"`
	SessionValid   bool `json:"sessionValid"`
}

func (a *App) status(c *gin.Context) {
	settings, err := a.settings.Get()
	if err != nil {
		Error(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		SetupCompleted: settings.Completed(),
		SessionValid:   a.validTokenCookie(c),
	})
}

func (a *App) validTokenCookie(c *gin.Context) bool {
	token, err := c.Cookie("token")
	if err != nil {
		return false
	}

	// If a session is returned, the token is valid
	_, err = a.sessions.GetSession(token)
	if err != nil {
		return false
	}

	return true
}
