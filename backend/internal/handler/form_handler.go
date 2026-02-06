package handler

import (
	"ahu-backend/internal/domain"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateFormTemplateRequest struct {
	Name        string `json:"name" binding:"required"`
	Period      string `json:"period" binding:"required"`
	Description string `json:"description"`

	Sections []struct {
		Code  string `json:"code" binding:"required"`
		Title string `json:"title" binding:"required"`

		Items []struct {
			Label    string `json:"label" binding:"required"`
			Type     string `json:"type" binding:"required"`
			Required bool   `json:"required"`
		} `json:"items" binding:"required"`
	} `json:"sections" binding:"required"`
}

func (h *Handlers) CreateFormTemplate(c *gin.Context) {
	var req CreateFormTemplateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	template := domain.FormTemplate{
		Name:        req.Name,
		Period:      req.Period,
		Description: strPtr(req.Description),
		IsActive:    true,
	}

	for secIndex, s := range req.Sections {
		section := domain.FormSection{
			Code:  s.Code,
			Title: s.Title,
			Order: secIndex + 1,
		}

		for itemIndex, i := range s.Items {
			section.Items = append(section.Items, domain.FormItem{
				ID:        uuid.NewString(),
				Label:     i.Label,
				InputType: i.Type,
				Required:  i.Required,
				Order:     itemIndex + 1,
			})
		}

		template.Sections = append(template.Sections, section)
	}

	if err := h.createFormTemplateUsecase.Execute(
		c.Request.Context(),
		&template,
	); err != nil {
		log.Println("[CREATE FORM TEMPLATE ERROR]:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Form template created",
		"id":      template.ID,
	})
}

func (h *Handlers) GetFormTemplateDetail(c *gin.Context) {
	templateID := c.Param("id")

	if templateID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "template id wajib diisi",
		})
		return
	}

	form, err := h.getFormTemplateDetailUsecase.Execute(templateID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": err.Error(),
		})
		return
	}

	if form == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "form template tidak ditemukan",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": form,
	})
}

func (h *Handlers) ListFormTemplates(c *gin.Context) {
	templates, err := h.listFormTemplateUsecase.Execute()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": templates,
	})
}

func (h *Handlers) SetFormTemplateActive(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Active bool `json:"active"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	if err := h.setFormTemplateActiveUsecase.Execute(id, req.Active); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "status updated",
	})
}

func (h *Handlers) CreateNewFormTemplateVersion(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Name        string  `json:"name"`
		Period      string  `json:"period"`
		Description *string `json:"description"`

		Sections []struct {
			Code  string `json:"code"`
			Title string `json:"title"`

			Items []struct {
				Label     string `json:"label"`
				InputType string `json:"input_type"`
				Required  bool   `json:"required"`
			} `json:"items"`
		} `json:"sections"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}

	template := domain.FormTemplate{
		Name:        req.Name,
		Period:      req.Period,
		Description: req.Description,
	}

	for si, s := range req.Sections {
		sec := domain.FormSection{
			Code:  s.Code,
			Title: s.Title,
			Order: si + 1,
		}

		for ii, i := range s.Items {
			sec.Items = append(sec.Items, domain.FormItem{
				Label:     i.Label,
				InputType: i.InputType,
				Required:  i.Required,
				Order:     ii + 1,
			})
		}

		template.Sections = append(template.Sections, sec)
	}

	if err := h.createNewFormTemplateVersionUsecase.Execute(
		c.Request.Context(),
		id,
		&template,
	); err != nil {
		log.Println("CREATE VERSION ERROR:", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "version created"})
}

func (h *Handlers) ListFormTemplateVersions(c *gin.Context) {
	id := c.Param("id")

	versions, err := h.listFormTemplateVersionsUsecase.Execute(id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"data": versions,
	})
}

func (h *Handlers) CompareFormTemplate(c *gin.Context) {
	fromID := c.Param("fromId")
	toID := c.Param("toId")

	log.Println("COMPARE:", fromID, toID) // 🔥 DEBUG

	diff, err := h.compareFormTemplateUsecase.Execute(fromID, toID)
	if err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"data": diff,
	})
}
