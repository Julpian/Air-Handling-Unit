package domain

import "time"

type AHU struct {
	ID               string    `json:"id"`
	BuildingID       string    `json:"building_id"`
	AreaID           string    `json:"area_id"`
	UnitCode         string    `json:"unit_code"`
	RoomName         *string   `json:"room_name,omitempty"`
	Vendor           *string   `json:"vendor,omitempty"`
	NFCUID           *string   `json:"nfc_uid,omitempty"`
	CleanlinessClass *string   `json:"cleanliness_class,omitempty"`
	IsActive         bool      `json:"is_active"`
	CreatedAt        time.Time `json:"created_at"`
}
