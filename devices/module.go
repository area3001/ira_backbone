package devices

import (
	"fmt"
	"github.com/area3001/goira/comm"
	"github.com/labstack/echo/v4"
)

func NewModule(client *comm.NatsClient, e *echo.Group) (*Module, error) {
	devices, err := client.JetStream.KeyValue("devices")
	if err != nil {
		return nil, fmt.Errorf("unable to get the deviced kv store from jetstream: %w", err)
	}

	service := &Service{
		n:      client,
		reader: &Reader{devices: devices},
	}

	RegisterEndpoints(e, service)

	return &Module{
		client:  client,
		Service: service,
		Logic:   &Logic{writer: &Writer{devices: devices}},
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
