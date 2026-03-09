package usecase

import (
	"ahu-backend/internal/domain"
	"ahu-backend/internal/repository"

	"github.com/google/uuid"
)

type BuildingUsecase struct {
	repo repository.BuildingRepository
	auditUC *AuditTrailUsecase
}

func NewBuildingUsecase(repo repository.BuildingRepository, auditUC *AuditTrailUsecase) *BuildingUsecase {
	return &BuildingUsecase{
		repo: repo,
		auditUC: auditUC,
	}
}

func (u *BuildingUsecase) Create(name string, description *string, adminID string, adminName string) error {
    b := &domain.Building{
        ID:          uuid.NewString(),
        Name:        name,
        Description: description,
    }

    err := u.repo.Create(b, adminID)
    if err == nil {
        u.auditUC.Log(&domain.AuditTrail{
            UserID:   adminID,
            Action:   "CREATE_BUILDING",
            Entity:   "Building",
            EntityID: b.ID,
            Metadata: map[string]interface{}{
                "name":       b.Name,
                "admin_name": adminName,
            },
        })
    }
    return err
}

func (u *BuildingUsecase) List() ([]domain.Building, error) {
	return u.repo.ListAll()
}

func (u *BuildingUsecase) Update(
	id string,
	name string,
	description *string,
) error {

	b := &domain.Building{
		ID:          id,
		Name:        name,
		Description: description,
	}

	return u.repo.Update(b)
}

func (u *BuildingUsecase) Deactivate(id string) error {
	return u.repo.Deactivate(id)
}
