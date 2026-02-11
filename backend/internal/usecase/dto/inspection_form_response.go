package dto

import "ahu-backend/internal/domain"

type InspectionFormResponse struct {
	ID       string               `json:"id"`
	Name     string               `json:"name"`
	AHUName  string               `json:"ahu_name"`
	Sections []domain.FormSection `json:"sections"`
}
