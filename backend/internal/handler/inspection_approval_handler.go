package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type approveRejectRequest struct {
	Note string `json:"note"`
}

func (h *Handlers) ApproveInspection(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString("user_id")

	if err := h.InspectionApprovalUC.ApproveInspection(
		id,
		userID,
	); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"status": "approved"})
}

func (h *Handlers) RejectInspection(c *gin.Context) {
	inspectionID := c.Param("id")
	supervisorID := c.GetString("user_id")

	var req approveRejectRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.Note == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "note is required"})
		return
	}

	err := h.InspectionApprovalUC.RejectInspection(
		inspectionID,
		supervisorID,
		req.Note,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "inspection rejected",
	})
}
