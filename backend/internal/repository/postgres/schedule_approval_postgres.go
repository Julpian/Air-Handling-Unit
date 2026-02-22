package postgres

import (
	"context"
	"time"

	"ahu-backend/internal/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ScheduleApprovalPostgres struct {
	db *pgxpool.Pool
}

func NewScheduleApprovalPostgres(db *pgxpool.Pool) *ScheduleApprovalPostgres {
	return &ScheduleApprovalPostgres{db: db}
}

func (r *ScheduleApprovalPostgres) GetByYear(year int) (*domain.ScheduleApproval, error) {
	row := r.db.QueryRow(context.Background(), `
		SELECT id,year,svp_id,svp_signed_at,svp_signature,
		       asmen_id,asmen_signed_at,asmen_signature,
		       pdf_path,status,verify_token,created_at
		FROM schedule_approvals WHERE year=$1
	`, year)

	var s domain.ScheduleApproval

	err := row.Scan(
		&s.ID,
		&s.Year,
		&s.SVPID,
		&s.SVPSignedAt,
		&s.SVPSignature,
		&s.AsmenID,
		&s.AsmenSignedAt,
		&s.AsmenSignature,
		&s.PDFPath,
		&s.Status,
		&s.VerifyToken,
		&s.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &s, nil
}

func (r *ScheduleApprovalPostgres) Create(year int) error {

	token := uuid.NewString()

	_, err := r.db.Exec(context.Background(), `
		INSERT INTO schedule_approvals (id,year,verify_token)
		VALUES (gen_random_uuid(),$1,$2)
	`, year, token)

	return err
}

func (r *ScheduleApprovalPostgres) SignSVP(year int, userID, signature string) error {
	_, err := r.db.Exec(context.Background(), `
		UPDATE schedule_approvals
		SET svp_id=$1,
		    svp_signature=$2,
		    svp_signed_at=$3,
		    status='signed_by_svp'
		WHERE year=$4
	`, userID, signature, time.Now(), year)

	return err
}

func (r *ScheduleApprovalPostgres) SignAsmen(year int, userID string, signature string) error {

	query := `
        UPDATE schedule_approvals
        SET
            asmen_id = $1,
            asmen_signature = $2,
            asmen_signed_at = now(),
            status = 'completed'
        WHERE year = $3
    `

	_, err := r.db.Exec(context.Background(), query,
		userID,
		signature,
		year,
	)

	return err
}

func (r *ScheduleApprovalPostgres) SetPDF(year int, path string) error {
	_, err := r.db.Exec(context.Background(), `
		UPDATE schedule_approvals SET pdf_path=$1 WHERE year=$2
	`, path, year)

	return err
}

func (r *ScheduleApprovalPostgres) SetVerifyToken(year int, token string) error {
	_, err := r.db.Exec(context.Background(), `
		UPDATE schedule_approvals
		SET verify_token=$1
		WHERE year=$2
	`, token, year)

	return err
}

func (r *ScheduleApprovalPostgres) GetByToken(token string) (*domain.ScheduleApproval, error) {
	row := r.db.QueryRow(context.Background(), `
		SELECT id,year,svp_id,svp_signed_at,svp_signature,
		       asmen_id,asmen_signed_at,asmen_signature,
		       pdf_path,status,verify_token,created_at
		FROM schedule_approvals
		WHERE verify_token=$1
	`, token)

	var s domain.ScheduleApproval

	err := row.Scan(
		&s.ID,
		&s.Year,
		&s.SVPID,
		&s.SVPSignedAt,
		&s.SVPSignature,
		&s.AsmenID,
		&s.AsmenSignedAt,
		&s.AsmenSignature,
		&s.PDFPath,
		&s.Status,
		&s.VerifyToken,
		&s.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &s, nil
}
