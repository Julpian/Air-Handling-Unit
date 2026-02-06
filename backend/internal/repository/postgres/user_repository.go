package postgres

import (
	"context"
	"errors"
	"time"

	"ahu-backend/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserPostgresRepository struct {
	db *pgxpool.Pool
}

func NewUserPostgresRepository(db *pgxpool.Pool) *UserPostgresRepository {
	return &UserPostgresRepository{db: db}
}

func (r *UserPostgresRepository) GetByEmail(email string) (*domain.User, error) {
	// email DIANGGAP SEBAGAI NPP
	row := r.db.QueryRow(context.Background(), `
		SELECT id, npp, name, jabatan, password_hash, role, is_active
		FROM users
		WHERE npp = $1
	`, email)

	var user domain.User
	err := row.Scan(
		&user.ID,
		&user.NPP,
		&user.Name,
		&user.Jabatan,
		&user.PasswordHash,
		&user.Role,
		&user.IsActive,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserPostgresRepository) Create(
	user *domain.User,
	createdBy string,
) error {
	_, err := r.db.Exec(context.Background(), `
		INSERT INTO users (
			id,
			npp,
			name,
			jabatan,
			password_hash,
			role,
			is_active,
			created_by
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
	`,
		user.ID,
		*user.NPP, // ⬅️ PENTING
		user.Name,
		user.Jabatan,
		user.PasswordHash,
		user.Role,
		user.IsActive,
		createdBy,
	)

	return err
}

func (r *UserPostgresRepository) SetActive(id string, active bool) error {
	_, err := r.db.Exec(
		context.Background(),
		`UPDATE users SET is_active=$1 WHERE id=$2`,
		active,
		id,
	)
	return err
}

func (r *UserPostgresRepository) GetByID(id string) (*domain.User, error) {
	row := r.db.QueryRow(context.Background(), `
		SELECT 
			id,
			npp,
			name,
			jabatan,
			role,
			birth_date,
			avatar_url
		FROM users
		WHERE id = $1
	`, id)

	var u domain.User
	err := row.Scan(
		&u.ID,
		&u.NPP,
		&u.Name,
		&u.Jabatan,
		&u.Role,
		&u.BirthDate,
		&u.AvatarURL,
	)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (r *UserPostgresRepository) UpdateName(userID, name string) error {
	_, err := r.db.Exec(
		context.Background(),
		`UPDATE users SET name = $1 WHERE id = $2`,
		name,
		userID,
	)
	return err
}

func (r *UserPostgresRepository) UpdatePassword(userID, passwordHash string) error {
	_, err := r.db.Exec(
		context.Background(),
		`UPDATE users SET password_hash = $1 WHERE id = $2`,
		passwordHash,
		userID,
	)
	return err
}

func (r *UserPostgresRepository) ListAll() ([]domain.User, error) {
	rows, err := r.db.Query(context.Background(), `
		SELECT id, npp, name, jabatan, role, is_active
		FROM users
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var u domain.User
		if err := rows.Scan(
			&u.ID,
			&u.NPP,
			&u.Name,
			&u.Jabatan,
			&u.Role,
			&u.IsActive,
		); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}

func (r *UserPostgresRepository) DeactivateUser(id string) error {
	_, err := r.db.Exec(
		context.Background(),
		`UPDATE users SET is_active = false WHERE id = $1`,
		id,
	)
	return err
}

func (r *UserPostgresRepository) ActivateUser(id string) error {
	_, err := r.db.Exec(
		context.Background(),
		`UPDATE users SET is_active = true WHERE id = $1`,
		id,
	)
	return err
}

func (r *UserPostgresRepository) Activate(
	userID string,
	adminID string,
) error {
	cmd, err := r.db.Exec(context.Background(), `
		UPDATE users
		SET is_active = true
		WHERE id = $1 AND is_active = false
	`, userID)

	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return errors.New("user tidak ditemukan atau sudah aktif")
	}

	return nil
}

func (r *UserPostgresRepository) Deactivate(
	userID string,
	adminID string,
) error {
	cmd, err := r.db.Exec(context.Background(), `
		UPDATE users
		SET is_active = false
		WHERE id = $1 AND is_active = true
	`, userID)

	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return errors.New("user tidak ditemukan atau sudah nonaktif")
	}

	return nil
}

func (r *UserPostgresRepository) UpdateProfile(
	userID string,
	name *string,
	jabatan *string,
) error {
	_, err := r.db.Exec(context.Background(), `
		UPDATE users SET
			name = COALESCE($1, name),
			jabatan = COALESCE($2, jabatan),
			updated_at = NOW()
		WHERE id = $3
	`, name, jabatan, userID)

	return err
}

func (r *UserPostgresRepository) UpdateBasicProfile(
	userID string,
	name *string,
	jabatan *string,
) error {
	_, err := r.db.Exec(context.Background(), `
		UPDATE users SET
			name = COALESCE($1, name),
			jabatan = COALESCE($2, jabatan),
			updated_at = NOW()
		WHERE id = $3
	`, name, jabatan, userID)

	return err
}

func (r *UserPostgresRepository) UpdatePersonalProfile(
	userID string,
	name *string,
	birthDate *time.Time,
) error {
	_, err := r.db.Exec(context.Background(), `
		UPDATE users SET
			name = COALESCE($1, name),
			birth_date = COALESCE($2, birth_date),
			updated_at = NOW()
		WHERE id = $3
	`, name, birthDate, userID)

	return err
}

func (r *UserPostgresRepository) UpdateAvatar(
	userID string,
	avatarURL string,
) error {
	_, err := r.db.Exec(
		context.Background(),
		`
		UPDATE users
		SET avatar_url = $1, updated_at = NOW()
		WHERE id = $2
		`,
		avatarURL,
		userID,
	)
	return err
}

func (r *UserPostgresRepository) UpdateBirthDate(
	userID string,
	birthDate time.Time,
) error {
	_, err := r.db.Exec(
		context.Background(),
		`
		UPDATE users
		SET birth_date = $1, updated_at = NOW()
		WHERE id = $2
		`,
		birthDate,
		userID,
	)
	return err
}

func (r *UserPostgresRepository) ListInspectors() ([]domain.User, error) {
	query := `
		SELECT id, name, role
		FROM users
		WHERE role = 'inspector'
		ORDER BY name ASC
	`

	rows, err := r.db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []domain.User

	for rows.Next() {
		var u domain.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Role); err != nil {
			return nil, err
		}
		result = append(result, u)
	}

	return result, nil
}

func (r *UserPostgresRepository) GetByNPP(npp string) (*domain.User, error) {
	row := r.db.QueryRow(context.Background(), `
		SELECT id, npp, password_hash, role, is_active
		FROM users
		WHERE npp = $1
	`, npp)

	var user domain.User
	err := row.Scan(
		&user.ID,
		&user.NPP,
		&user.PasswordHash,
		&user.Role,
		&user.IsActive,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
