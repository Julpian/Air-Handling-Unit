package usecase

import (
	"errors"

	"ahu-backend/internal/domain"
	"ahu-backend/internal/repository"
	"ahu-backend/internal/security"
)

type AuthUsecase struct {
	userRepo repository.UserRepository
	auditUC  *AuditTrailUsecase // 👈 Tambahkan ini agar bisa log
}

// Update Constructor agar menerima auditUC
func NewAuthUsecase(userRepo repository.UserRepository, auditUC *AuditTrailUsecase) *AuthUsecase {
	return &AuthUsecase{
		userRepo: userRepo,
		auditUC:  auditUC,
	}
}

// 🔐 LOGIN DENGAN NPP
func (u *AuthUsecase) Login(npp, password string) (string, string, error) {
	user, err := u.userRepo.GetByNPP(npp)
	
	// Cek jika user tidak ditemukan
	if err != nil {
		u.auditUC.Log(&domain.AuditTrail{
			Action:   "LOGIN_FAILED",
			Entity:   "auth",
			EntityID: npp,
			Metadata: map[string]interface{}{"reason": "npp_not_found"},
		})
		return "", "", errors.New("NPP atau password salah")
	}

	if !user.IsActive {
		u.auditUC.Log(&domain.AuditTrail{
			UserID:   user.ID,
			Action:   "LOGIN_FAILED",
			Entity:   "auth",
			EntityID: user.ID,
			Metadata: map[string]interface{}{"reason": "account_inactive"},
		})
		return "", "", errors.New("akun dinonaktifkan")
	}

	// Cek Password (menggunakan package security kamu)
	if err := security.CheckPassword(user.PasswordHash, password); err != nil {
		u.auditUC.Log(&domain.AuditTrail{
			UserID:   user.ID,
			Action:   "LOGIN_FAILED",
			Entity:   "auth",
			EntityID: user.ID,
			Metadata: map[string]interface{}{"reason": "wrong_password"},
		})
		return "", "", errors.New("NPP atau password salah")
	}

	// Generate Token
	token, err := security.GenerateToken(
		user.ID,
		user.Role,
		user.IsActive,
	)
	if err != nil {
		return "", "", err
	}

	// ✅ CATAT LOGIN BERHASIL
	u.auditUC.Log(&domain.AuditTrail{
		UserID:   user.ID,
		Action:   "LOGIN_SUCCESS",
		Entity:   "auth",
		EntityID: user.ID,
		Metadata: map[string]interface{}{
			"role": user.Role,
			"name": user.Name,
		},
	})

	return token, user.Role, nil
}

func (u *AuthUsecase) Logout(userID string, name string) {
    u.auditUC.Log(&domain.AuditTrail{
        UserID:   userID,
        Action:   "LOGOUT",
        Entity:   "auth",
        EntityID: userID,
        Metadata: map[string]interface{}{
            "name":   name,
            "status": "success",
        },
    })
}