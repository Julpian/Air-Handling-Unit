package repository

import "ahu-backend/internal/domain"

type AuditTrailRepository interface {
	Save(a *domain.AuditTrail) error

	ListAll(limit int) ([]domain.AuditTrailView, error)
	ListByEntity(entity, entityID string) ([]domain.AuditTrailView, error)
}
