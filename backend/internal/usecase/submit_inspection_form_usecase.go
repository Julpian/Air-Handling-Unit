package usecase

import (
	"errors"
	"time"

	"ahu-backend/internal/domain"
	"ahu-backend/internal/repository"
	"ahu-backend/internal/usecase/dto"

	"github.com/google/uuid"
)

type SubmitInspectionFormUsecase struct {
	resultRepo   repository.InspectionResultRepository
	inspectRepo  repository.InspectionRepository
	formRepo     repository.FormRepository
	scheduleRepo repository.ScheduleRepository //
}

func NewSubmitInspectionFormUsecase(
	resultRepo repository.InspectionResultRepository,
	inspectionRepo repository.InspectionRepository,
	formRepo repository.FormRepository,
	scheduleRepo repository.ScheduleRepository, // 🔥 TAMBAH
) *SubmitInspectionFormUsecase {
	return &SubmitInspectionFormUsecase{
		resultRepo:   resultRepo,
		inspectRepo:  inspectionRepo,
		formRepo:     formRepo,
		scheduleRepo: scheduleRepo,
	}
}

func (uc *SubmitInspectionFormUsecase) Execute(
	inspectionID string,
	req dto.SubmitInspectionFormRequest,
) error {

	// ===============================
	// GET INSPECTION
	// ===============================
	inspection, err := uc.inspectRepo.GetByID(inspectionID)
	if err != nil {
		return err
	}
	if inspection == nil {
		return errors.New("inspection tidak ditemukan")
	}

	if inspection.FormTemplateID == "" {
		return errors.New("inspection belum memiliki form template")
	}

	// ===============================
	// GET FORM TEMPLATE
	// ===============================
	form, err := uc.formRepo.GetTemplateByID(inspection.FormTemplateID)
	if err != nil {
		return err
	}

	// ===============================
	// MAP FORM ITEMS
	// ===============================
	itemMap := map[string]domain.FormItem{}
	for _, sec := range form.Sections {
		for _, item := range sec.Items {
			itemMap[item.ID] = item
		}
	}

	// ===============================
	// PROCESS ANSWERS
	// ===============================
	var results []domain.InspectionResult
	finalResult := "pass"

	for _, ans := range req.Items {
		item, ok := itemMap[ans.FormItemID]
		if !ok {
			return errors.New("form item tidak valid")
		}

		res := domain.InspectionResult{
			ID:           uuid.NewString(),
			InspectionID: inspectionID,
			FormItemID:   ans.FormItemID,
			ValueText:    ans.ValueText,
			ValueNumber:  ans.ValueNumber,
			ValueBool:    ans.ValueBool,
			CreatedAt:    time.Now(),
			Result:       "pass",
		}

		// ===============================
		// VALIDATION REQUIRED
		// ===============================
		if item.Required {
			isEmpty :=
				ans.ValueText == nil &&
					ans.ValueNumber == nil &&
					ans.ValueBool == nil

			if isEmpty {
				res.Result = "fail"
				finalResult = "fail"
			}
		}

		results = append(results, res)
	}

	// ===============================
	// SAVE RESULTS
	// ===============================
	if err := uc.resultRepo.SaveMany(results); err != nil {
		return err
	}

	// ===============================
	// UPDATE SCHEDULE STATUS
	// ===============================

	scheduleStatus := domain.ScheduleStatusSelesai

	if finalResult == "fail" {
		scheduleStatus = domain.ScheduleStatusRevisi
	}

	if err := uc.scheduleRepo.UpdateStatus(
		inspection.ScheduleID,
		scheduleStatus,
	); err != nil {
		return err
	}

	// ===============================
	// UPDATE INSPECTION STATUS
	// ===============================
	status := "inspected"
	if finalResult == "fail" {
		status = "rejected"
	}

	uc.inspectRepo.ClearScanToken(inspectionID)

	return uc.inspectRepo.UpdateStatus(
		inspectionID,
		status,
		nil,
	)
}

func (uc *SubmitInspectionFormUsecase) GetInspection(id string) (*domain.Inspection, error) {
	return uc.inspectRepo.GetByID(id)
}
