package main

import (
	"github.com/atmxlab/vpn/pkg/errors"
	"github.com/sirupsen/logrus"
)

func exitIF(err error, msg string, a ...any) {
	if err != nil {
		exit(errors.Wrapf(err, msg, a...))
	}
}

func exit(err error) {
	logrus.Fatal(err.Error())
}
