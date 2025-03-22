package errors

func ValidateErr(msg string) error {
	return Wrap(ErrInvalidArgument, msg)
}
