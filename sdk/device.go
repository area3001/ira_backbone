package sdk

import (
	"fmt"
	"github.com/area3001/goira/comm"
	"github.com/area3001/goira/core"
	"github.com/nats-io/nats.go"
	"time"
)

type Device struct {
	Meta  *DeviceMeta
	store nats.KeyValue
	nc    *comm.NatsClient
}

type DeviceMeta struct {
	MAC          string    `json:"mac"`
	IP           string    `json:"ip"`
	Hardware     Hardware  `json:"hardware"`
	Mode         string    `json:"mode"`
	ExternalMode string    `json:"external_mode"`
	LastBeat     time.Time `json:"last_beat"`
}

type Hardware struct {
	Kind    string `json:"kind"`
	Version string `json:"revision"`
}

func (d *Device) SetMode(mac string, mode *core.Mode) error {
	if mode == nil {
		return fmt.Errorf("no mode provided")
	}

	return d.nc.Call("mode", mac, []byte(fmt.Sprintf("%d", mode.Code)))
}

func (d *Device) Blink(times int) error {
	return d.nc.Call("blink", d.Meta.MAC, []byte(fmt.Sprintf("%d", times)))
}

func (d *Device) Reset(mac string, delayMs int) error {
	return d.nc.Call("reset", mac, []byte(fmt.Sprintf("%d", delayMs)))
}

func (d *Device) SendDmx(mac string, data []byte) error {
	return d.nc.Call("dmx", mac, data)
}

func (d *Device) SendDeltaDmx(mac string, data []byte) error {
	return d.nc.Call("deltadmx", mac, data)
}

func (d *Device) SendRgb(mac string, data []byte) error {
	return d.nc.Call("rgb", mac, data)
}

func (d *Device) SendFx(mac string, effect *core.Effect) error {
	return d.nc.Call("fx", mac, effect.ToBytes())
}
