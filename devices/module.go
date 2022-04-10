package devices

import (
	"fmt"
	"github.com/area3001/goira/comm"
)

func NewModule(client *comm.NatsClient) (*Module, error) {
	devices, err := client.JetStream.KeyValue("devices")
	if err != nil {
		return nil, fmt.Errorf("unable to get the deviced kv store from jetstream: %w", err)
	}

	return &Module{
		client: client,
		Service: &Service{
			n:      client,
			reader: &Reader{devices: devices},
		},
		Logic: &Logic{writer: &Writer{devices: devices}},
	}, nil
}

type Module struct {
	client  *comm.NatsClient
	Service *Service
	Logic   *Logic
}

func (m *Module) Start() error {
	return m.Logic.Connect(m.client)
}
