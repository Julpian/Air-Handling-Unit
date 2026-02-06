package repository

import "ahu-backend/internal/domain"

type InspectionResultRepository interface {
	SaveMany(results []domain.InspectionResult) error
}
