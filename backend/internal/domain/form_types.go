package domain

const (
	InputTypeText   = "text"
	InputTypeNumber = "number"
	InputTypeSelect = "select"

	InputTypeBooleanClean  = "boolean_clean"
	InputTypeBooleanOK     = "boolean_ok"
	InputTypeBooleanNormal = "boolean_normal"
)

func IsValidInputType(t string) bool {
	switch t {
	case InputTypeText,
		InputTypeNumber,
		InputTypeSelect,
		InputTypeBooleanClean,
		InputTypeBooleanOK,
		InputTypeBooleanNormal:
		return true
	default:
		return false
	}
}