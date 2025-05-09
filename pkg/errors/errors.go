package errors

import "errors"

var (
	ErrInvalidArgument = errors.New("invalid argument")
	ErrNotFound        = errors.New("not found")
	ErrAlreadyExists   = errors.New("already exists")
)

func NotFoundf(msg string, a ...any) error {
	return Wrapf(ErrNotFound, msg, a)
}

func NotFound(msg string) error {
	return Wrap(ErrNotFound, msg)
}

func InvalidArgument(msg string, a ...any) error {
	return Wrapf(ErrInvalidArgument, msg, a)
}

func AlreadyExists(msg string, a ...any) error {
	return Wrapf(ErrAlreadyExists, msg, a)
}
