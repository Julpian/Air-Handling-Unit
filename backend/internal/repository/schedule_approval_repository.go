package repository

import "ahu-backend/internal/domain"

type ScheduleApprovalRepository interface {
	GetByYear(year int) (*domain.ScheduleApproval, error)
	Create(year int) error

	SignSVP(year int, userID string, signature string) error
	SignAsmen(year int, userID string, signature string) error

	SetPDF(year int, path string) error
	SetVerifyToken(year int, token string) error
	GetByToken(token string) (*domain.ScheduleApproval, error)
}
