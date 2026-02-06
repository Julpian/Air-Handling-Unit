package usecase

import "errors"

func ErrValidation(msg string) error {
	return errors.New(msg)
}
