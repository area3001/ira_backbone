package sdk

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
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

func (d *Device) SetMode(mode *core.Mode) error {
	if mode == nil {
		return fmt.Errorf("no mode provided")
	}

	return d.nc.Call("mode", d.Meta.MAC, []byte(fmt.Sprintf("%d", mode.Code)))
}

func (d *Device) SetConfig(param string, value string) error {
	buf := new(bytes.Buffer)

	b, fnd := ConfigParams[param]
	if !fnd {
		return fmt.Errorf("%s: invalid parameter", param)
	}
	buf.Write(b)
	buf.WriteString(value)

	return d.nc.Call("config", d.Meta.MAC, []byte(base64.StdEncoding.EncodeToString(buf.Bytes())))
}

func (d *Device) Blink(times int) error {
	return d.nc.Call("blink", d.Meta.MAC, []byte(fmt.Sprintf("%d", times)))
}

func (d *Device) Reset(delayMs int) error {
	return d.nc.Call("reset", d.Meta.MAC, []byte(fmt.Sprintf("%d", delayMs)))
}

func (d *Device) SendDmx(data []byte) error {
	return d.nc.Call("dmx", d.Meta.MAC, data)
}

func (d *Device) SendDeltaDmx(data []byte) error {
	return d.nc.Call("deltadmx", d.Meta.MAC, data)
}

func (d *Device) SendRgbRaw(data []byte) error {
	return d.nc.Call("rgb", d.Meta.MAC, []byte(base64.StdEncoding.EncodeToString(data)))
}

func (d *Device) SendRgbPixels(offset int, hexCodes []string) error {
	// -- construct the packet
	bytesPerPixel := len(hexCodes[0])
	if bytesPerPixel != 6 && bytesPerPixel != 8 {
		return fmt.Errorf("wrong number of pixel channels")
	}

	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, uint16(offset))
	binary.Write(buf, binary.LittleEndian, uint16(len(hexCodes)))

	for idx, arg := range hexCodes {
		b, err := hex.DecodeString(arg)
		if err != nil {
			return fmt.Errorf("data %s for pixel %d is invalid: %s", arg, idx, err)
		}

		if len(b) != bytesPerPixel {
			return fmt.Errorf("data %s for pixel %d is invalid: inconsistent number of bytes per pixel", arg, idx)
		}

		buf.Write(b)
	}

	return d.nc.Call("rgb", d.Meta.MAC, []byte(base64.StdEncoding.EncodeToString(buf.Bytes())))
}

func (d *Device) SendFx(effect *core.Effect) error {
	return d.nc.Call("fx", d.Meta.MAC, effect.ToBytes())
}
