package postgres

import (
	"context"
	"strings"

	"ahu-backend/internal/domain"

	"github.com/jackc/pgx/v5"
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
			created_at,
			sp.period         -- 🔥 Tambahkan ini
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
			&s.Period,
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
	a.unit_code,
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

	status = strings.TrimSpace(strings.ToLower(status))

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
		s.inspector_id,
		s.status,

		COALESCE(sp.form_template_id, (
			SELECT id FROM form_templates
			WHERE is_active = true
			LIMIT 1
		)) AS form_template_id

	FROM schedules s
	JOIN schedule_plans sp ON sp.id = s.plan_id

	WHERE s.ahu_id = $1
	  AND s.status = 'siap_diperiksa'
	  AND CURRENT_DATE BETWEEN s.start_date AND s.end_date

	ORDER BY s.created_at DESC
	LIMIT 1
	`

	var s domain.Schedule

	err := r.db.QueryRow(context.Background(), query, ahuID).Scan(
		&s.ID,
		&s.PlanID,
		&s.InspectorID,
		&s.Status,
		&s.FormTemplateID,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}

	return &s, err
}

func (r *SchedulePostgresRepository) GetByYear(year int) ([]domain.Schedule, error) {

	query := `
	SELECT
		s.id,
		s.plan_id,
		s.ahu_id,
		s.start_date,
		s.end_date,
		s.inspector_id,
		s.status,
		s.nfc_bypass,
		s.created_at
	FROM schedules s
	WHERE EXTRACT(YEAR FROM s.start_date) = $1
	ORDER BY s.start_date ASC
	`

	rows, err := r.db.Query(context.Background(), query, year)
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

func (r *SchedulePostgresRepository) ListWithDetailByYear(year int) ([]*domain.ScheduleWithDetail, error) {

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
	a.unit_code,
	a.room_name,
	a.nfc_uid,

	s.inspector_id,
	u.name AS inspector_name

	FROM schedules s
	JOIN schedule_plans sp ON sp.id = s.plan_id
	JOIN ahus a ON a.id = s.ahu_id
	LEFT JOIN users u ON u.id = s.inspector_id   -- 🔥 UBAH KE LEFT JOIN

	WHERE EXTRACT(YEAR FROM s.start_date) = $1

	ORDER BY s.start_date ASC
	`

	rows, err := r.db.Query(context.Background(), query, year)
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

func (r *SchedulePostgresRepository) SetPDFHash(year int, hash string) error {
	_, err := r.db.Exec(context.Background(), `
        UPDATE schedule_approvals
        SET pdf_hash=$1
        WHERE year=$2
    `, hash, year)

	return err
}

func (r *SchedulePostgresRepository) ListByYear(year int) ([]*domain.ScheduleWithDetail, error) {

	rows, err := r.db.Query(context.Background(), `
		SELECT
			s.id,
			s.start_date,
			s.end_date,

			a.unit_code,
			u.name

		FROM schedules s
		JOIN ahus a ON a.id = s.ahu_id
		LEFT JOIN users u ON u.id = s.inspector_id

		WHERE EXTRACT(YEAR FROM s.start_date) = $1
		ORDER BY s.start_date
	`, year)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var result []*domain.ScheduleWithDetail

	for rows.Next() {
		var d domain.ScheduleWithDetail

		rows.Scan(
			&d.ID,
			&d.StartDate,
			&d.EndDate,
			&d.UnitCode,
			&d.InspectorName,
		)

		result = append(result, &d)
	}

	return result, nil
}

func (r *SchedulePostgresRepository) ListTasksByInspector(
	inspectorID string,
) ([]domain.InspectionTask, error) {

	rows, err := r.db.Query(context.Background(), `
	SELECT
		s.id,
		s.start_date,
		s.end_date,
		s.status,
		a.unit_code,
		sp.period
	FROM schedules s
	JOIN ahus a ON a.id = s.ahu_id
	JOIN schedule_plans sp ON sp.id = s.plan_id
	WHERE s.inspector_id = $1
	AND s.status IN ('siap_diperiksa','dalam_pemeriksaan')
	ORDER BY s.start_date ASC
	`, inspectorID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []domain.InspectionTask

	for rows.Next() {
		var t domain.InspectionTask
		err := rows.Scan(
			&t.ID,
			&t.StartDate,
			&t.EndDate,
			&t.Status,
			&t.UnitCode,
			&t.Period,
		)
		if err != nil {
			return nil, err
		}

		result = append(result, t)
	}

	return result, nil
}

func (r *SchedulePlanPostgresRepository) GetByPeriod(
	period string,
	ahuID string,
) (*domain.SchedulePlan, error) {

	var plan domain.SchedulePlan

	err := r.db.QueryRow(context.Background(), `
	SELECT id, ahu_id, period
	FROM schedule_plans
	WHERE period=$1 AND ahu_id=$2
	LIMIT 1
	`,
		period,
		ahuID,
	).Scan(
		&plan.ID,
		&plan.AHUId,
		&plan.Period,
	)

	if err != nil {

		if err == pgx.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &plan, nil
}