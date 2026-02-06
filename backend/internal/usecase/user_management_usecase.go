package usecase

import (
	"errors"
	"time"

	"ahu-backend/internal/domain"
	"ahu-backend/internal/repository"
	"ahu-backend/internal/security"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserManagementUsecase struct {
	userRepo  repository.UserRepository
	auditRepo repository.AuditTrailRepository
}

func NewUserManagementUsecase(
	userRepo repository.UserRepository,
	auditRepo repository.AuditTrailRepository,
) *UserManagementUsecase {
	return &UserManagementUsecase{
		userRepo:  userRepo,
		auditRepo: auditRepo,
	}
}

// ================= LIST =================
func (u *UserManagementUsecase) ListUsers() ([]domain.User, error) {
	return u.userRepo.ListAll()
}

// ================= CREATE =================
func (u *UserManagementUsecase) CreateUser(
	npp,
	name,
	jabatan,
	password,
	role,
	adminID string,
) error {
	if npp == "" {
		return errors.New("NPP wajib diisi")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &domain.User{
		ID:           uuid.NewString(),
		NPP:          &npp,
		Name:         name,
		Jabatan:      &jabatan, // sementara pakai kolom email
		PasswordHash: string(hash),
		Role:         role,
		IsActive:     true,
	}

	return u.userRepo.Create(user, adminID)
}

// ================= ACTIVATE =================
func (u *UserManagementUsecase) ActivateUser(
	userID, adminID string,
) error {
	if userID == adminID {
		return errors.New("tidak boleh mengaktifkan diri sendiri")
	}
	return u.userRepo.Activate(userID, adminID)
}

// ================= DEACTIVATE =================
func (u *UserManagementUsecase) DeactivateUser(
	userID, adminID string,
) error {
	if userID == adminID {
		return errors.New("tidak boleh menonaktifkan diri sendiri")
	}
	return u.userRepo.Deactivate(userID, adminID)
}

func (uc *UserManagementUsecase) GetProfile(userID string) (*domain.User, error) {
	return uc.userRepo.GetByID(userID)
}

func (u *UserManagementUsecase) GetUserByID(id string) (*domain.User, error) {
	return u.userRepo.GetByID(id)
}

func (u *UserManagementUsecase) UpdateProfile(
	userID string,
	name string,
) error {
	user, err := u.userRepo.GetByID(userID)
	if err != nil {
		return err
	}

	if err := u.userRepo.UpdateName(userID, name); err != nil {
		return err
	}

	u.auditRepo.Save(&domain.AuditTrail{
		ID:       uuid.NewString(),
		UserID:   userID,
		Action:   "update_profile",
		Entity:   "user",
		EntityID: userID,
		Metadata: map[string]interface{}{
			"old_name": user.Name,
			"new_name": name,
		},
		CreatedAt: time.Now(),
	})

	return nil
}

func (u *UserManagementUsecase) ChangePassword(
	userID, oldPass, newPass string,
) error {
	user, err := u.userRepo.GetByID(userID)
	if err != nil {
		return err
	}

	if err := security.CheckPassword(user.PasswordHash, oldPass); err != nil {
		return errors.New("password lama salah")
	}

	hash, _ := security.HashPassword(newPass)
	return u.userRepo.UpdatePassword(userID, hash)
}

func (uc *UserManagementUsecase) GetMyProfile(userID string) (*domain.User, error) {
	return uc.userRepo.GetByID(userID)
}

func (uc *UserManagementUsecase) ChangeMyPassword(
	userID string,
	newHash string,
) error {
	return uc.userRepo.UpdatePassword(userID, newHash)
}

func (u *UserManagementUsecase) UpdateAvatar(
	userID string,
	avatarURL string,
) error {
	return u.userRepo.UpdateAvatar(userID, avatarURL)
}

func (u *UserManagementUsecase) UpdateBirthDate(
	userID string,
	birthDate time.Time,
) error {
	return u.userRepo.UpdateBirthDate(userID, birthDate)
}

func (u *UserManagementUsecase) UpdateName(
	userID string,
	name string,
) error {
	return u.userRepo.UpdateName(userID, name)
}

func (u *UserManagementUsecase) UpdateUser(
	userID string,
	name *string,
	jabatan *string,
	password *string,
) error {

	if name != nil || jabatan != nil {
		if err := u.userRepo.UpdateBasicProfile(
			userID,
			name,
			jabatan,
		); err != nil {
			return err
		}
	}

	if password != nil && *password != "" {
		hash, _ := security.HashPassword(*password)
		if err := u.userRepo.UpdatePassword(userID, hash); err != nil {
			return err
		}
	}

	return nil
}

func (uc *UserManagementUsecase) UpdateMyProfile(
	userID string,
	name string,
	birth *time.Time,
) error {

	var namePtr *string
	if name != "" {
		namePtr = &name
	}

	return uc.userRepo.UpdatePersonalProfile(
		userID,
		namePtr,
		birth,
	)
}
