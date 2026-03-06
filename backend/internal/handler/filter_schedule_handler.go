package handler

import (
	"time"

	"github.com/gin-gonic/gin"
)

type CreateFilterScheduleRequest struct {
	AHUID     string `json:"ahu_id"`
	StartDate string `json:"start_date"`
}

func (h *Handlers) CreateFilterSchedule(c *gin.Context) {

	var req CreateFilterScheduleRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "payload tidak valid"})
		return
	}

	println("DATE RECEIVED:", req.StartDate)

	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		c.JSON(400, gin.H{"error": "format tanggal salah"})
		return
	}

	err = h.CreateFilterScheduleUC.Execute(
		req.AHUID,
		startDate,
	)

	if err != nil {
		println("FILTER ERROR:", err.Error())

		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	println("AHU:", req.AHUID)
	println("DATE:", req.StartDate)

	c.JSON(201, gin.H{
		"message": "schedule ganti filter berhasil dibuat",
	})
}