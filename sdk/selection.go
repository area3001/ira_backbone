package sdk

import "sync"

type DeviceSelection []*Device

type Performer func(dev *Device)

func (s DeviceSelection) Perform(performer Performer) {
	wg := sync.WaitGroup{}
	wg.Add(len(s))

	for _, dev := range s {
		go func(dev *Device) {
			defer wg.Add(-1)
			performer(dev)

		}(dev)
	}

	wg.Wait()
}
