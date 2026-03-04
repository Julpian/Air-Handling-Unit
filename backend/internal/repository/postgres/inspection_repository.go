package postgres

import (
	"context"
	"fmt"
	"time"

	"ahu-backend/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type InspectionPostgresRepository struct {
	db *pgxpool.Pool
}

// ✅ CONSTRUCTOR
func NewInspectionPostgresRepository(db *pgxpool.Pool) *InspectionPostgresRepository {
	return &InspectionPostgresRepository{db: db}
}

// ================= CREATE =================

func (r *InspectionPostgresRepository) Create(i *domain.Inspection) error {
	query := `
	INSERT INTO inspections (
		id,
		schedule_id,
		inspector_id,
		form_template_id,
		status,
		scanned_nfc_uid,
		inspected_at,
		scan_token,
		scan_token_expires_at,
		parent_id,
		created_at
	)
	VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,now())
	`

	_, err := r.db.Exec(
		context.Background(),
		query,
		i.ID,
		i.ScheduleID,
		i.InspectorID,
		i.FormTemplateID,
		i.Status,
		i.ScannedNFCUID,
		i.InspectedAt,      // ✅
		i.ScanToken,        // ✅
		i.ScanTokenExpires, // ✅
		i.ParentID,
	)

	return err
}

// ================= GET =================

