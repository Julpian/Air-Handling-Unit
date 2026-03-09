package usecase

import (
	"errors"

	"ahu-backend/internal/domain"
	"ahu-backend/internal/repository"
)

type AHUUsecase struct {
	repo repository.AHURepository
	auditUC *AuditTrailUsecase
}

func NewAHUUsecase(r repository.AHURepository, audit *AuditTrailUsecase) *AHUUsecase {
    return &AHUUsecase{repo: r, auditUC: audit}
}

func (u *AHUUsecase) Create(ahu *domain.AHU, adminID string, adminName string) error {
	if ahu.UnitCode == "" {
		return errors.New("unit code AHU wajib diisi") // ✅ konsisten
	}

	if ahu.NFCUID != nil && *ahu.NFCUID == "" {
		ahu.NFCUID = nil
	}

	if ahu.RoomName != nil && *ahu.RoomName == "" {
		ahu.RoomName = nil
	}

	if ahu.CleanlinessClass != nil && !isValidCleanClass(*ahu.CleanlinessClass) {
		return errors.New("kelas kebersihan harus E, F, atau G")
	}

	err := u.repo.Create(ahu)
    if err == nil {
        u.auditUC.Log(&domain.AuditTrail{
            UserID:   adminID,
            Action:   "CREATE_AHU",
            Entity:   "AHU",
            EntityID: ahu.ID,
            Metadata: map[string]interface{}{
                "unit_code": ahu.UnitCode,
                "name":      adminName,
            },
        })
    }
    return err
}

func (u *AHUUsecase) ListAll() ([]domain.AHU, error) {
	return u.repo.ListAll()
}

func (u *AHUUsecase) GetByID(id string) (*domain.AHU, error) {
	return u.repo.GetByID(id)
}

func (u *AHUUsecase) Update(ahu *domain.AHU) error {
	if ahu.UnitCode == "" {
		return errors.New("nama AHU wajib diisi")
	}

	if ahu.RoomName != nil && *ahu.RoomName == "" {
		ahu.RoomName = nil
	}

	if ahu.CleanlinessClass != nil && !isValidCleanClass(*ahu.CleanlinessClass) {
		return errors.New("kelas kebersihan harus E, F, atau G")
	}

	return u.repo.Update(ahu)
}

func (u *AHUUsecase) Deactivate(id string, adminID string) error {
	return u.repo.Deactivate(id)
}

func isValidCleanClass(v string) bool {
	return v == "E" || v == "F" || v == "G"
}
