package dto

type SubmitInspectionFormRequest struct {
	Items []InspectionFormItemDTO `json:"items"`
	Notes []string                `json:"notes,omitempty"`
}

type InspectionFormItemDTO struct {
	FormItemID  string   `json:"form_item_id"`
	ValueText   *string  `json:"value_text,omitempty"`
	ValueNumber *float64 `json:"value_number,omitempty"`
	ValueBool   *bool    `json:"value_bool,omitempty"`
}
