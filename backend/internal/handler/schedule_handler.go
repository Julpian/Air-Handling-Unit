package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"ahu-backend/internal/usecase"
)

type ScheduleHandler struct {
	bypassUC   *usecase.ScheduleBypassNFCUsecase
	assignUC   *usecase.ScheduleAssignUsecase
	queryUC    *usecase.ScheduleQueryUsecase
	approvalUC *usecase.ScheduleApprovalUsecase
}

func NewScheduleHandler(
	bypassUC *usecase.ScheduleBypassNFCUsecase,
	assignUC *usecase.ScheduleAssignUsecase,
	queryUC *usecase.ScheduleQueryUsecase,
	approvalUC *usecase.ScheduleApprovalUsecase, // ✅ TAMBAH
) *ScheduleHandler {
	return &ScheduleHandler{
		bypassUC:   bypassUC,
		assignUC:   assignUC,
		queryUC:    queryUC,
		approvalUC: approvalUC, // ✅ TAMBAH
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

	year := time.Now().Year()

	approval, _ := h.approvalUC.Get(year)

	if approval != nil && approval.Status == "completed" {
		c.JSON(400, gin.H{"error": "schedule locked"})
		return
	}

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

func (h *ScheduleHandler) Verify(c *gin.Context) {

	token := c.Param("token")

	approval, err := h.approvalUC.GetByToken(token)
	if err != nil || approval == nil {
		c.JSON(404, gin.H{
			"valid":   false,
			"message": "Dokumen tidak valid",
		})
		return
	}

	c.JSON(200, gin.H{
		"valid":    true,
		"year":     approval.Year,
		"status":   approval.Status,
		"pdf_path": approval.PDFPath,
	})
}
