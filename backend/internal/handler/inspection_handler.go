package handler

import (
	"net/http"

	"ahu-backend/internal/usecase"
	"ahu-backend/internal/usecase/dto"

	"github.com/gin-gonic/gin"
)

type scanNFCRequest struct {
	NFCUID      string `json:"nfc_uid"`
	ScheduleID  string `json:"schedule_id"`
	InspectorID string `json:"inspector_id"`
}

type InspectionHandler struct {
	inspectionUC   *usecase.InspectionUsecase
	queryUC        *usecase.InspectionQueryUsecase
	scanNFCUsecase *usecase.ScanNFCUsecase // 🔥 TAMBAH INI
}

func NewInspectionHandler(
	inspectionUC *usecase.InspectionUsecase,
	queryUC *usecase.InspectionQueryUsecase,
	scanNFCUsecase *usecase.ScanNFCUsecase, // 🔥 TAMBAH PARAMETER
) *InspectionHandler {
	return &InspectionHandler{
		inspectionUC:   inspectionUC,
		queryUC:        queryUC,
		scanNFCUsecase: scanNFCUsecase, // 🔥 SIMPAN
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
	var req dto.ScanNFCRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}

	inspectorID := c.GetString("user_id")

	// 🔥 DEBUG
	println("NFC:", req.NFCUID)
	println("INSPECTOR:", inspectorID)

	res, err := h.scanNFCUsecase.Execute(req, inspectorID)
	if err != nil {
		// 🔥 INI PENTING
		println("SCAN ERROR:", err.Error())

		c.JSON(400, gin.H{"message": err.Error()})
		return
	}

	c.JSON(200, res)
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
