package usecase

import (
	"errors"

	"ahu-backend/internal/domain"
	"ahu-backend/internal/repository"
)

type AreaUsecase struct {
	repo repository.AreaRepository
	auditUC *AuditTrailUsecase
}

func NewAreaUsecase(repo repository.AreaRepository, auditUC *AuditTrailUsecase) *AreaUsecase {
	return &AreaUsecase{
		repo: repo,
		auditUC: auditUC,
	}
}

func (u *AreaUsecase) Create(area *domain.Area, adminID string, adminName string) error {
    if area.Name == "" || area.BuildingID == "" {
        return errors.New("nama area dan gedung wajib diisi")
    }
    
    err := u.repo.Create(area)
    if err == nil {
        // PERBAIKAN: Gunakan struct domain.AuditTrail sesuai permintaan fungsi Log
        u.auditUC.Log(&domain.AuditTrail{
            UserID:   adminID,
            Action:   "CREATE_AREA",
            Entity:   "Area",
            EntityID: area.ID,
            Metadata: map[string]interface{}{
                "name":       area.Name,
                "admin_name": adminName,
            },
        })
    }
    return err
}

func (u *AreaUsecase) ListAll() ([]domain.Area, error) {
	return u.repo.ListAll()
}

func (u *AreaUsecase) Update(area *domain.Area) error {
	return u.repo.Update(area)
}

func (u *AreaUsecase) Deactivate(id string) error {
	return u.repo.Deactivate(id)
}
