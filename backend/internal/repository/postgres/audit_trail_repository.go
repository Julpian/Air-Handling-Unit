package postgres

import (
	"context"
	"time"

	"ahu-backend/internal/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuditTrailPostgresRepository struct {
	db *pgxpool.Pool
}

func NewAuditTrailPostgresRepository(db *pgxpool.Pool) *AuditTrailPostgresRepository {
	return &AuditTrailPostgresRepository{db: db}
}

// ================= SAVE =================
func (r *AuditTrailPostgresRepository) Save(audit *domain.AuditTrail) error {

	loc, _ := time.LoadLocation("Asia/Jakarta")

	query := `
		INSERT INTO audit_trails (
			id, user_id, action, entity, entity_id, metadata, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.db.Exec(
		context.Background(),
		query,
		uuid.NewString(),
		audit.UserID,
		audit.Action,
		audit.Entity,
		audit.EntityID,
		audit.Metadata,
		time.Now().In(loc), // ✅ waktu dari Go, bukan dari DB
	)

	return err
}

// ================= LIST ALL (ADMIN) =================
func (r *AuditTrailPostgresRepository) ListAll(
	limit int,
) ([]domain.AuditTrailView, error) {

	query := `
	SELECT
		a.id,
		a.created_at,
		u.id,
		u.name,
		u.role,
		a.action,
		a.entity,
		a.entity_id,
		a.metadata
	FROM audit_trails a
	JOIN users u ON u.id = a.user_id
	ORDER BY a.created_at DESC
	LIMIT $1
	`

	rows, err := r.db.Query(context.Background(), query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []domain.AuditTrailView

	for rows.Next() {
		var a domain.AuditTrailView

		if err := rows.Scan(
			&a.ID,
			&a.CreatedAt,
			&a.UserID,
			&a.Name,
			&a.Role,
			&a.Action,
			&a.Entity,
			&a.EntityID,
			&a.Metadata,
		); err != nil {
			return nil, err
		}

		list = append(list, a)
	}

	return list, nil
}

// ================= LIST BY ENTITY =================
func (r *AuditTrailPostgresRepository) ListByEntity(
	entity string,
	entityID string,
) ([]domain.AuditTrailView, error) {

	query := `
	SELECT
		a.id,
		a.created_at,
		u.id,
		u.name,
		u.role,
		a.action,
		a.entity,
		a.entity_id,
		a.metadata
	FROM audit_trails a
	JOIN users u ON u.id = a.user_id
	WHERE a.entity = $1 AND a.entity_id = $2
	ORDER BY a.created_at ASC
	`

	rows, err := r.db.Query(context.Background(), query, entity, entityID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []domain.AuditTrailView

	for rows.Next() {
		var a domain.AuditTrailView

		if err := rows.Scan(
			&a.ID,
			&a.CreatedAt,
			&a.UserID,
			&a.Name,
			&a.Role,
			&a.Action,
			&a.Entity,
			&a.EntityID,
			&a.Metadata,
		); err != nil {
			return nil, err
		}

		list = append(list, a)
	}

	return list, nil
}