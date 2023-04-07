package app

import (
	"errors"
	"net/http"

	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
	"github.com/tgiv014/sugarcube/session"
	"golang.org/x/crypto/bcrypt"
)

type loginRequest struct {
	Password string `json:"password"`
}

func (l loginRequest) valid() bool {
	if len(l.Password) < 8 {
		return false
	}

	return true
}

func (a *App) login(c *gin.Context) {
	var req loginRequest
	err := c.BindJSON(&req)
	if err != nil {
		log.Warn("couldn't bind login request", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	newSession, err := a.sessions.Login(req.Password)
	if errors.Is(err, session.ErrIncorrectPassword) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}
	if err != nil {
		log.Warn("couldn't log in", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.SetCookie("token", newSession.Token, 24*60*60, "", "", false, true)
	c.Status(http.StatusOK)
}

func (a *App) signup(c *gin.Context) {
	var req loginRequest
	err := c.BindJSON(&req)
	if err != nil {
		log.Warn("couldn't bind signup request", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if !req.valid() {
		log.Warn("signup request invalid")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid password",
		})
		return
	}

	settings, err := a.settings.Get()
	if err != nil {
		log.Warn("couldn't get settings", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	if len(settings.HashedPassword) > 0 {
		log.Warn("signup already completed")
		c.Status(http.StatusForbidden)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Warn("failed to generate bcrypt hash", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	settings.HashedPassword = hashedPassword

	err = a.settings.Save(settings)
	if err != nil {
		log.Warn("failed to save password to settings", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.Status(http.StatusOK)
}
