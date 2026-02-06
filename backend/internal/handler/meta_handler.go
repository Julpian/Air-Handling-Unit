package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handlers) GetStatusMeta(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"inspection_status": []gin.H{
			{"value": "sedang_diisi", "label": "Sedang Diisi"},
			{"value": "terkirim", "label": "Terkirim"},
			{"value": "revisi", "label": "Perlu Revisi"},
			{"value": "disetujui", "label": "Disetujui"},
		},
		"schedule_status": []gin.H{
			{"value": "siap_diperiksa", "label": "Siap Diperiksa"},
			{"value": "dalam_pemeriksaan", "label": "Dalam Pemeriksaan"},
			{"value": "selesai", "label": "Selesai"},
		},
	})
}
