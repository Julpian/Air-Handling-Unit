package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handlers) SVPSignSchedule(c *gin.Context) {

	year, _ := strconv.Atoi(c.Param("year"))
	userID := c.GetString("user_id")

	var req struct {
		Signature string `json:"signature"`
	}

	c.ShouldBindJSON(&req)

	err := h.ScheduleApprovalUC.SignSVP(year, userID, req.Signature)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "SVP signed"})
}

func (h *Handlers) AsmenSignSchedule(c *gin.Context) {
	year, _ := strconv.Atoi(c.Param("year"))

	userID := c.GetString("user_id") // karena kamu set ini di auth middleware

	var body struct {
		Signature string `json:"signature"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": "invalid body"})
		return
	}

	err := h.ScheduleApprovalUC.SignAsmen(year, userID, body.Signature)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"status": "completed"})
}

func (h *Handlers) GetScheduleApproval(c *gin.Context) {
	year, _ := strconv.Atoi(c.Param("year"))

	data, err := h.ScheduleApprovalUC.Get(year)
	if err != nil {
		c.JSON(404, gin.H{"error": "not found"})
		return
	}

	c.JSON(200, data)
}

func (h *Handlers) DownloadSchedulePDF(c *gin.Context) {
	year, _ := strconv.Atoi(c.Param("year"))

	s, err := h.ScheduleApprovalUC.Get(year)
	if err != nil {
		c.JSON(404, gin.H{"error": "not found"})
		return
	}

	if s.PDFPath == nil {
		c.JSON(404, gin.H{"error": "pdf not generated yet"})
		return
	}

	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", "attachment; filename=schedule-2026.pdf")
	c.File(*s.PDFPath)
}

func (h *Handlers) VerifySchedule(c *gin.Context) {
	token := c.Param("token")

	data, err := h.ScheduleApprovalUC.GetByToken(token)
	if err != nil {
		c.JSON(404, gin.H{"error": "invalid"})
		return
	}

	c.JSON(200, data)
}
