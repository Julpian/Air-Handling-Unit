package postgres

import (
	"context"

	"ahu-backend/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type InspectionResultPostgresRepository struct {
	db *pgxpool.Pool
}

func NewInspectionResultPostgresRepository(
	db *pgxpool.Pool,
) *InspectionResultPostgresRepository {
	return &InspectionResultPostgresRepository{db: db}
}

func (r *InspectionResultPostgresRepository) SaveMany(
	results []domain.InspectionResult,
) error {

	tx, err := r.db.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	for _, res := range results {
		_, err := tx.Exec(
			context.Background(),
			`
			INSERT INTO inspection_results (
				id, inspection_id, form_item_id,
				value_text, value_number, value_bool,
				result, created_at
			)
			VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
			`,
			res.ID,
			res.InspectionID,
			res.FormItemID,
			res.ValueText,
			res.ValueNumber,
			res.ValueBool,
			res.Result,
			res.CreatedAt,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit(context.Background())
}
