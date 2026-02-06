package handler

import (
	"net/http"

	"ahu-backend/internal/usecase"

	"github.com/gin-gonic/gin"
)

type scanNFCRequest struct {
	NFCUID      string `json:"nfc_uid"`
	ScheduleID  string `json:"schedule_id"`
	InspectorID string `json:"inspector_id"`
}

type InspectionHandler struct {
	inspectionUC *usecase.InspectionUsecase
	queryUC      *usecase.InspectionQueryUsecase
}

func NewInspectionHandler(
	inspectionUC *usecase.InspectionUsecase,
	queryUC *usecase.InspectionQueryUsecase,
) *InspectionHandler {
	return &InspectionHandler{
		inspectionUC: inspectionUC,
		queryUC:      queryUC,
	}
}

// ================= LIST =================
func (h *InspectionHandler) List(c *gin.Context) {
	status := c.Query("status")

	list, err := h.queryUC.ListByStatus(status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, list)
}

// ================= DETAIL =================
func (h *InspectionHandler) Detail(c *gin.Context) {
	id := c.Param("inspection_id")

	inspection, err := h.inspectionUC.GetByID(id)
	if err != nil || inspection == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "inspection tidak ditemukan",
		})
		return
	}

	c.JSON(http.StatusOK, inspection)
}

// ================= SCAN NFC =================
func (h *InspectionHandler) ScanNFC(c *gin.Context) {
	var req scanNFCRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	inspection, err := h.inspectionUC.ScanNFC(
		req.NFCUID,
		req.ScheduleID,
		req.InspectorID,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, inspection)
}

// ================= DROPDOWN (ADMIN) =================
func (h *InspectionHandler) ListDropdown(c *gin.Context) {
	data, err := h.queryUC.ListDropdown()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, data)
}

func (h *InspectionHandler) Dashboard(c *gin.Context) {
	userID := c.GetString("user_id")

	// sementara response sederhana dulu
	// nanti bisa dikembangkan (total inspeksi, jadwal hari ini, dll)
	c.JSON(http.StatusOK, gin.H{
		"message": "Inspector dashboard",
		"user_id": userID,
	})
}
