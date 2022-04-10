package devices

import (
	"encoding/json"
	"fmt"
	"github.com/nats-io/nats.go"
)

type Writer struct {
	devices nats.KeyValue
}

func (w *Writer) Register(dev *Device) error {
	b, err := json.Marshal(dev)
	if err != nil {
		return fmt.Errorf("unable to encode device data: %w", err)
	}

	_, err = w.devices.Put(dev.MAC, b)
	if err != nil {
		return fmt.Errorf("unable to store device data: %w", err)
	}

	return nil
}
