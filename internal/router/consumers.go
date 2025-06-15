package router

import (
	"context"

	"github.com/sirupsen/logrus"
)

func (r *Router) consumeTun(ctx context.Context) error {
	log := logrus.WithField("Namespace", "TUN")

	for packet := range r.tunPackets {
		log = log.WithField("Len", len(packet.Payload()))
		log.Debug("Read packet")

		if err := r.tunHandler.Handle(ctx, packet); err != nil {
			log.
				WithField("Len", len(packet.Payload())).
				WithError(err).
				Warn("Failed to handle packet")
		}
	}

	return nil
}

func (r *Router) consumeTunnel(ctx context.Context) error {
	log := logrus.WithField("Namespace", "TUNNEL")

	for packet := range r.tunnelPackets {
		log = log.
			WithField("Namespace", "TUNNEL").
			WithField("Flag", packet.Header().Flag()).
			WithField("Len", len(packet.Payload()))

		log.Debug("Read packet")

		handler, ok := r.tunnelHandlerByFlag[packet.Header().Flag()]
		if !ok {
			log.Errorf("Failed to find handler")
			continue
		}

		if err := handler.Handle(ctx, packet); err != nil {
			log.
				WithError(err).
				Warn("Failed to handle packet")
		}
	}

	return nil
}
