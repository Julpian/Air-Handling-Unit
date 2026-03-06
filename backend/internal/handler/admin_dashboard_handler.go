package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handlers) Dashboard(c *gin.Context) {

	data, err := h.dashboardUsecase.GetStats(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, data)
}

func (h *Handlers) FilterPressureChart(c *gin.Context) {

	data, err := h.dashboardUsecase.GetFilterPressureChart()
	if err != nil {

		println("FILTER CHART ERROR:", err.Error())

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, data)
}