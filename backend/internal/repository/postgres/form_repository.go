package postgres

import (
	"context"
	"log"

	"ahu-backend/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type FormPostgresRepository struct {
	db *pgxpool.Pool
}

func NewFormPostgresRepository(db *pgxpool.Pool) *FormPostgresRepository {
	return &FormPostgresRepository{db: db}
}

/*
=====================================================
GET TEMPLATE BY SCHEDULE
=====================================================
*/
func (r *FormPostgresRepository) GetTemplateBySchedule(scheduleID string) (*domain.FormTemplate, error) {

	var template domain.FormTemplate

	err := r.db.QueryRow(context.Background(), `
		SELECT ft.id, ft.name, ft.period
		FROM schedules s
		JOIN schedule_plans sp ON sp.id = s.plan_id
		JOIN form_templates ft ON ft.period = sp.period
		WHERE s.id = $1 AND ft.is_active = true
		ORDER BY ft.version DESC
		LIMIT 1
	`, scheduleID).Scan(&template.ID, &template.Name, &template.Period)

	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(context.Background(), `
		SELECT
			sec.id, sec.code, sec.title, sec.order_no,
			item.id, item.label, item.input_type, item.required, item.order_no
		FROM form_template_sections sec
		LEFT JOIN form_template_items item ON item.section_id = sec.id
		WHERE sec.form_template_id = $1
		ORDER BY sec.order_no, item.order_no
	`, template.ID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	sectionMap := map[string]*domain.FormSection{}

	for rows.Next() {
		var secID, secCode, secTitle string
		var secOrder int
		var item domain.FormItem

		if err := rows.Scan(
			&secID,
			&secCode,
			&secTitle,
			&secOrder,
			&item.ID,
			&item.Label,
			&item.InputType,
			&item.Required,
			&item.Order,
		); err != nil {
			return nil, err
		}

		if _, ok := sectionMap[secID]; !ok {
			sectionMap[secID] = &domain.FormSection{
				ID:    secID,
				Code:  secCode,
				Title: secTitle,
				Order: secOrder,
			}
		}

		sectionMap[secID].Items = append(sectionMap[secID].Items, item)
	}

	for _, sec := range sectionMap {
		template.Sections = append(template.Sections, *sec)
	}

	return &template, nil
}

/*
=====================================================
GET TEMPLATE BY ID
=====================================================
*/
func (r *FormPostgresRepository) GetTemplateByID(templateID string) (*domain.FormTemplate, error) {

	var template domain.FormTemplate

	err := r.db.QueryRow(context.Background(), `
		SELECT id,name,period,description,version,is_active,created_at
		FROM form_templates
		WHERE id=$1
	`, templateID).Scan(
		&template.ID,
		&template.Name,
		&template.Period,
		&template.Description,
		&template.Version,
		&template.IsActive,
		&template.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(context.Background(), `
		SELECT
			sec.id, sec.code, sec.title, sec.order_no,
			item.id, item.label, item.input_type, item.required, item.order_no
		FROM form_template_sections sec
		LEFT JOIN form_template_items item ON item.section_id = sec.id
		WHERE sec.form_template_id = $1
		ORDER BY sec.order_no, item.order_no
	`, template.ID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	sectionMap := map[string]*domain.FormSection{}

	for rows.Next() {
		var secID, secCode, secTitle string
		var secOrder int
		var item domain.FormItem

		if err := rows.Scan(
			&secID,
			&secCode,
			&secTitle,
			&secOrder,
			&item.ID,
			&item.Label,
			&item.InputType,
			&item.Required,
			&item.Order,
		); err != nil {
			return nil, err
		}

		if _, ok := sectionMap[secID]; !ok {
			sectionMap[secID] = &domain.FormSection{
				ID:    secID,
				Code:  secCode,
				Title: secTitle,
				Order: secOrder,
			}
		}

		sectionMap[secID].Items = append(sectionMap[secID].Items, item)
	}

	for _, sec := range sectionMap {
		template.Sections = append(template.Sections, *sec)
	}

	return &template, nil
}

/*
=====================================================
CREATE TEMPLATE
=====================================================
*/
func (r *FormPostgresRepository) CreateTemplate(ctx context.Context, template *domain.FormTemplate) error {

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	err = tx.QueryRow(ctx, `
		INSERT INTO form_templates (name,period,description,version,is_active)
		VALUES ($1,$2,$3,1,true)
		RETURNING id
	`, template.Name, template.Period, template.Description).Scan(&template.ID)

	if err != nil {
		return err
	}

	for _, sec := range template.Sections {

		var sectionID string

		if err := tx.QueryRow(ctx, `
			INSERT INTO form_template_sections (form_template_id,code,title,order_no)
			VALUES ($1,$2,$3,$4)
			RETURNING id
		`, template.ID, sec.Code, sec.Title, sec.Order).Scan(&sectionID); err != nil {
			return err
		}

		for _, item := range sec.Items {
			if _, err := tx.Exec(ctx, `
				INSERT INTO form_template_items
				(section_id,label,input_type,required,order_no)
				VALUES ($1,$2,$3,$4,$5)
			`, sectionID, item.Label, item.InputType, item.Required, item.Order); err != nil {
				return err
			}
		}
	}

	return tx.Commit(ctx)
}

/*
=====================================================
LIST TEMPLATES
=====================================================
*/
func (r *FormPostgresRepository) ListTemplates() ([]domain.FormTemplate, error) {

	rows, err := r.db.Query(context.Background(), `
		SELECT id,name,period,description,version,is_active,created_at
		FROM form_templates
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []domain.FormTemplate

	for rows.Next() {
		var t domain.FormTemplate
		if err := rows.Scan(
			&t.ID,
			&t.Name,
			&t.Period,
			&t.Description,
			&t.Version,
			&t.IsActive,
			&t.CreatedAt,
		); err != nil {
			log.Println("SCAN ERROR:", err)
			return nil, err
		}

		result = append(result, t)
	}

	return result, nil
}

/*
=====================================================
SET ACTIVE
=====================================================
*/
func (r *FormPostgresRepository) SetActive(id string, active bool) error {
	_, err := r.db.Exec(context.Background(),
		`UPDATE form_templates SET is_active=$2 WHERE id=$1`, id, active)
	return err
}

/*
=====================================================
CREATE NEW VERSION
=====================================================
*/
func (r *FormPostgresRepository) CreateNewVersion(ctx context.Context, oldID string, t *domain.FormTemplate) error {

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var parentID string
	var lastVersion int

	if err := tx.QueryRow(ctx,
		`SELECT COALESCE(parent_id,id),version FROM form_templates WHERE id=$1`,
		oldID,
	).Scan(&parentID, &lastVersion); err != nil {
		return err
	}

	if _, err := tx.Exec(ctx,
		`UPDATE form_templates SET is_active=false WHERE id=$1 OR parent_id=$1`,
		parentID); err != nil {
		return err
	}

	var newID string

	if err := tx.QueryRow(ctx, `
		INSERT INTO form_templates (name,period,description,version,parent_id,is_active)
		VALUES ($1,$2,$3,$4,$5,true)
		RETURNING id
	`, t.Name, t.Period, t.Description, lastVersion+1, parentID).Scan(&newID); err != nil {
		return err
	}

	for si, sec := range t.Sections {

		var secID string

		if err := tx.QueryRow(ctx, `
			INSERT INTO form_template_sections (form_template_id,code,title,order_no)
			VALUES ($1,$2,$3,$4)
			RETURNING id
		`, newID, sec.Code, sec.Title, si+1).Scan(&secID); err != nil {
			return err
		}

		for ii, item := range sec.Items {
			if _, err := tx.Exec(ctx, `
				INSERT INTO form_template_items
				(section_id,label,input_type,required,order_no)
				VALUES ($1,$2,$3,$4,$5)
			`, secID, item.Label, item.InputType, item.Required, ii+1); err != nil {
				return err
			}
		}
	}

	return tx.Commit(ctx)
}

/*
=====================================================
LIST TEMPLATE VERSIONS
=====================================================
*/
func (r *FormPostgresRepository) ListTemplateVersions(templateID string) ([]domain.FormTemplate, error) {

	var rootID string

	err := r.db.QueryRow(context.Background(), `
	WITH RECURSIVE root AS (
	  SELECT id,parent_id FROM form_templates WHERE id=$1
	  UNION ALL
	  SELECT ft.id,ft.parent_id
	  FROM form_templates ft
	  JOIN root r ON ft.id = r.parent_id
	)
	SELECT id FROM root WHERE parent_id IS NULL LIMIT 1
	`, templateID).Scan(&rootID)

	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(context.Background(), `
	SELECT id,name,period,description,version,is_active,created_at
	FROM form_templates
	WHERE id=$1 OR parent_id=$1
	ORDER BY version ASC
	`, rootID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []domain.FormTemplate

	for rows.Next() {
		var t domain.FormTemplate
		rows.Scan(
			&t.ID,
			&t.Name,
			&t.Period,
			&t.Description,
			&t.Version,
			&t.IsActive,
			&t.CreatedAt,
		)
		result = append(result, t)
	}

	return result, nil
}

func (r *FormPostgresRepository) GetActiveByPeriod(
	period string,
) (*domain.FormTemplate, error) {

	query := `
	SELECT id, name, period
	FROM form_templates
	WHERE is_active = true AND period = $1
	LIMIT 1
	`

	var f domain.FormTemplate

	err := r.db.QueryRow(context.Background(), query, period).
		Scan(&f.ID, &f.Name, &f.Period)

	if err != nil {
		return nil, err
	}

	return &f, nil
}
