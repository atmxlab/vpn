package errors

import "errors"

func As[T error](err error) (target T, ok bool) {
	return target, errors.As(err, &target)
}
