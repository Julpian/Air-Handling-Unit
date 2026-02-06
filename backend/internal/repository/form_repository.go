package repository

import (
	"ahu-backend/internal/domain"
	"context"
)

type FormRepository interface {
	CreateTemplate(ctx context.Context, template *domain.FormTemplate) error

	GetTemplateBySchedule(scheduleID string) (*domain.FormTemplate, error)
	GetTemplateByID(templateID string) (*domain.FormTemplate, error)

	ListTemplates() ([]domain.FormTemplate, error)
	SetActive(id string, active bool) error

	CreateNewVersion(
		ctx context.Context,
		oldID string,
		template *domain.FormTemplate,
	) error

	ListTemplateVersions(templateID string) ([]domain.FormTemplate, error)
}

type sectionWithItems struct {
	Section domain.FormSection
}
