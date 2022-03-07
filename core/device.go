package core

import (
	"encoding/json"
	"github.com/nats-io/nats.go"
	"time"
)

const (
	MacVarName     = "mac_address"
	IpVarName      = "IP"
	HwTypeVarName  = "HWTYPE"
	HwRevVarName   = "HWREV"
	ModeVarName    = "MODE"
	ExtModeVarName = "EXTMODE"
)

type Device struct {
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

func ParseDevice(m *nats.Msg) (*Device, error) {
	var data map[string]interface{}
	if err := json.Unmarshal(m.Data, &data); err != nil {
		return nil, err
	}

	return &Device{
		MAC:          MapString(data, MacVarName, ""),
		IP:           MapString(data, IpVarName, ""),
		ExternalMode: MapString(data, ExtModeVarName, ""),
		Mode:         MapString(data, ModeVarName, ""),
		Hardware: Hardware{
			Kind:    MapString(data, HwTypeVarName, ""),
			Version: MapString(data, HwRevVarName, ""),
		},
		LastBeat: time.Now(),
	}, nil
}
