package devices

import (
	"fmt"
	"github.com/area3001/goira/comm"
	"github.com/nats-io/nats.go"
	"log"
)

type Logic struct {
	writer *Writer
}

func (l *Logic) Connect(client *comm.NatsClient) error {
	return client.HandleBroadcast("announce", func(msg *nats.Msg) {
		if err := l.handleAnnounce(msg); err != nil {
			errMsg := fmt.Sprintf("handling device announcement failed: %v", err)
			log.Println(errMsg)
		}
	})
}

func (l *Logic) handleAnnounce(msg *nats.Msg) error {
	dev, err := ParseDevice(msg)
	if err != nil {
		return fmt.Errorf("unable to parse device data: %w", err)
	}

	if err := l.writer.Register(dev); err != nil {
		return fmt.Errorf("unable to register device %s: %w", dev.MAC, err)
	}

	return nil
}
