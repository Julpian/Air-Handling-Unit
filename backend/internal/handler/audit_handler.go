package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ================= LIST ALL =================
func (h *Handlers) ListAuditTrails(c *gin.Context) {
	list, err := h.AuditUC.ListAll(100)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, list)
}

// ================= LIST BY ENTITY =================
func (h *Handlers) ListAuditByEntity(c *gin.Context) {
	entity := c.Param("entity")
	id := c.Param("id")

	list, err := h.AuditUC.ListByEntity(entity, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, list)
}
