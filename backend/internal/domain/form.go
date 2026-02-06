package domain

import "time"

type FormTemplate struct {
	ID          string        `json:"id"`
	Name        string        `json:"name"`
	Period      string        `json:"period"`
	Description *string       `json:"description"`
	Version     int           `json:"version"`
	IsActive    bool          `json:"is_active"`
	CreatedAt   time.Time     `json:"created_at"`
	Sections    []FormSection `json:"sections"`
}

type FormSection struct {
	ID    string     `json:"id"`
	Code  string     `json:"code"`
	Title string     `json:"title"`
	Order int        `json:"order"`
	Items []FormItem `json:"items"`
}

type FormItem struct {
	ID        string `json:"id"`
	Label     string `json:"label"`
	InputType string `json:"input_type"`
	Required  bool   `json:"required"`
	Order     int    `json:"order"`
}
