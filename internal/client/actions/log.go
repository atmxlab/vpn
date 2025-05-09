package actions

import (
	"github.com/sirupsen/logrus"
)

func log(action string) *logrus.Entry {
	return logrus.
		WithField("Namespace", "TUNNEL|ACTION").
		WithField("Action", action)
}
