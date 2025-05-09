package cmd

import (
	"runtime/debug"

	"github.com/atmxlab/vpn/pkg/errors"
	"github.com/sirupsen/logrus"
)

func Exitf(err error, msg string, a ...any) {
	if err != nil {
		Exit(errors.Wrapf(err, msg, a...))
	}
}

func Exit(err error) {
	logrus.Fatal(err.Error())
}

func Recover() {
	if err := recover(); err != nil {
		logrus.Errorf("Panic recovered: %v", err)
		logrus.Fatalf("Stack trace:\n%s", debug.Stack())
	}
}
