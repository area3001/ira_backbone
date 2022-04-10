package devices

import (
	"fmt"
	"github.com/area3001/goira/comm"
	"github.com/area3001/goira/core"
)

type Service struct {
	n      *comm.NatsClient
	reader *Reader
}

func (c *Service) ListDevices() ([]Device, error) {
	return c.reader.List()
}

func (c *Service) Ping() error {
	return c.n.Broadcast("ping")
}

func (c *Service) SetMode(mac string, mode *core.Mode) error {
	if mode == nil {
		return fmt.Errorf("no mode provided")
	}

	return c.n.Call("mode", mac, []byte(fmt.Sprintf("%d", mode.Code)))
}

func (c *Service) Blink(mac string, times int) error {
	return c.n.Call("blink", mac, []byte(fmt.Sprintf("%d", times)))
}

func (c *Service) Reset(mac string, delayMs int) error {
	return c.n.Call("reset", mac, []byte(fmt.Sprintf("%d", delayMs)))
}

func (c *Service) SendDmx(mac string, data []byte) error {
	return c.n.Call("dmx", mac, data)
}

func (c *Service) SendDeltaDmx(mac string, data []byte) error {
	return c.n.Call("deltadmx", mac, data)
}

func (c *Service) SendRgb(mac string, data []byte) error {
	return c.n.Call("rgb", mac, data)
}

func (c *Service) SendFx(mac string, effect *core.Effect) error {
	return c.n.Call("fx", mac, effect.ToBytes())
}
