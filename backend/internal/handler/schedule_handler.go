package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"ahu-backend/internal/usecase"
)

type ScheduleHandler struct {
	bypassUC *usecase.ScheduleBypassNFCUsecase
	assignUC *usecase.ScheduleAssignUsecase
	queryUC  *usecase.ScheduleQueryUsecase
}

func NewScheduleHandler(
	bypassUC *usecase.ScheduleBypassNFCUsecase,
	assignUC *usecase.ScheduleAssignUsecase,
	queryUC *usecase.ScheduleQueryUsecase, // ⬅️ TAMBAH
) *ScheduleHandler {
	return &ScheduleHandler{
		bypassUC: bypassUC,
		assignUC: assignUC,
		queryUC:  queryUC, // ⬅️ TAMBAH
	}
}

func (h *ScheduleHandler) BypassNFC(c *gin.Context) {
	scheduleID := c.Param("id")

	var req struct {
		Enabled bool   `json:"enabled"`
		Reason  string `json:"reason"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetString("user_id")

	if err := h.bypassUC.Execute(
		c.Request.Context(),
		scheduleID,
		req.Enabled,
		userID,
		req.Reason,
	); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "NFC bypass updated"})
}

type assignInspectorRequest struct {
	InspectorID string `json:"inspector_id"`
}

func (h *ScheduleHandler) AssignInspector(c *gin.Context) {
	scheduleID := c.Param("id")

	var req assignInspectorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	adminID := c.GetString("user_id")

	err := h.assignUC.AssignInspector(
		scheduleID,
		req.InspectorID,
		adminID,
	)

	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "inspector berhasil di-assign"})
}

func (h *ScheduleHandler) List(c *gin.Context) {
	data, err := h.queryUC.List()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, data)
}
