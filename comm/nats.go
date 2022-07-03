package comm

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"log"
	"time"
)

type HandlerContext struct {
	JetStream nats.JetStreamContext
}

type NatsClientOpts struct {
	Root             string
	NatsUrl          string
	NatsOptions      []nats.Option
	JetStreamOptions []nats.JSOpt
}

func Dial(opts *NatsClientOpts) (*NatsClient, error) {
	result := &NatsClient{
		root: opts.Root,
	}

	// Connect to NATS
	nc, err := nats.Connect(opts.NatsUrl, opts.NatsOptions...)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to comm: %v", err)
	}
	result.nc = nc

	// Create JetStream Context
	js, err := nc.JetStream(opts.JetStreamOptions...)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to jetstream: %v", err)
	}
	result.JetStream = js

	result.subscriptions = map[string]*nats.Subscription{}

	return result, nil
}

type NatsClient struct {
	nc        *nats.Conn
	JetStream nats.JetStreamContext
	root      string

	subscriptions map[string]*nats.Subscription
}

func (c *NatsClient) Close() error {
	return c.nc.Drain()
}

func (c *NatsClient) IsConnected() bool {
	return c.nc != nil && c.JetStream != nil
}

func (c *NatsClient) Broadcast(channel string) error {
	subject := fmt.Sprintf("%s.%s", c.root, channel)

	if err := c.nc.Publish(subject, []byte{}); err != nil {
		return err
	}

	return c.nc.Flush()
}

func (c *NatsClient) Request(channel string, device string, data []byte) ([]byte, error) {
	subject := fmt.Sprintf("%s.%s.%s", c.root, device, channel)
	response, err := c.nc.Request(subject, data, 10*time.Second)
	if err != nil {
		return nil, fmt.Errorf("calling channel %s for device %s failed: %v", channel, device, err)
	}

	if response != nil {
		return response.Data, nil
	}

	return nil, nil
}

func (c *NatsClient) Call(channel string, device string, data []byte) error {
	subject := fmt.Sprintf("%s.%s.%s", c.root, device, channel)
	_, err := c.nc.Request(subject, data, 10*time.Second)
	if err != nil {
		return fmt.Errorf("calling channel %s for device %s failed: %v", channel, device, err)
	}

	return nil
}

func (c *NatsClient) HandleBroadcast(channel string, handler nats.MsgHandler) error {
	subject := fmt.Sprintf("%s.%s", c.root, channel)

	if _, fnd := c.subscriptions[subject]; fnd {
		return fmt.Errorf("a broadcast listener for %s is already registered", subject)
	}

	log.Println("listening for broadcast messages on " + subject)
	s, err := c.JetStream.QueueSubscribe(subject, fmt.Sprintf("%s-broadcast-handlers", channel), handler)

	if err != nil {
		return err
	}
	c.subscriptions[subject] = s

	return nil
}
