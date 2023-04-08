package app

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
	"github.com/tgiv014/sugarcube/session"
	"golang.org/x/crypto/bcrypt"
)

type loginRequest struct {
	Password string `json:"password"`
}

func (l loginRequest) validate() error {
	if len(l.Password) < 8 {
		return errors.New("password is too short")
	}

	return nil
}

func (a *App) login(c *gin.Context) {
	var req loginRequest
	err := c.Bind(&req)
	if err != nil {
		err = fmt.Errorf("could not bind request: %w", err)
		log.Warn("couldn't bind login request", "err", err)
		Error(c, http.StatusInternalServerError, err)
		return
	}

	newSession, err := a.sessions.Login(req.Password)
	if errors.Is(err, session.ErrIncorrectPassword) {
		Error(c, http.StatusUnauthorized, err)
		return
	}
	if err != nil {
		log.Warn("couldn't log in", "err", err)
		Error(c, http.StatusInternalServerError, err)
		return
	}

	c.SetCookie("token", newSession.Token, 24*60*60, "", "", false, true)
	c.Status(http.StatusOK)
}

func (a *App) signup(c *gin.Context) {
	var req loginRequest
	err := c.Bind(&req)
	if err != nil {
		log.Warn("couldn't bind signup request", "err", err)
		Error(c, http.StatusInternalServerError, err)
		return
	}

	err = req.validate()
	if err != nil {
		log.Warn("signup request invalid", "err", err)
		Error(c, http.StatusBadRequest, err)
		return
	}

	settings, err := a.settings.Get()
	if err != nil {
		log.Warn("couldn't get settings", "err", err)
		Error(c, http.StatusInternalServerError, err)
		return
	}
	if len(settings.HashedPassword) > 0 {
		err = errors.New("signup already completed")
		log.Warn("signup already completed")
		Error(c, http.StatusForbidden, err)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Warn("failed to generate bcrypt hash", "err", err)
		Error(c, http.StatusInternalServerError, err)
		return
	}

	settings.HashedPassword = hashedPassword

	err = a.settings.Save(settings)
	if err != nil {
		log.Warn("failed to save password to settings", "err", err)
		Error(c, http.StatusInternalServerError, err)
		return
	}

	newSession, err := a.sessions.Login(req.Password)
	if errors.Is(err, session.ErrIncorrectPassword) {
		Error(c, http.StatusUnauthorized, err)
		return
	}
	if err != nil {
		log.Warn("couldn't log in", "err", err)
		Error(c, http.StatusInternalServerError, err)
		return
	}

	c.SetCookie("token", newSession.Token, 24*60*60, "", "", false, true)

	c.Status(http.StatusOK)
}
