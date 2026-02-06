package domain

import "time"

// User merepresentasikan pengguna sistem
type User struct {
	ID           string     `db:"id"`
	NPP          *string    `db:"npp"`
	Name         string     `db:"name"`
	Jabatan      *string    `db:"jabatan"`
	PasswordHash string     `db:"password_hash"`
	Role         string     `db:"role"`
	IsActive     bool       `db:"is_active"`
	CreatedBy    *string    `db:"created_by"` // admin yang membuat user
	CreatedAt    time.Time  `db:"created_at"`
	UpdatedAt    *time.Time `db:"updated_at"`
	BirthDate    *time.Time
	AvatarURL    *string `db:"avatar_url"`
}
