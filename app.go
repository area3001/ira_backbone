package goira

import (
	"github.com/area3001/goira/comm"
	"github.com/area3001/goira/devices"
	"log"
)

func NewApp(natsClient *comm.NatsClient) *App {
	devs, err := devices.NewModule(natsClient)
	if err != nil {
		log.Panicf("unable to create the devices module: %v", err)
	}

	if err := devs.Start(); err != nil {
		log.Panicln(err.Error())
	}

	return &App{
		Devices: devs,
	}
}

type App struct {
	Devices *devices.Module
}
