package app

import (
	"net/http"
	"time"

	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
)

func (a *App) getReadings(c *gin.Context) {
	startQ := c.Query("start")
	endQ := c.Query("end")

	// Defaults
	end := time.Now().UTC()
	start := end.Add(-time.Hour)

	if startQ != "" {
		parsedStartQuery, err := time.Parse(time.RFC3339, startQ)
		if err != nil {
			Error(c, http.StatusBadRequest, err)
			return
		}
		start = parsedStartQuery
	}
	if endQ != "" {
		parsedEndQuery, err := time.Parse(time.RFC3339, endQ)
		if err != nil {
			Error(c, http.StatusBadRequest, err)
			return
		}
		end = parsedEndQuery
	}
	readings, err := a.glucose.GetReadings(
		start,
		end,
	)
	if err != nil {
		log.Warn("failed to get readings", "err", err)
		Error(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, readings)
}
