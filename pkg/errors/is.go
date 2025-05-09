package errors

import "errors"

func Is(err error, target error) bool {
	return errors.Is(err, target)
}

func IsSomeBut(err error, references ...error) bool {
	return err != nil && !IsAny(err, references...)
}

func IsAny(err error, references ...error) bool {
	for _, ref := range references {
		if errors.Is(err, ref) {
			return true
		}
	}

	return false
}
