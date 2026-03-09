package handler

import (
	"ahu-backend/internal/domain"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ================= AHU =================

func (h *Handlers) ListAHUs(c *gin.Context) {
	list, err := h.AHUUC.ListAll()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, list)
}

type updateAHURequest struct {
	UnitCode         *string `json:"unit_code"`
	RoomName         *string `json:"room_name"`
	Vendor           *string `json:"vendor"`
	NFCUID           *string `json:"nfc_uid"`
	IsActive         *bool   `json:"is_active"`
	AreaID           *string `json:"area_id"`
	CleanlinessClass *string `json:"cleanliness_class"`
}

func (h *Handlers) UpdateAHU(c *gin.Context) {
	id := c.Param("id")

	var req updateAHURequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	ahu, err := h.AHUUC.GetByID(id)
	if err != nil {
		c.JSON(404, gin.H{"error": "unit code AHU tidak ditemukan"})
		return
	}

	if req.UnitCode != nil {
		ahu.UnitCode = *req.UnitCode
	}
	if req.Vendor != nil {
		ahu.Vendor = req.Vendor
	}
	if req.NFCUID != nil {
		ahu.NFCUID = req.NFCUID
	}
	if req.IsActive != nil {
		ahu.IsActive = *req.IsActive
	}
	if req.AreaID != nil {
		ahu.AreaID = *req.AreaID
	}

	if req.CleanlinessClass != nil {
		ahu.CleanlinessClass = req.CleanlinessClass
	}

	if req.RoomName != nil {
		ahu.RoomName = req.RoomName
	}

	if err := h.AHUUC.Update(ahu); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "unit code AHU berhasil diperbarui"})
}

func (h *Handlers) CreateAHU(c *gin.Context) {
	adminID := c.GetString("user_id")
	adminName := c.GetString("user_name")

	var req struct {
		BuildingID       string  `json:"building_id"`
		AreaID           string  `json:"area_id"`
		UnitCode         string  `json:"unit_code"`
		RoomName         *string `json:"room_name"`
		Vendor           *string `json:"vendor"`
		NFCUID           *string `json:"nfc_uid"`
		CleanlinessClass *string `json:"cleanliness_class"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "format request tidak valid"})
		return
	}

	ahu := &domain.AHU{
		ID:               uuid.NewString(),
		BuildingID:       req.BuildingID,
		AreaID:           req.AreaID,
		UnitCode:         req.UnitCode,
		RoomName:         req.RoomName,
		Vendor:           req.Vendor,
		NFCUID:           req.NFCUID,
		CleanlinessClass: req.CleanlinessClass,
		IsActive:         true,
		CreatedAt:        time.Now(),
	}

	if err := h.AHUUC.Create(ahu, adminID, adminName); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "unit code AHU berhasil dibuat",
		"id":      ahu.ID,
	})
}

func (h *Handlers) DeactivateAHU(c *gin.Context) {
	id := c.Param("id")
	adminID := c.GetString("user_id")

	if err := h.AHUUC.Deactivate(id, adminID); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "unit code AHU berhasil dinonaktifkan"})
}

func (h *Handlers) GetAHUDetail(c *gin.Context) {
	id := c.Param("id")

	ahu, err := h.AHUUC.GetByID(id)
	if err != nil {
		c.JSON(404, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, ahu)
}
