package postgres

import (
	"context"

	"ahu-backend/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type SchedulePostgresRepository struct {
	db *pgxpool.Pool
}

func NewSchedulePostgresRepository(db *pgxpool.Pool) *SchedulePostgresRepository {
	return &SchedulePostgresRepository{db: db}
}

/* ================= CREATE ================= */

func (r *SchedulePostgresRepository) Create(s *domain.Schedule) error {
	query := `
		INSERT INTO schedules (
			id,
			plan_id,
			ahu_id,
			start_date,
			end_date,
			inspector_id,
			status,
			nfc_bypass,
			created_at
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
	`

	_, err := r.db.Exec(context.Background(), query,
		s.ID,
		s.PlanID,
		s.AHUId,
		s.StartDate,
		s.EndDate,
		s.InspectorID,
		s.Status,
		s.NFCBypass,
		s.CreatedAt,
	)

	return err
}

/* ================= DELETE ================= */

func (r *SchedulePostgresRepository) DeleteByPlan(planID string) error {
	_, err := r.db.Exec(
		context.Background(),
		`DELETE FROM schedules WHERE plan_id = $1`,
		planID,
	)
	return err
}

/* ================= ASSIGN INSPECTOR ================= */

func (r *SchedulePostgresRepository) AssignInspector(
	scheduleID string,
	inspectorID string,
) error {

	_, err := r.db.Exec(
		context.Background(),
		`
		UPDATE schedules
		SET inspector_id = $1
		WHERE id = $2
		`,
		inspectorID,
		scheduleID,
	)

	return err
}

/* ================= LIST SIMPLE ================= */

func (r *SchedulePostgresRepository) ListAll() ([]domain.Schedule, error) {
	query := `
		SELECT
			id,
			plan_id,
			ahu_id,
			start_date,
			end_date,
			inspector_id,
			status,
			nfc_bypass,
			created_at
		FROM schedules
		ORDER BY start_date ASC
	`

	rows, err := r.db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []domain.Schedule

	for rows.Next() {
		var s domain.Schedule
		if err := rows.Scan(
			&s.ID,
			&s.PlanID,
			&s.AHUId,
			&s.StartDate,
			&s.EndDate,
			&s.InspectorID,
			&s.Status,
			&s.NFCBypass,
			&s.CreatedAt,
		); err != nil {
			return nil, err
		}
		result = append(result, s)
	}

	return result, nil
}

/* ================= LIST WITH DETAIL (ADMIN UI) ================= */

func (r *SchedulePostgresRepository) ListWithDetail() ([]*domain.ScheduleWithDetail, error) {
	query := `
		SELECT
		s.id,
		s.start_date,
		s.end_date,
		s.status,
		s.nfc_bypass,

		sp.id AS plan_id,
		sp.period,
		sp.week_of_month,
		sp.month,

		a.id AS ahu_id,
		a.unit_code AS unit_code,
		a.room_name,
		a.nfc_uid,

		s.inspector_id,
		u.name AS inspector_name

		FROM schedules s
		JOIN schedule_plans sp ON sp.id = s.plan_id
		JOIN ahus a ON a.id = s.ahu_id
		LEFT JOIN users u ON u.id = s.inspector_id
		ORDER BY s.start_date ASC
	`

	rows, err := r.db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*domain.ScheduleWithDetail

	for rows.Next() {
		var d domain.ScheduleWithDetail

		if err := rows.Scan(
			&d.ID,
			&d.StartDate,
			&d.EndDate,
			&d.Status,
			&d.NFCBypass,

			&d.PlanID,
			&d.Period,
			&d.WeekOfMonth,
			&d.Month,

			&d.AHUID,
			&d.UnitCode,
			&d.RoomName,
			&d.NFCUID,

			&d.InspectorID,
			&d.InspectorName,
		); err != nil {
			return nil, err
		}

		result = append(result, &d)
	}

	return result, nil
}

/* ================= GET BY ID ================= */

func (r *SchedulePostgresRepository) GetByID(id string) (*domain.Schedule, error) {
	query := `
		SELECT
			id,
			plan_id,
			ahu_id,
			start_date,
			end_date,
			inspector_id,
			status,
			nfc_bypass,
			created_at
		FROM schedules
		WHERE id = $1
	`

	var s domain.Schedule

	err := r.db.QueryRow(
		context.Background(),
		query,
		id,
	).Scan(
		&s.ID,
		&s.PlanID,
		&s.AHUId,
		&s.StartDate,
		&s.EndDate,
		&s.InspectorID,
		&s.Status,
		&s.NFCBypass,
		&s.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &s, nil
}

/* ================= UPDATE ================= */

func (r *SchedulePostgresRepository) Update(s *domain.Schedule) error {
	_, err := r.db.Exec(
		context.Background(),
		`
		UPDATE schedules
		SET
			start_date   = $1,
			end_date     = $2,
			inspector_id = $3,
			status       = $4,
			nfc_bypass   = $5
		WHERE id = $6
		`,
		s.StartDate,
		s.EndDate,
		s.InspectorID,
		s.Status,
		s.NFCBypass,
		s.ID,
	)

	return err
}

/* ================= UPDATE STATUS ================= */

func (r *SchedulePostgresRepository) UpdateStatus(
	scheduleID string,
	status string,
) error {

	_, err := r.db.Exec(
		context.Background(),
		`
		UPDATE schedules
		SET status = $1
		WHERE id = $2
		`,
		status,
		scheduleID,
	)

	return err
}

func (r *SchedulePostgresRepository) GetActiveByAHU(
	ahuID string,
) (*domain.Schedule, error) {

	query := `
		SELECT
			s.id,
			s.plan_id,
			s.status,
			s.created_at
		FROM schedules s
		JOIN schedule_plans sp ON sp.id = s.plan_id
		WHERE sp.ahu_id = $1
		  AND s.status = 'siap_diperiksa'
		ORDER BY s.created_at DESC
		LIMIT 1
	`

	row := r.db.QueryRow(context.Background(), query, ahuID)

	var s domain.Schedule
	err := row.Scan(
		&s.ID,
		&s.PlanID,
		&s.Status,
		&s.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &s, nil
}
