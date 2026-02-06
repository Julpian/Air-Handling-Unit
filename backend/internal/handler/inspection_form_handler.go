package handler

import (
	"github.com/gin-gonic/gin"

	"ahu-backend/internal/usecase"
	"ahu-backend/internal/usecase/dto"
)

type InspectionFormHandler struct {
	getFormUC *usecase.GetFormByInspectionUsecase
	submitUC  *usecase.SubmitInspectionFormUsecase
}

func NewInspectionFormHandler(
	getFormUC *usecase.GetFormByInspectionUsecase,
	submitUC *usecase.SubmitInspectionFormUsecase,
) *InspectionFormHandler {
	return &InspectionFormHandler{getFormUC: getFormUC, submitUC: submitUC}
}

func (h *InspectionFormHandler) Submit(
	c *gin.Context,
) {
	inspectionID := c.Param("inspection_id")

	var req dto.SubmitInspectionFormRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := h.submitUC.Execute(inspectionID, req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "inspection submitted"})
}

func (h *InspectionFormHandler) GetForm(
	c *gin.Context,
) {
	inspectionID := c.Param("inspection_id")

	form, err := h.getFormUC.Execute(inspectionID)
	if err != nil {
		c.JSON(404, gin.H{"error": "form not found"})
		return
	}

	c.JSON(200, form)
}
