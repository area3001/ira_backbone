package sdk

import (
	"fmt"
	"github.com/area3001/goira/comm"
)

func NewClient(opts *comm.NatsClientOpts) (*Client, error) {
	nc, err := comm.Dial(opts)
	if err != nil {
		return nil, err
	}

	deviceStore, err := nc.JetStream.KeyValue("devices")
	if err != nil {
		return nil, fmt.Errorf("unable to get the deviced kv store from jetstream: %w", err)
	}

	return &Client{nc, &Devices{deviceStore, nc}}, nil
}

type Client struct {
	nc      *comm.NatsClient
	Devices *Devices
}
