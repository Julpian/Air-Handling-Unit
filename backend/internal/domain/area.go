package domain

type Area struct {
	ID               string  `json:"id"`
	Name             string  `json:"name"`
	BuildingID       string  `json:"building_id"`
	Floor            *string `json:"floor"`             // ⬅ pointer
	CleanlinessClass *string `json:"cleanliness_class"` // ⬅ pointer
	IsActive         bool    `json:"is_active"`
}
