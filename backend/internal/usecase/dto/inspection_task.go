package dto

type InspectionTaskDTO struct {
	ScheduleID string  `json:"schedule_id"`
	StartDate  string  `json:"start_date"`
	EndDate    string  `json:"end_date"`
	UnitCode   string  `json:"unit_code"`
	RoomName   *string `json:"room_name"`
	Period     string  `json:"period"`
	Status     string  `json:"status"`
}
