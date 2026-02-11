package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"ahu-backend/internal/usecase"
	"ahu-backend/internal/usecase/dto"
)

type ScanNFCHandler struct {
	uc *usecase.ScanNFCUsecase
}

func (h *ScanNFCHandler) Scan(c *gin.Context) {

	var req dto.ScanNFCRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// 🔥 ambil inspector dari JWT
	inspectorID := c.GetString("user_id")

	if inspectorID == "" {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}

	res, err := h.uc.Execute(req, inspectorID) // ⬅️ KIRIM inspectorID
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}
