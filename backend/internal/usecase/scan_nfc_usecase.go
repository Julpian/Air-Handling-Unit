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
	formRepo       repository.FormRepository // ✅ TAMBAH
}

func NewScanNFCUsecase(
	ahuRepo repository.AHURepository,
	scheduleRepo repository.ScheduleRepository,
	inspectionRepo repository.InspectionRepository,
	formRepo repository.FormRepository, // ✅ TAMBAH
) *ScanNFCUsecase {
	return &ScanNFCUsecase{
		ahuRepo:        ahuRepo,
		scheduleRepo:   scheduleRepo,
		inspectionRepo: inspectionRepo,
		formRepo:       formRepo, // ✅ SIMPAN
	}
}

func (uc *ScanNFCUsecase) Execute(
	req dto.ScanNFCRequest,
	inspectorID string,
) (*dto.ScanNFCResponse, error) {

	if req.NFCUID == "" {
		return nil, errors.New("nfc uid wajib")
	}

	// 1️⃣ Ambil AHU
	ahu, err := uc.ahuRepo.GetByNFCUID(req.NFCUID)
	if err != nil {
		return nil, err
	}
	if ahu == nil {
		return nil, errors.New("nfc tidak terdaftar")
	}

	// 2️⃣ Ambil schedule aktif
	schedule, err := uc.scheduleRepo.GetActiveByAHU(ahu.ID)
	if err != nil {
		return nil, err
	}
	if schedule == nil {
		return nil, errors.New("tidak ada jadwal aktif untuk AHU ini")
	}

	// 3️⃣ Validasi inspector
	if schedule.InspectorID == nil {
		return nil, errors.New("schedule belum punya inspector")
	}
	if *schedule.InspectorID != inspectorID {
		return nil, errors.New("schedule bukan milik anda")
	}

	// 4️⃣ Validasi status
	if schedule.Status != domain.ScheduleStatusSiapDiperiksa {
		return nil, errors.New("schedule sudah discan atau selesai")
	}

	// 5️⃣ Cek inspection existing
	existing, err := uc.inspectionRepo.GetByScheduleID(schedule.ID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("inspection sudah dibuat")
	}

	// ===============================
	// 🔥 AMBIL FORM BERDASARKAN PERIOD
	// ===============================
	form, err := uc.formRepo.GetTemplateBySchedule(schedule.ID)
	if err != nil {
		return nil, err
	}
	if form == nil {
		return nil, errors.New("form template tidak ditemukan untuk schedule ini")
	}

	// 6️⃣ Generate token
	token := uuid.NewString()
	exp := time.Now().Add(10 * time.Minute)

	now := time.Now()

	inspection := &domain.Inspection{
		ID:             uuid.NewString(),
		ScheduleID:     schedule.ID,
		InspectorID:    inspectorID,
		FormTemplateID: form.ID, // ✅ FIX UTAMA

		Status: domain.InspectionStatusSedangDiisi,

		ScannedNFCUID: &req.NFCUID,
		InspectedAt:   &now,

		ScanToken:        &token,
		ScanTokenExpires: &exp,
	}

	// 7️⃣ Create inspection
	if err := uc.inspectionRepo.Create(inspection); err != nil {
		return nil, err
	}

	// 8️⃣ Update schedule → dalam pemeriksaan
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
