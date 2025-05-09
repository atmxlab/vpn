package cmd

import (
	"runtime/debug"

	"github.com/sirupsen/logrus"
)

func Recover() {
	if err := recover(); err != nil {
		logrus.Errorf("Panic recovered: %v", err)
		logrus.Fatalf("Stack trace:\n%s", debug.Stack())
	}
}
