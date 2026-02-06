package postgres

import (
	"context"
	"database/sql"

	"ahu-backend/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type SchedulePlanPostgresRepository struct {
	db *pgxpool.Pool
}

func NewSchedulePlanPostgresRepository(db *pgxpool.Pool) *SchedulePlanPostgresRepository {
	return &SchedulePlanPostgresRepository{db: db}
}

/* ================= CREATE ================= */

func (r *SchedulePlanPostgresRepository) Create(plan *domain.SchedulePlan) error {
	query := `
		INSERT INTO schedule_plans
		(id, ahu_id, unit_code, period, week_of_month, month, status)
		VALUES ($1,$2,$3,$4,$5,$6,'draft')
	`

	_, err := r.db.Exec(
		context.Background(),
		query,
		plan.ID,
		plan.AHUId,
		plan.UnitCode,
		plan.Period,
		plan.WeekOfMonth,
		plan.Month,
	)

	return err
}

/* ================= LIST (DRAFT ONLY) ================= */

func (r *SchedulePlanPostgresRepository) ListAllWithAHU() (
	[]domain.SchedulePlanWithAHU,
	error,
) {
	query := `
		SELECT
			sp.id,
			sp.ahu_id,
			a.unit_code,
			sp.period,
			sp.week_of_month,
			sp.month,
			sp.status,
			sp.created_at
		FROM schedule_plans sp
		JOIN ahus a ON a.id = sp.ahu_id
		ORDER BY sp.created_at DESC
	`

	rows, err := r.db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []domain.SchedulePlanWithAHU

	for rows.Next() {
		var p domain.SchedulePlanWithAHU
		var month sql.NullInt32

		err := rows.Scan(
			&p.ID,
			&p.AHUId,
			&p.UnitCode,
			&p.Period,
			&p.WeekOfMonth,
			&month,
			&p.Status,
			&p.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		if month.Valid {
			m := int(month.Int32)
			p.Month = &m
		}

		list = append(list, p)
	}

	return list, nil
}

/* ================= GET ACTIVE ================= */

func (r *SchedulePlanPostgresRepository) GetActiveByYear(
	year int,
) ([]domain.SchedulePlan, error) {

	query := `
		SELECT
			id, ahu_id, unit_code, period, week_of_month, month, created_at
		FROM schedule_plans
		WHERE
			status = 'draft'
			AND EXTRACT(YEAR FROM created_at) = $1
	`

	rows, err := r.db.Query(context.Background(), query, year)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []domain.SchedulePlan
	for rows.Next() {
		var p domain.SchedulePlan
		err := rows.Scan(
			&p.ID,
			&p.AHUId,
			&p.UnitCode,
			&p.Period,
			&p.WeekOfMonth,
			&p.Month,
			&p.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		list = append(list, p)
	}

	return list, nil
}

/* ================= UPDATE ================= */

func (r *SchedulePlanPostgresRepository) Update(plan *domain.SchedulePlan) error {
	query := `
		UPDATE schedule_plans
		SET period=$1, week_of_month=$2, month=$3
		WHERE id=$4 AND status='draft'
	`
	_, err := r.db.Exec(
		context.Background(),
		query,
		plan.Period,
		plan.WeekOfMonth,
		plan.Month,
		plan.ID,
	)
	return err
}

/* ================= DELETE ================= */

func (r *SchedulePlanPostgresRepository) Delete(id string) error {
	_, err := r.db.Exec(
		context.Background(),
		`DELETE FROM schedule_plans WHERE id = $1 AND status='draft'`,
		id,
	)
	return err
}

/* ================= MARK GENERATED ================= */

func (r *SchedulePlanPostgresRepository) MarkGenerated(planID string) error {
	_, err := r.db.Exec(
		context.Background(),
		`UPDATE schedule_plans SET status = 'generated' WHERE id = $1`,
		planID,
	)
	return err
}
