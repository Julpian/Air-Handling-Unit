package usecase

import (
	"errors"
	"time"

	"ahu-backend/internal/domain"
	"ahu-backend/internal/repository"
	"ahu-backend/internal/usecase/dto"

	"github.com/google/uuid"
)

type ScanNFCUsecase struct {
	ahuRepo        repository.AHURepository
	scheduleRepo   repository.ScheduleRepository
	inspectionRepo repository.InspectionRepository
}

func NewScanNFCUsecase(
	ahuRepo repository.AHURepository,
	scheduleRepo repository.ScheduleRepository,
	inspectionRepo repository.InspectionRepository,
) *ScanNFCUsecase {
	return &ScanNFCUsecase{
		ahuRepo:        ahuRepo,
		scheduleRepo:   scheduleRepo,
		inspectionRepo: inspectionRepo,
	}
}

func (uc *ScanNFCUsecase) Execute(
	req dto.ScanNFCRequest,
	inspectorID string,
) (*dto.ScanNFCResponse, error) {

	if req.NFCUID == "" {
		return nil, errors.New("nfc uid wajib")
	}

	// 1️⃣ AHU
	ahu, err := uc.ahuRepo.GetByNFCUID(req.NFCUID)
	if err != nil || ahu == nil {
		return nil, errors.New("nfc tidak terdaftar")
	}

	// 2️⃣ Schedule aktif
	schedule, err := uc.scheduleRepo.GetActiveByAHUAndInspector(
		ahu.ID,
		inspectorID,
	)
	if err != nil || schedule == nil {
		return nil, errors.New("tidak ada jadwal aktif")
	}

	// 3️⃣ Inspection
	inspection, err := uc.inspectionRepo.GetByScheduleID(schedule.ID)
	if err != nil {
		return nil, err
	}

	if schedule.InspectorID == nil {
		return nil, errors.New("inspector belum ditentukan")
	}

	// 4️⃣ JIKA BELUM ADA → BUAT BARU
	if inspection == nil {

		if schedule.FormTemplateID == "" {
			return nil, errors.New(
				"schedule tidak punya form_template_id (cek schedule_plans)",
			)
		}

		if schedule.InspectorID == nil {
			return nil, errors.New("inspector belum ditentukan")
		}

		inspection = &domain.Inspection{
			ID:             uuid.NewString(),
			ScheduleID:     schedule.ID,
			InspectorID:    *schedule.InspectorID,
			FormTemplateID: schedule.FormTemplateID,
			Status:         "draft",
			CreatedAt:      time.Now(),
		}

		if err := uc.inspectionRepo.Create(inspection); err != nil {
			return nil, err
		}
	}

	// 5️⃣ Generate token
	token := uuid.NewString()
	exp := time.Now().Add(10 * time.Minute)

	if err := uc.inspectionRepo.SetScanToken(
		inspection.ID,
		token,
		exp,
		req.NFCUID,
	); err != nil {
		return nil, err
	}

	// 🔥 INI YANG HILANG SELAMA INI
	if err := uc.scheduleRepo.UpdateStatus(
		schedule.ID,
		domain.ScheduleStatusDalamPemeriksaan,
	); err != nil {
		return nil, err
	}

	return &dto.ScanNFCResponse{
		InspectionID: inspection.ID,
		ScanToken:    token,
	}, nil
}
