package devices

import (
	"encoding/json"
	"fmt"
	"github.com/nats-io/nats.go"
)

type Reader struct {
	devices nats.KeyValue
}

func (r *Reader) Keys() ([]string, error) {
	return r.devices.Keys()
}

func (r *Reader) List() ([]Device, error) {
	result := make([]Device, 0)

	keys, err := r.devices.Keys()
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve the list of devices: %w", err)
	}

	for _, key := range keys {
		entry, err := r.devices.Get(key)
		if err != nil {
			return nil, fmt.Errorf("unable to retrieve info for device %s: %w", key, err)
		}

		var dev Device
		if err := json.Unmarshal(entry.Value(), &dev); err != nil {
			return nil, fmt.Errorf("unable to decode the info for device %s: %w", key, err)
		}

		result = append(result, dev)
	}

	return result, nil
}

func (r *Reader) Get(key string) (*Device, error) {
	entry, err := r.devices.Get(key)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve device with key %s: %w", key, err)
	}

	var dev Device
	if err := json.Unmarshal(entry.Value(), &dev); err != nil {
		return nil, fmt.Errorf("unable to decode the info for device %s: %w", key, err)
	}

	return &dev, nil
}
