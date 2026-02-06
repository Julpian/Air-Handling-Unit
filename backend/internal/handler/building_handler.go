package handler

import (
	"ahu-backend/internal/usecase"

	"github.com/gin-gonic/gin"
)

type BuildingHandler struct {
	uc *usecase.BuildingUsecase
}

func NewBuildingHandler(uc *usecase.BuildingUsecase) *BuildingHandler {
	return &BuildingHandler{uc: uc}
}

func (h *BuildingHandler) Create(c *gin.Context) {
	var req struct {
		Name        string  `json:"name"`
		Description *string `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetString("user_id")

	if err := h.uc.Create(req.Name, req.Description, userID); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, gin.H{"message": "building created"})
}

func (h *BuildingHandler) List(c *gin.Context) {
	data, err := h.uc.List()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, data)
}

func (h *BuildingHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Name        string  `json:"name"`
		Description *string `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if req.Name == "" {
		c.JSON(400, gin.H{"error": "name wajib"})
		return
	}

	if err := h.uc.Update(id, req.Name, req.Description); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "gedung berhasil diupdate"})
}

func (h *BuildingHandler) Deactivate(c *gin.Context) {
	id := c.Param("id")

	if err := h.uc.Deactivate(id); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "gedung dinonaktifkan"})
}
