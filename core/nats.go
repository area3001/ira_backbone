package core

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"time"
)

type NatsClient struct {
	nc   *nats.Conn
	js   nats.JetStream
	root string
}

func (c *NatsClient) Connect() error {
	// Connect to NATS
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		return fmt.Errorf("unable to connect to nats: %v", err)
	}
	c.nc = nc

	// Create JetStream Context
	js, err := nc.JetStream(nats.PublishAsyncMaxPending(256))
	if err != nil {
		return fmt.Errorf("unable to connect to jetstream: %v", err)
	}
	c.js = js

	js.Subscribe("ORDERS.*", func(m *nats.Msg) {
		fmt.Printf("Received a JetStream message: %s\n", string(m.Data))
	}, nats.DeliverLastPerSubject())

	return nil
}

func (c *NatsClient) IsConnected() bool {
	return c.nc != nil && c.js != nil
}

func (c *NatsClient) Ping() error {
	subject := fmt.Sprintf("%s.ping", c.root)
	return c.nc.Publish(subject, []byte{})
}

func (c *NatsClient) SetMode(mac string, mode *Mode) error {
	if mode == nil {
		return fmt.Errorf("no mode provided")
	}

	subject := fmt.Sprintf("%s.%s.mode", c.root, mac)
	_, err := c.nc.Request(subject, []byte(fmt.Sprintf("%d", mode.Code)), 10*time.Second)
	if err != nil {
		return fmt.Errorf("unable to set mode %q: %v", mode.Name, err)
	}

	return nil
}

func (c *NatsClient) SendBlink(mac string, times int) error {
	subject := fmt.Sprintf("%s.%s.blink", c.root, mac)
	_, err := c.nc.Request(subject, []byte(fmt.Sprintf("%d", times)), 10*time.Second)
	if err != nil {
		return fmt.Errorf("unable to send blink: %v", err)
	}

	return nil
}

func (c *NatsClient) SendReset(mac string, delayMs int) error {
	subject := fmt.Sprintf("%s.%s.reset", mac, c.root)
	return c.nc.Publish(subject, []byte(fmt.Sprintf("%d", delayMs)))
}

func (c *NatsClient) SendDmx(mac string, data []byte) error {
	subject := fmt.Sprintf("%s.%s.dmx", mac, c.root)
	if len(data) > 513 {
		return fmt.Errorf("invalid dmx data length")
	}

	return c.nc.Publish(subject, data)
}

func (c *NatsClient) SendDeltaDmx(mac string, data []byte) error {
	subject := fmt.Sprintf("%s.%s.deltadmx", mac, c.root)

	return c.nc.Publish(subject, data)
}

func (c *NatsClient) SendRgb(mac string, data []byte) error {
	subject := fmt.Sprintf("%s.%s.rgb", mac, c.root)
	return c.nc.Publish(subject, data)
}

func (c *NatsClient) SendFx(mac string, effect *Effect) error {
	subject := fmt.Sprintf("%s.%s.fx", mac, c.root)

	_, err := c.nc.Request(subject, effect.toBytes(), 10*time.Second)
	if err != nil {
		return fmt.Errorf("unable to send effect: %v", err)
	}

	return nil
}
