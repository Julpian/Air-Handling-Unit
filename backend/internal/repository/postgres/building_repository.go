package postgres

import (
	"context"

	"ahu-backend/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type BuildingPostgresRepository struct {
	db *pgxpool.Pool
}

func NewBuildingPostgresRepository(db *pgxpool.Pool) *BuildingPostgresRepository {
	return &BuildingPostgresRepository{db: db}
}

func (r *BuildingPostgresRepository) Create(
	b *domain.Building,
	createdBy string,
) error {
	_, err := r.db.Exec(context.Background(), `
		INSERT INTO buildings (
			id,
			name,
			description,
			is_active,
			created_by,
			created_at
		) VALUES ($1, $2, $3, true, $4, now())
	`,
		b.ID,
		b.Name,
		b.Description,
		createdBy,
	)

	return err
}

func (r *BuildingPostgresRepository) ListAll() ([]domain.Building, error) {
	rows, err := r.db.Query(context.Background(), `
		SELECT id, name, description, is_active, created_at
		FROM buildings
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var buildings []domain.Building

	for rows.Next() {
		var b domain.Building
		var desc *string // ← penampung NULL

		err := rows.Scan(
			&b.ID,
			&b.Name,
			&desc,
			&b.IsActive,
			&b.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		b.Description = desc
		buildings = append(buildings, b)
	}

	return buildings, nil
}

func (r *BuildingPostgresRepository) Update(b *domain.Building) error {
	_, err := r.db.Exec(context.Background(), `
		UPDATE buildings
		SET name = $1,
		    description = $2
		WHERE id = $3
	`, b.Name, b.Description, b.ID)

	return err
}

func (r *BuildingPostgresRepository) Deactivate(id string) error {
	_, err := r.db.Exec(context.Background(), `
		UPDATE buildings
		SET is_active = false
		WHERE id = $1
	`, id)

	return err
}

func (r *BuildingPostgresRepository) GetByID(
	id string,
) (*domain.Building, error) {

	var b domain.Building

	err := r.db.QueryRow(context.Background(), `
		SELECT
			id,
			name,
			description,
			is_active,
			created_at
		FROM buildings
		WHERE id = $1
	`, id).Scan(
		&b.ID,
		&b.Name,
		&b.Description,
		&b.IsActive,
		&b.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &b, nil
}
