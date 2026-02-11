package usecase

import (
	"errors"

	"ahu-backend/internal/repository"
	"ahu-backend/internal/usecase/dto"
)

type GetFormByInspectionUsecase struct {
	inspectionRepo repository.InspectionRepository
	formRepo       repository.FormRepository
	ahuRepo        repository.AHURepository
	scheduleRepo   repository.ScheduleRepository
}

func NewGetFormByInspectionUsecase(
	inspectionRepo repository.InspectionRepository,
	formRepo repository.FormRepository,
	ahuRepo repository.AHURepository,
	scheduleRepo repository.ScheduleRepository,
) *GetFormByInspectionUsecase {
	return &GetFormByInspectionUsecase{
		inspectionRepo: inspectionRepo,
		formRepo:       formRepo,
		ahuRepo:        ahuRepo,
		scheduleRepo:   scheduleRepo,
	}
}

func (u *GetFormByInspectionUsecase) Execute(
	inspectionID string,
) (*dto.InspectionFormResponse, error) {

	// 1️⃣ inspection
	inspection, err := u.inspectionRepo.GetByID(inspectionID)
	if err != nil || inspection == nil {
		return nil, errors.New("inspection tidak ditemukan")
	}

	// 2️⃣ form template
	form, err := u.formRepo.GetTemplateByID(inspection.FormTemplateID)
	if err != nil || form == nil {
		return nil, errors.New("form tidak ditemukan")
	}

	// 3️⃣ schedule
	schedule, err := u.scheduleRepo.GetByID(inspection.ScheduleID)
	if err != nil || schedule == nil {
		return nil, errors.New("schedule tidak ditemukan")
	}

	// 4️⃣ AHU
	ahu, err := u.ahuRepo.GetByID(schedule.AHUId)
	if err != nil || ahu == nil {
		return nil, errors.New("ahu tidak ditemukan")
	}

	ahuName := ""

	if ahu.RoomName != nil {
		ahuName = *ahu.RoomName
	} else {
		ahuName = ahu.UnitCode
	}

	// 5️⃣ response ke frontend
	return &dto.InspectionFormResponse{
		ID:       form.ID,
		Name:     form.Name,
		AHUName:  ahuName,
		Sections: form.Sections,
	}, nil
}
