package app

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type getReadingsRequest struct {
	Start *int64 `json:"start,omitempty"`
	End   *int64 `json:"end,omitempty"`
}

func (a *App) getReadings(c *gin.Context) {
	var request getReadingsRequest
	err := c.BindJSON(&request)
	if err != nil {
		Error(c, http.StatusInternalServerError, err)
		return
	}
	start := time.Now().Add(-time.Hour * 2).Unix()
	end := time.Now().Unix()

	if request.Start != nil {
		start = *request.Start
	}
	if request.End != nil {
		end = *request.End
	}
	readings, err := a.glucose.GetReadings(
		time.Unix(start, 0),
		time.Unix(end, 0),
	)
	if err != nil {
		Error(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, readings)
}
