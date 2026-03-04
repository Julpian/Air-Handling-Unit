package handler

import (
	"fmt"
	"net/http"
	"strings"

	"ahu-backend/internal/domain"
	"ahu-backend/internal/middleware"
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
	taskUC         *usecase.InspectionTaskUsecase
	signUsecase    *usecase.SignInspectionUsecase

	approveUC *usecase.ApproveInspectionUsecase // 🔥 TAMBAH
	pdf       *usecase.InspectionPDFService
}

func NewInspectionHandler(
	inspectionUC *usecase.InspectionUsecase,
	queryUC *usecase.InspectionQueryUsecase,
	scanNFCUsecase *usecase.ScanNFCUsecase,
	taskUC *usecase.InspectionTaskUsecase,
	signUsecase *usecase.SignInspectionUsecase,
	approveUC *usecase.ApproveInspectionUsecase,
	pdf *usecase.InspectionPDFService,
) *InspectionHandler {
	return &InspectionHandler{
		inspectionUC:   inspectionUC,
		queryUC:        queryUC,
		scanNFCUsecase: scanNFCUsecase,
		taskUC:         taskUC,
		signUsecase:    signUsecase,

		approveUC: approveUC, // 🔥
		pdf:       pdf,       // 🔥
	}
}

// ================= LIST =================
func (h *InspectionHandler) List(c *gin.Context) {
	status := c.Query("status")

	// 🔥 PERBAIKAN: Gunakan middleware.GetUser agar Role tidak kosong
	user := middleware.GetUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	role := strings.ToLower(user.Role)
	userID := user.ID

	// DEBUG untuk memastikan role sudah terisi (cek di terminal nanti)
	fmt.Printf("DEBUG: Request List - Role: %s, UserID: %s, Status: %s\n", role, userID, status)

	var list []domain.Inspection
	var err error

	if role == "admin" || role == "spv" {
		// Admin melihat semua laporan tanpa filter inspector_id
		list, err = h.queryUC.ListByStatus(status, "")
	} else {
		// Inspector hanya melihat miliknya sendiri
		list, err = h.queryUC.ListByStatus(status, userID)
	}

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, list)
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

func (h *InspectionHandler) Tasks(c *gin.Context) {
	inspectorID := c.GetString("user_id")

	list, err := h.taskUC.ListByInspector(inspectorID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// 🔥 DEBUG
	println("TASK COUNT:", len(list))

	c.JSON(200, list)
}

func (h *InspectionHandler) SignInspection(c *gin.Context) {
	id := c.Param("inspection_id")

	var req dto.SignInspectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := h.signUsecase.Execute(id, req.Signature); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"success": true})
}

func (h *InspectionHandler) Approve(c *gin.Context) {
	id := c.Param("inspection_id")

	// 🔥 FIX 1: Definisikan variabel 'req'
	var req struct {
		Signature string `json:"signature"`
	}

	// Bind JSON dari frontend ke variabel req
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Signature is required"})
		return
	}

	// 🔥 FIX 2: Definisikan variabel 'user' menggunakan middleware
	user := middleware.GetUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// 1. Jalankan Usecase (Hanya sekali)
	if err := h.approveUC.Execute(id, user.ID, req.Signature); err != nil {
		fmt.Println("DEBUG: Usecase Error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 2. Generate PDF
	if err := h.pdf.GenerateInspectionPDF(id); err != nil {
		fmt.Println("Error generating PDF:", err)
		c.JSON(http.StatusOK, gin.H{"ok": true, "warning": "Approved, but PDF generation failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func (h *InspectionHandler) VerifyInspection(c *gin.Context) {
	id := c.Param("id")
	fmt.Println("DEBUG: Hit VerifyInspection for ID:", id)

	// Panggil repository/usecase untuk ambil data verifikasi
	data, err := h.queryUC.GetVerificationData(id)
	if err != nil {
		fmt.Println("DEBUG: DB Error or Not Found:", err)
		c.JSON(404, gin.H{"error": "Dokumen tidak ditemukan atau belum di-approve"})
		return
	}

	// Format Periode agar rapi (sama dengan logika di PDF)
	displayPeriod := data.Period
	p := strings.ToLower(data.Period)
	if strings.Contains(p, "bulan") || strings.Contains(p, "month") {
		displayPeriod = "Monthly (1 Month)"
	} else if strings.Contains(p, "6") {
		displayPeriod = "6 Months"
	} else if strings.Contains(p, "tahun") || strings.Contains(p, "year") {
		displayPeriod = "Yearly (1 Year)"
	}

	c.JSON(200, gin.H{
		"inspection_id":  data.InspectionID,
		"unit_code":      data.UnitCode,
		"period":         displayPeriod,
		"inspector_name": data.Inspector,
		"spv_name":       data.SPVName,
		"inspected_at":   data.InspectedAt.Format("02 January 2006"),
	})
}
