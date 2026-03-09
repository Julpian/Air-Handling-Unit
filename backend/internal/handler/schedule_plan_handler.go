package handler

import (
	"ahu-backend/internal/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SchedulePlanHandler struct {
	SchedulePlanUC *usecase.SchedulePlanUsecase
}

func NewSchedulePlanHandler(
	uc *usecase.SchedulePlanUsecase,
) *SchedulePlanHandler {
	return &SchedulePlanHandler{
		SchedulePlanUC: uc,
	}
}

type createSchedulePlanRequest struct {
	AHUID       string `json:"ahu_id"`
	Period      string `json:"period"`
	WeekOfMonth int    `json:"week_of_month"`
	Month       *int   `json:"month"`
}

// ================= CREATE =================
func (h *SchedulePlanHandler) Create(c *gin.Context) {
	var req createSchedulePlanRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "payload tidak valid"})
		return
	}

	if req.AHUID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ahu wajib diisi"})
		return
	}

	if req.WeekOfMonth < 1 || req.WeekOfMonth > 4 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "minggu tidak valid"})
		return
	}

	if req.Period != "bulanan" && req.Month == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bulan wajib diisi"})
		return
	}

	validPeriods := map[string]bool{
		"bulanan":    true,
		"enam_bulan": true,
		"tahunan":    true,
	}

	if !validPeriods[req.Period] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "period tidak valid"})
		return
	}

	// 🔥 ambil user dari JWT middleware
	adminID := c.GetString("user_id")
	adminName := c.GetString("user_name")

	if err := h.SchedulePlanUC.Create(
		req.AHUID,
		req.Period,
		req.WeekOfMonth,
		req.Month,
		adminID,
		adminName,
	); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "schedule plan berhasil dibuat",
	})
}

// ================= LIST =================
func (h *SchedulePlanHandler) List(c *gin.Context) {
	list, err := h.SchedulePlanUC.ListAllWithAHU()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, list)
}

func (h *Handlers) UpdateSchedulePlan(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Period      string `json:"period"`
		WeekOfMonth int    `json:"week_of_month"`
		Month       *int   `json:"month"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "format request tidak valid"})
		return
	}

	adminID := c.GetString("user_id")
	adminName := c.GetString("user_name")

	if err := h.SchedulePlanUC.Update(
		id,
		req.Period,
		req.WeekOfMonth,
		req.Month,
		adminID,
		adminName,
	); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "schedule plan berhasil diupdate"})
}

func (h *Handlers) DeleteSchedulePlan(c *gin.Context) {
	id := c.Param("id")

	adminID := c.GetString("user_id")
	adminName := c.GetString("user_name")

	if err := h.SchedulePlanUC.Delete(
		id,
		adminID,
		adminName,
	); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "schedule plan berhasil dihapus"})
}
