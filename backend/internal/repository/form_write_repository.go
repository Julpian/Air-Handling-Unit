package repository

import (
	"context"
)

type CreateFormTemplateRequest struct {
	Name        string
	Period      string
	Description string

	Sections []struct {
		Code  string
		Title string
		Order int
		Items []struct {
			Label    string
			Type     string
			Required bool
			Options  any
			Order    int
		}
	}
}

type FormWriteRepository interface {
	CreateFormTemplate(
		ctx context.Context,
		req CreateFormTemplateRequest,
	) error
}
