package postgres

import (
	"context"

	"ahu-backend/internal/repository"

	"github.com/jackc/pgx/v5/pgxpool"
)

type FormWritePostgresRepository struct {
	db *pgxpool.Pool
}

func NewFormWritePostgresRepository(db *pgxpool.Pool) *FormWritePostgresRepository {
	return &FormWritePostgresRepository{db: db}
}

func (r *FormWritePostgresRepository) CreateFormTemplate(
	ctx context.Context,
	req repository.CreateFormTemplateRequest,
) error {

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var formID string
	err = tx.QueryRow(ctx, `
		INSERT INTO form_templates (name, period, description)
		VALUES ($1, $2, $3)
		RETURNING id
	`, req.Name, req.Period, req.Description).Scan(&formID)
	if err != nil {
		return err
	}

	for _, sec := range req.Sections {
		var sectionID string
		err = tx.QueryRow(ctx, `
			INSERT INTO form_template_sections
			(form_template_id, code, title, order_no)
			VALUES ($1, $2, $3, $4)
			RETURNING id
		`, formID, sec.Code, sec.Title, sec.Order).Scan(&sectionID)
		if err != nil {
			return err
		}

		for _, item := range sec.Items {
			_, err = tx.Exec(ctx, `
				INSERT INTO form_template_items
				(section_id, label, input_type, required, options, order_no)
				VALUES ($1, $2, $3, $4, $5, $6)
			`, sectionID, item.Label, item.Type, item.Required, item.Options, item.Order)
			if err != nil {
				return err
			}
		}
	}

	return tx.Commit(ctx)
}
