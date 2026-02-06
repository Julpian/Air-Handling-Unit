package usecase

import (
	"errors"

	"ahu-backend/internal/repository"
	"ahu-backend/internal/security"
)

type AuthUsecase struct {
	userRepo repository.UserRepository
}

func NewAuthUsecase(userRepo repository.UserRepository) *AuthUsecase {
	return &AuthUsecase{userRepo: userRepo}
}

// 🔐 LOGIN DENGAN NPP
func (u *AuthUsecase) Login(npp, password string) (string, string, error) {
	user, err := u.userRepo.GetByNPP(npp)
	if err != nil {
		return "", "", errors.New("NPP atau password salah")
	}

	if !user.IsActive {
		return "", "", errors.New("akun dinonaktifkan")
	}

	if err := security.CheckPassword(user.PasswordHash, password); err != nil {
		return "", "", errors.New("NPP atau password salah")
	}

	token, err := security.GenerateToken(
		user.ID,
		user.Role,
		user.IsActive,
	)
	if err != nil {
		return "", "", err
	}

	return token, user.Role, nil
}
