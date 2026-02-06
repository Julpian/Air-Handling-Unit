package postgres

import (
	"context"
	"time"

	"ahu-backend/internal/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AreaPostgresRepository struct {
	db *pgxpool.Pool
}

func NewAreaPostgresRepository(db *pgxpool.Pool) *AreaPostgresRepository {
	return &AreaPostgresRepository{db: db}
}

// ================= CREATE =================
func (r *AreaPostgresRepository) Create(area *domain.Area) error {
	query := `
		INSERT INTO areas (
			id,
			building_id,
			name,
			floor,
			cleanliness_class,
			is_active,
			created_at
		) VALUES ($1, $2, $3, $4, $5, true, $6)
	`

	_, err := r.db.Exec(
		context.Background(),
		query,
		uuid.NewString(),
		area.BuildingID,
		area.Name,
		area.Floor,
		area.CleanlinessClass,
		time.Now(),
	)

	return err
}

// ================= LIST =================
func (r *AreaPostgresRepository) ListAll() ([]domain.Area, error) {
	query := `
		SELECT
			id,
			building_id,
			name,
			floor,
			cleanliness_class,
			is_active
		FROM areas
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []domain.Area

	for rows.Next() {
		var a domain.Area

		// ⚠️ URUTAN HARUS SAMA DENGAN SELECT
		if err := rows.Scan(
			&a.ID,
			&a.BuildingID,
			&a.Name,
			&a.Floor,
			&a.CleanlinessClass,
			&a.IsActive,
		); err != nil {
			return nil, err
		}

		list = append(list, a)
	}

	return list, nil
}

// ================= GET BY ID =================
func (r *AreaPostgresRepository) GetByID(id string) (*domain.Area, error) {
	query := `
		SELECT
			id,
			building_id,
			name,
			floor,
			cleanliness_class,
			is_active
		FROM areas
		WHERE id = $1
	`

	var a domain.Area

	err := r.db.QueryRow(context.Background(), query, id).Scan(
		&a.ID,
		&a.BuildingID,
		&a.Name,
		&a.Floor,
		&a.CleanlinessClass,
		&a.IsActive,
	)
	if err != nil {
		return nil, err
	}

	return &a, nil
}

// ================= UPDATE =================
func (r *AreaPostgresRepository) Update(area *domain.Area) error {
	query := `
		UPDATE areas
		SET
			building_id = $1,
			name = $2,
			floor = $3,
			cleanliness_class = $4
		WHERE id = $5
	`

	_, err := r.db.Exec(
		context.Background(),
		query,
		area.BuildingID,
		area.Name,
		area.Floor,
		area.CleanlinessClass,
		area.ID,
	)

	return err
}

// ================= DEACTIVATE =================
func (r *AreaPostgresRepository) Deactivate(id string) error {
	_, err := r.db.Exec(
		context.Background(),
		`UPDATE areas SET is_active = false WHERE id = $1`,
		id,
	)
	return err
}
