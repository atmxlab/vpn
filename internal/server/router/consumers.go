package router

import (
	"context"

	"github.com/sirupsen/logrus"
)

func (r *Router) consumeTun(ctx context.Context) error {
	for packet := range r.tunPackets {
		if err := r.tunHandler.Handle(ctx, packet); err != nil {
			logrus.Errorf("failed to handle TUN packet %+v: %v", packet, err)
		}
	}

	return nil
}

func (r *Router) consumeTunnel(ctx context.Context) error {
	for packet := range r.tunnelPackets {
		logrus.Debugf("Readed from tunnel channel: flag=[%s]", packet.Header().Flag())
		handler, ok := r.tunnelHandlerByFlag[packet.Header().Flag()]
		if !ok {
			logrus.Errorf("failed to find handler by flag=[%s]", packet.Header().Flag())
			continue
		}

		if err := handler.Handle(ctx, packet); err != nil {
			logrus.Errorf("failed to handle TUNNEL packet %+v: %v", packet, err)
		}
	}

	return nil
}
