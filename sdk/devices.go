package sdk

import (
	"encoding/json"
	"github.com/area3001/goira/comm"
	"github.com/nats-io/nats.go"
	"strings"
)

type Devices struct {
	store nats.KeyValue
	nc    *comm.NatsClient
}

func (d *Devices) Keys() ([]string, error) {
	return d.store.Keys()
}

func (d *Devices) List() ([]*Device, error) {
	keys, err := d.store.Keys()
	if err != nil {
		return nil, err
	}

	result := make([]*Device, len(keys))
	for idx, k := range keys {
		entry, err := d.store.Get(k)
		if err != nil {
			return nil, err
		}

		var dev DeviceMeta
		if err := json.Unmarshal(entry.Value(), &dev); err != nil {
			return nil, err
		}

		result[idx] = &Device{
			Meta:  &dev,
			store: d.store,
			nc:    d.nc,
		}
	}

	return result, nil
}

func (d *Devices) Select(selector string) (DeviceSelection, error) {
	if strings.ToUpper(selector) == "ALL" {
		result, err := d.List()
		if err != nil {
			return nil, err
		}

		return result, nil
	}

	selectors := strings.Split(selector, ",")
	result := make([]*Device, len(selectors))
	for idx, s := range selectors {
		dev, err := d.Device(strings.TrimSpace(s))
		if err != nil {
			return nil, err
		}

		result[idx] = dev
	}

	return result, nil
}

func (d *Devices) Device(key string) (*Device, error) {
	kv, err := d.store.Get(key)
	if err != nil {
		return nil, err
	}

	var meta DeviceMeta
	if err := json.Unmarshal(kv.Value(), &meta); err != nil {
		return nil, err
	}

	return &Device{
		Meta:  &meta,
		store: d.store,
		nc:    d.nc,
	}, nil
}

func (d *Devices) Forget(key string) error {
	return d.store.Purge(key)
}

func (d *Devices) Sync() error {
	return d.nc.Broadcast("ping")
}
