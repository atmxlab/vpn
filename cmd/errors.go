package cmd

import (
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
