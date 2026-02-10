package postgres

import (
	"context"

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
			parent_id,
			created_at
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,now())
	`

	_, err := r.db.Exec(
		context.Background(),
		query,
		i.ID,
		i.ScheduleID,
		i.InspectorID,
		i.FormTemplateID, // 🔥 INI PENTING
		i.Status,
		i.ScannedNFCUID,
		i.InspectedAt,
		i.ParentID,
	)

	return err
}

// ================= GET =================

func (r *InspectionPostgresRepository) GetByID(id string) (*domain.Inspection, error) {
	query := `
		SELECT
			id,
			schedule_id,
			inspector_id,
			form_template_id,
			status,
			note,
			scanned_nfc_uid,
			inspected_at,
			parent_id,
			created_at
		FROM inspections
		WHERE id = $1
	`

	var i domain.Inspection

	err := r.db.QueryRow(context.Background(), query, id).Scan(
		&i.ID,
		&i.ScheduleID,
		&i.InspectorID,
		&i.FormTemplateID,
		&i.Status,
		&i.Note,
		&i.ScannedNFCUID,
		&i.InspectedAt,
		&i.ParentID,
		&i.CreatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &i, nil
}

func (r *InspectionPostgresRepository) GetByScheduleID(scheduleID string) (*domain.Inspection, error) {
	return nil, nil
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

func (r *InspectionPostgresRepository) ListByStatus(status string) ([]domain.Inspection, error) {
	query := `
		SELECT
			id,
			schedule_id,
			inspector_id,
			form_template_id,
			status,
			note,
			scanned_nfc_uid,
			inspected_at,
			parent_id,
			created_at
		FROM inspections
		WHERE status = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(context.Background(), query, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var inspections []domain.Inspection

	for rows.Next() {
		var i domain.Inspection
		err := rows.Scan(
			&i.ID,
			&i.ScheduleID,
			&i.InspectorID,
			&i.FormTemplateID,
			&i.Status,
			&i.Note,
			&i.ScannedNFCUID,
			&i.InspectedAt,
			&i.ParentID,
			&i.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		inspections = append(inspections, i)
	}

	return inspections, nil
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
