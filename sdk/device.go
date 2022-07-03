package sdk

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"github.com/area3001/goira/comm"
	"github.com/area3001/goira/core"
	"github.com/nats-io/nats.go"
	"image/color"
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

	response, err := d.nc.Request("mode", d.Meta.MAC, []byte(fmt.Sprintf("%d", mode.Code)))
	if err != nil {
		return err
	}

	if string(response) != "+OK" {
		return fmt.Errorf("unexpected response: %s", response)
	}

	return nil
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

func (d *Device) SendRgbPixels(offset uint16, colors []color.RGBA) error {
	buf := toRgbBuffer(offset, colors)

	resp, err := d.nc.Request("rgb", d.Meta.MAC, []byte(base64.StdEncoding.EncodeToString(buf.Bytes())))
	if err != nil {
		return err
	}

	if string(resp) != "+OK" {
		return fmt.Errorf("operation failed: %s", resp)
	}

	return nil
}

func (d *Device) SendFx(effect *core.Effect) error {
	return d.nc.Call("fx", d.Meta.MAC, effect.ToBytes())
}

func toRgbBuffer(offset uint16, colors []color.RGBA) *bytes.Buffer {
	buf := new(bytes.Buffer)

	//of := make([]byte, 4)
	//binary.LittleEndian.PutUint16(of, offset)
	//buf.Write(of)

	binary.Write(buf, binary.BigEndian, offset)
	binary.Write(buf, binary.BigEndian, uint16(len(colors)))

	for _, c := range colors {
		binary.Write(buf, binary.BigEndian, c.R)
		binary.Write(buf, binary.BigEndian, c.G)
		binary.Write(buf, binary.BigEndian, c.B)
	}

	return buf
}
