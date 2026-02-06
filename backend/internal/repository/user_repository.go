package repository

import (
	"ahu-backend/internal/domain"
	"time"
)

type UserRepository interface {
	Create(user *domain.User, createdBy string) error
	ListAll() ([]domain.User, error)

	GetByEmail(email string) (*domain.User, error)
	GetByNPP(npp string) (*domain.User, error)
	GetByID(id string) (*domain.User, error)

	Activate(userID, adminID string) error
	Deactivate(userID, adminID string) error

	UpdateName(userID, name string) error
	UpdatePassword(userID, passwordHash string) error
	UpdateBirthDate(userID string, birthDate time.Time) error
	UpdateAvatar(userID string, avatarURL string) error

	// 🔥 ADMIN EDIT USER
	UpdateBasicProfile(
		userID string,
		name *string,
		jabatan *string,
	) error

	// 🔥 USER EDIT PROFILE SENDIRI
	UpdatePersonalProfile(
		userID string,
		name *string,
		birthDate *time.Time,
	) error

	ListInspectors() ([]domain.User, error)
}
