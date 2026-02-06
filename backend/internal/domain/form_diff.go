package domain

type FormTemplateDiff struct {
	From     FormTemplateMeta `json:"from"`
	To       FormTemplateMeta `json:"to"`
	Sections []SectionDiff    `json:"sections"`
}

type FormTemplateMeta struct {
	ID      string `json:"id"`
	Version int    `json:"version"`
}

type SectionDiff struct {
	Code   string     `json:"code"`
	Title  string     `json:"title"`
	Status string     `json:"status"` // added | removed | modified | unchanged
	Items  []ItemDiff `json:"items"`
}

type ItemDiff struct {
	Label  string    `json:"label"`
	Status string    `json:"status"`
	From   *FormItem `json:"from,omitempty"`
	To     *FormItem `json:"to,omitempty"`
}
