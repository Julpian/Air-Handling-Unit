package postgres

import (
	"context"
	"time"

	"ahu-backend/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AHUPostgresRepository struct {
	db *pgxpool.Pool
}

func NewAHUPostgresRepository(db *pgxpool.Pool) *AHUPostgresRepository {
	return &AHUPostgresRepository{db: db}
}

func (r *AHUPostgresRepository) Create(ahu *domain.AHU) error {
	query := `
		INSERT INTO ahus
		(id, area_id, unit_code, room_name, vendor, nfc_uid, cleanliness_class, is_active, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,true,$8)
	`
	_, err := r.db.Exec(
		context.Background(),
		query,
		ahu.ID,
		ahu.AreaID,
		ahu.UnitCode,
		ahu.RoomName,
		ahu.Vendor,
		ahu.NFCUID,
		ahu.CleanlinessClass,
		time.Now(),
	)
	return err
}

func (r *AHUPostgresRepository) ListAll() ([]domain.AHU, error) {
	query := `
		SELECT
		id,
		area_id,
		unit_code,
		room_name,
		vendor,
		nfc_uid,
		cleanliness_class,
		is_active,
		created_at
		FROM ahus
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []domain.AHU
	for rows.Next() {
		var a domain.AHU
		if err := rows.Scan(
			&a.ID,
			&a.AreaID,
			&a.UnitCode,
			&a.RoomName,
			&a.Vendor,
			&a.NFCUID,
			&a.CleanlinessClass,
			&a.IsActive,
			&a.CreatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, a)
	}
	return list, nil
}

func (r *AHUPostgresRepository) GetByID(id string) (*domain.AHU, error) {
	query := `
		SELECT
		id,
		area_id,
		unit_code,
		room_name,
		vendor,
		nfc_uid,
		cleanliness_class,
		is_active,
		created_at
		FROM ahus WHERE id=$1
	`
	var a domain.AHU
	err := r.db.QueryRow(context.Background(), query, id).Scan(
		&a.ID,
		&a.AreaID,
		&a.UnitCode,
		&a.RoomName,
		&a.Vendor,
		&a.NFCUID,
		&a.CleanlinessClass,
		&a.IsActive,
		&a.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *AHUPostgresRepository) Update(ahu *domain.AHU) error {
	query := `
		UPDATE ahus
		SET
		area_id=$1,
		unit_code=$2,
		room_name=$3,
		vendor=$4,
		nfc_uid=$5,
		cleanliness_class=$6,
		is_active=$7
		WHERE id=$8
	`
	_, err := r.db.Exec(
		context.Background(),
		query,
		ahu.AreaID,
		ahu.UnitCode,
		ahu.RoomName,
		ahu.Vendor,
		ahu.NFCUID,
		ahu.CleanlinessClass,
		ahu.IsActive,
		ahu.ID,
	)
	return err
}

func (r *AHUPostgresRepository) Deactivate(id string) error {
	_, err := r.db.Exec(
		context.Background(),
		`UPDATE ahus SET is_active=false WHERE id=$1`,
		id,
	)
	return err
}

func (r *AHUPostgresRepository) GetByNFCUID(nfcUID string) (*domain.AHU, error) {
	query := `
		SELECT
			id,
			area_id,
			unit_code,
			room_name,
			vendor,
			nfc_uid,
			cleanliness_class,
			is_active,
			created_at
		FROM ahus
		WHERE LOWER(TRIM(nfc_uid)) = LOWER(TRIM($1))
		AND is_active = true
		LIMIT 1
	`

	var a domain.AHU
	err := r.db.QueryRow(context.Background(), query, nfcUID).Scan(
		&a.ID,
		&a.AreaID,
		&a.UnitCode,
		&a.RoomName,
		&a.Vendor,
		&a.NFCUID,
		&a.CleanlinessClass,
		&a.IsActive,
		&a.CreatedAt,
	)

	if err != nil {
		return nil, nil // NFC tidak ditemukan
	}

	return &a, nil
}

func (r *AHUPostgresRepository) ListAllOrderedByCreatedAt() ([]domain.AHU, error) {
	query := `
		SELECT id, area_id, unit_code, vendor, nfc_uid, is_active, created_at
		FROM ahus
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []domain.AHU
	for rows.Next() {
		var a domain.AHU
		if err := rows.Scan(
			&a.ID,
			&a.AreaID,
			&a.UnitCode,
			&a.Vendor,
			&a.NFCUID,
			&a.IsActive,
			&a.CreatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, a)
	}

	return list, nil
}
