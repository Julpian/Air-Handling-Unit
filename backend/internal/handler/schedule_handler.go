package handler

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"ahu-backend/internal/usecase"
)

func hashFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}

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
		c.String(404, "Dokumen tidak valid")
		return
	}

	// pastikan PDFPath & PDFHash ada
	if approval.PDFPath == nil || approval.PDFHash == nil {
		c.String(500, "Hash PDF belum tersedia")
		return
	}

	// ===== HITUNG ULANG HASH PDF =====
	pdfHash, err := hashFile(*approval.PDFPath)
	if err != nil {
		c.String(500, "Gagal membaca PDF")
		return
	}

	validIntegrity := pdfHash == *approval.PDFHash

	c.HTML(200, "verify.html", gin.H{
		"year":       approval.Year,
		"svp":        approval.SVPID,
		"asmen":      approval.AsmenID,
		"svp_time":   approval.SVPSignedAt,
		"asmen_time": approval.AsmenSignedAt,
		"pdf":        *approval.PDFPath,

		"integrity": validIntegrity,
		"hash":      pdfHash,
	})
}

func (h *ScheduleHandler) ListByYear(c *gin.Context) {
	year, _ := strconv.Atoi(c.Param("year"))

	data, err := h.queryUC.ListByYear(year)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, data)
}
