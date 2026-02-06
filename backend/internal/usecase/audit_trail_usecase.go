package usecase

import (
	"ahu-backend/internal/domain"
	"ahu-backend/internal/repository"
)

type AuditTrailUsecase struct {
	repo repository.AuditTrailRepository
}

func NewAuditTrailUsecase(
	repo repository.AuditTrailRepository,
) *AuditTrailUsecase {
	return &AuditTrailUsecase{repo: repo}
}

// dipanggil dari bisnis logic (scan, approve, dll)
func (u *AuditTrailUsecase) Log(a *domain.AuditTrail) {
	_ = u.repo.Save(a)
}

// admin list
func (u *AuditTrailUsecase) ListAll(
	limit int,
) ([]domain.AuditTrailView, error) {
	return u.repo.ListAll(limit)
}

// detail entity
func (u *AuditTrailUsecase) ListByEntity(
	entity, entityID string,
) ([]domain.AuditTrailView, error) {
	return u.repo.ListByEntity(entity, entityID)
}
