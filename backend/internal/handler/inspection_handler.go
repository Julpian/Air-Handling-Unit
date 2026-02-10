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
	var req struct {
		NFCUID string `json:"nfc_uid"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.NFCUID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "nfc_uid wajib"})
		return
	}

	// ambil inspector dari JWT
	inspectorID := c.GetString("user_id")

	inspection, err := h.inspectionUC.ScanNFC(
		req.NFCUID,
		inspectorID,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// VALIDASI inspector
	if inspection.InspectorID != "" && inspection.InspectorID != inspectorID {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "jadwal ini bukan milik anda",
		})
		return
	}

	// 🔥 RESPONSE MINIMAL UNTUK FRONTEND
	c.JSON(http.StatusOK, gin.H{
		"id": inspection.ID,
	})
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