func (r *InspectionPostgresRepository) GetByID(id string) (*domain.Inspection, error) {
	query := `
    SELECT
        id, schedule_id, inspector_id, form_template_id, status,
        note, scanned_nfc_uid, inspected_at, scan_token, scan_token_expires_at,
        parent_id, 
        COALESCE(inspector_signature, ''), 
        COALESCE(spv_signature, ''), 
        approved_by, 
        created_at
    FROM inspections
    WHERE id = $1
    `

	var i domain.Inspection
	// 🔥 FIX: Scan harus berjumlah 15 variabel sesuai urutan SELECT
	err := r.db.QueryRow(context.Background(), query, id).Scan(
		&i.ID,                 // 1
		&i.ScheduleID,         // 2
		&i.InspectorID,        // 3
		&i.FormTemplateID,     // 4
		&i.Status,             // 5
		&i.Note,               // 6
		&i.ScannedNFCUID,      // 7
		&i.InspectedAt,        // 8
		&i.ScanToken,          // 9
		&i.ScanTokenExpires,   // 10
		&i.ParentID,           // 11
		&i.InspectorSignature, // 12
		&i.SPVSignature,       // 13
		&i.ApprovedBy,         // 14
		&i.CreatedAt,          // 15
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &i, nil
}
func (r *InspectionPostgresRepository) GetByScheduleID(
	scheduleID string,
) (*domain.Inspection, error) {

	var i domain.Inspection

	err := r.db.QueryRow(context.Background(), `
	SELECT id, schedule_id, status, form_template_id
	FROM inspections
	WHERE schedule_id=$1
	LIMIT 1
	`, scheduleID).Scan(
		&i.ID,
		&i.ScheduleID,
		&i.Status,
		&i.FormTemplateID,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &i, nil
}

func (r *InspectionPostgresRepository) UpdateStatus(
	id string,
	status string,
	note *string,
) error {

	query := `
		UPDATE inspections
		SET status = $1
		WHERE id = $2
	`

	_, err := r.db.Exec(
		context.Background(),
		query,
		status,
		id,
	)

	return err
}

func (r *InspectionPostgresRepository) Approve(
	id string,
	approverID string,
	approvedAt any,
	metadata map[string]any,
) error {
	return nil
}

func (r *InspectionPostgresRepository) ExistsApproved(scheduleID string) (bool, error) {
	return false, nil
}

func (r *InspectionPostgresRepository) GetLastByScheduleID(scheduleID string) (*domain.Inspection, error) {
	return nil, nil
}

func (r *InspectionPostgresRepository) ListByStatus(status string, inspectorID string) ([]domain.Inspection, error) {
	var rows pgx.Rows
	var err error

	// Gunakan ::text pada kolom inspector_id agar bisa dibandingkan dengan string kosong
	query := `
        SELECT id, schedule_id, inspector_id, form_template_id, status, created_at
        FROM inspections
        WHERE ($1 = '' OR status = $1)
          AND ($2 = '' OR inspector_id::text = $2) -- 🔥 FIX: Cast UUID ke text
        ORDER BY created_at DESC
    `

	rows, err = r.db.Query(context.Background(), query, status, inspectorID)
	if err != nil {
		fmt.Println("DATABASE ERROR:", err) // Tambahkan log untuk debug
		return nil, err
	}
	defer rows.Close()

	var list []domain.Inspection
	for rows.Next() {
		var i domain.Inspection
		err := rows.Scan(&i.ID, &i.ScheduleID, &i.InspectorID, &i.FormTemplateID, &i.Status, &i.CreatedAt)
		if err != nil {
			return nil, err
		}
		list = append(list, i)
	}
	return list, nil
}

func (r *InspectionPostgresRepository) SaveResult(
	res *domain.InspectionResult,
) error {

	_, err := r.db.Exec(context.Background(), `
	INSERT INTO inspection_results (
		id,
		inspection_id,
		item_id,
		value,
		created_at
	)
	VALUES ($1,$2,$3,$4,now())
	`,
		res.ID,
		res.InspectionID,
		res.ItemID,
		res.Value,
	)

	return err
}

func (r *InspectionPostgresRepository) ClearScanToken(id string) error {
	_, err := r.db.Exec(context.Background(), `
	UPDATE inspections
	SET scan_token = NULL,
	    scan_token_expires_at = NULL
	WHERE id = $1
	`, id)

	return err
}

func (r *InspectionPostgresRepository) SetScanToken(
	id string,
	token string,
	exp time.Time,
	uid string,
) error {

	_, err := r.db.Exec(context.Background(), `
	UPDATE inspections
	SET scan_token=$1,
	    scan_token_expires_at=$2,
	    scanned_nfc_uid=$3
	WHERE id=$4
	`,
		token,
		exp,
		uid,
		id,
	)

	return err
}

func (r *InspectionPostgresRepository) GetInspectionReport(id string) (*domain.InspectionReport, error) {
	rows, err := r.db.Query(context.Background(), `
    SELECT
        i.id, s.id, i.status, a.unit_code, u.name,
        COALESCE(i.inspector_signature, ''),
        COALESCE(i.spv_signature, ''),
        
        -- Master Data AHU & Area
        ar.name as area_name,
        a.cleanliness_class,
        a.vendor,
        a.room_name,
        i.inspected_at,
        sp.period,          -- 🔥 1. Tambahkan kolom period dari schedule_plans

        fs.code, fs.title,
        fi.label, COALESCE(ir.value_text,''), ir.result
    FROM inspections i
    JOIN schedules s ON s.id = i.schedule_id
    JOIN schedule_plans sp ON sp.id = s.plan_id -- 🔥 2. Join ke tabel schedule_plans
    JOIN ahus a ON a.id = s.ahu_id
    JOIN areas ar ON ar.id = a.area_id
    JOIN users u ON u.id = i.inspector_id
    JOIN inspection_results ir ON ir.inspection_id = i.id
    JOIN form_template_items fi ON fi.id = ir.form_item_id
    JOIN form_template_sections fs ON fs.id = fi.section_id
    WHERE i.id = $1
    ORDER BY fs.code
    `, id)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	report := domain.InspectionReport{}
	sectionMap := map[string]int{}

	for rows.Next() {
		var secCode, secTitle string
		var label, value, result, signature, spvSignature string
		var areaName, cleanliness, vendor, roomName string
		var inspectedAt *time.Time
		var period string // 🔥 3. Tambahkan variabel penampung period

		// 🔥 4. Update urutan Scan (Total 18 kolom sekarang)
		err := rows.Scan(
			&report.InspectionID, &report.ScheduleID, &report.Status, &report.UnitCode, &report.Inspector, // 1-5
			&signature, &spvSignature, // 6-7
			&areaName, &cleanliness, &vendor, &roomName, &inspectedAt, // 8-12
			&period,             // 🔥 13. Scan data period di sini
			&secCode, &secTitle, // 14-15
			&label, &value, &result, // 16-18
		)
		if err != nil {
			return nil, err
		}

		// Masukkan data master ke struct
		report.Signature = signature
		report.SPVSignature = spvSignature
		report.AreaName = areaName
		report.CleanlinessClass = cleanliness
		report.Vendor = vendor
		report.RoomName = roomName
		report.InspectedAt = inspectedAt
		report.Period = period // 🔥 5. Set nilai Period ke struct report

		// grouping section
		idx, ok := sectionMap[secCode]
		if !ok {
			report.Sections = append(report.Sections, domain.InspectionReportSection{
				Code: secCode, Title: secTitle,
			})
			idx = len(report.Sections) - 1
			sectionMap[secCode] = idx
		}
		report.Sections[idx].Items = append(report.Sections[idx].Items, domain.InspectionReportItem{
			Label: label, Value: value, Result: result,
		})
	}
	return &report, nil
}

func (r *InspectionPostgresRepository) SignInspection(id string, signature string) error {
	_, err := r.db.Exec(context.Background(), `
    UPDATE inspections
    SET inspector_signature=$1,
        inspector_signed_at=now(),
        status='waiting_spv' -- 🔥 UBAH DARI 'signed' KE 'waiting_spv'
    WHERE id=$2
    `,
		signature,
		id,
	)
	return err
}

func (r *InspectionPostgresRepository) SaveSignature(id string, signature string) error {

	_, err := r.db.Exec(context.Background(), `
	UPDATE inspections
	SET inspector_signature = $1
	WHERE id = $2
	`, signature, id)

	return err
}

func (r *InspectionPostgresRepository) ApproveInspection(id string, spvID string, signature string) error {
	res, err := r.db.Exec(context.Background(), `
        UPDATE inspections 
        SET spv_signature=$1, spv_signed_at=now(), approved_by=$2, status='approved' 
        WHERE id=$3 AND status = 'waiting_spv'
    `, signature, spvID, id)

	if err != nil {
		return err
	}

	// 🔥 Tambahan: Beri tahu jika tidak ada data yang terupdate (ID salah atau status bukan waiting_spv)
	if res.RowsAffected() == 0 {
		return fmt.Errorf("inspeksi tidak ditemukan atau sudah diproses")
	}

	return nil
}

func (r *InspectionPostgresRepository) GetVerificationData(id string) (*domain.InspectionReport, error) {
	// 🔥 PERBAIKAN: Tambahkan JOIN agar tabel a, sp, u_ins, dan u_spv terbaca
	query := `
    SELECT 
        i.id, a.unit_code, sp.period, 
        u_ins.name as inspector_name, 
        COALESCE(u_spv.name, 'Awaiting Approval') as spv_name, 
        i.inspected_at,
        i.status 
    FROM inspections i
    JOIN schedules s ON s.id = i.schedule_id           -- 🔥 Wajib ada
    JOIN schedule_plans sp ON sp.id = s.plan_id         -- 🔥 Wajib ada
    JOIN ahus a ON a.id = s.ahu_id                      -- 🔥 Wajib ada
    JOIN users u_ins ON u_ins.id = i.inspector_id       -- 🔥 Wajib ada
    LEFT JOIN users u_spv ON u_spv.id = i.approved_by   -- 🔥 Wajib ada (LEFT JOIN karena mungkin belum ada SPV)
    WHERE i.id = $1 AND i.status IN ('waiting_spv', 'approved')
    `
	var d domain.InspectionReport
	var inspectedAt *time.Time

	// Scan sekarang sudah benar berjumlah 7 variabel
	err := r.db.QueryRow(context.Background(), query, id).Scan(
		&d.InspectionID, // 1
		&d.UnitCode,     // 2
		&d.Period,       // 3
		&d.Inspector,    // 4
		&d.SPVName,      // 5
		&inspectedAt,    // 6
		&d.Status,       // 7
	)
	if err != nil {
		return nil, err
	}

	d.InspectedAt = inspectedAt
	return &d, nil
}
