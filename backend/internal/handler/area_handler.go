package handler

import (
	"github.com/gin-gonic/gin"

	"ahu-backend/internal/domain"
	"ahu-backend/internal/usecase"
)

type AreaHandler struct {
	uc *usecase.AreaUsecase
}

func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func NewAreaHandler(uc *usecase.AreaUsecase) *AreaHandler {
	return &AreaHandler{uc: uc}
}

func (h *AreaHandler) Create(c *gin.Context) {
	adminID := c.GetString("user_id")
	adminName := c.GetString("user_name")

	var req struct {
		BuildingID       string `json:"building_id"`
		Name             string `json:"name"`
		Floor            string `json:"floor"`
		CleanlinessClass string `json:"cleanliness_class"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	area := &domain.Area{
		BuildingID:       req.BuildingID,
		Name:             req.Name,
		Floor:            strPtr(req.Floor),
		CleanlinessClass: strPtr(req.CleanlinessClass),
	}

	if err := h.uc.Create(area, adminID, adminName); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, gin.H{"message": "area berhasil dibuat"})
}

func (h *AreaHandler) List(c *gin.Context) {
	list, err := h.uc.ListAll()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, list)
}

func (h *AreaHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		BuildingID       string `json:"building_id"`
		Name             string `json:"name"`
		Floor            string `json:"floor"`
		CleanlinessClass string `json:"cleanliness_class"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	area := &domain.Area{
		ID:               id,
		BuildingID:       req.BuildingID,
		Name:             req.Name,
		Floor:            strPtr(req.Floor),
		CleanlinessClass: strPtr(req.CleanlinessClass),
	}

	if err := h.uc.Update(area); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "area diperbarui"})
}

func (h *AreaHandler) Deactivate(c *gin.Context) {
	id := c.Param("id")

	if err := h.uc.Deactivate(id); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "area dinonaktifkan"})
}
