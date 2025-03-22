package errors

import (
	"fmt"
)

func Wrap(err error, format string, a ...any) error {
	return fmt.Errorf("%s: %w", fmt.Sprintf(format, a...), err)
}

func Fatalf(format string, a ...any) {
	panic(fmt.Sprintf(format, a...))
}
