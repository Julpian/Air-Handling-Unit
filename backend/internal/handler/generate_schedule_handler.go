package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type GenerateScheduleRequest struct {
	Year int `json:"year"`
}

func (h *Handlers) GenerateSchedule(c *gin.Context) {
	year := time.Now().Year()

	if err := h.GenerateUC.Generate(year); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Jadwal tahun berhasil di-generate",
		"year":    year,
	})
}
