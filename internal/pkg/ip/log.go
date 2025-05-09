package ip

import (
	"github.com/sirupsen/logrus"
	"golang.org/x/net/ipv4"
)

func LogHeader(frame []byte) {
	header, err := ipv4.ParseHeader(frame)
	if err != nil {
		logrus.
			WithError(err).
			Warn("Failed to parse IP header")

		return
	}

	logrus.
		WithField("SRC", header.Src).
		WithField("DST", header.Dst).
		WithField("VERSION", header.Version).
		WithField("PROTOCOL", header.Protocol).
		WithField("CHECKSUM", header.Checksum).
		WithField("LEN", header.Len).
		Debugf("%s -> %s", header.Src, header.Dst)
}
