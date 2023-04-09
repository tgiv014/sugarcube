package app

import (
	"net/http"
	"time"

	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
)

type getReadingsRequest struct {
	Start *int64 `json:"start,omitempty"`
	End   *int64 `json:"end,omitempty"`
}

func (a *App) getReadings(c *gin.Context) {
	var request getReadingsRequest
	err := c.Bind(&request)
	if err != nil {
		log.Warn("failed to unmarshal request", "err", err)
		Error(c, http.StatusInternalServerError, err)
		return
	}
	readings, err := a.glucose.GetReadings(
		time.Now().Add(-time.Hour*4),
		time.Now(),
	)
	if err != nil {
		log.Warn("failed to get readings", "err", err)
		Error(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, readings)
}
