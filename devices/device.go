package devices

import (
	"encoding/json"
	"github.com/area3001/goira/core"
	"github.com/nats-io/nats.go"
	"time"
)

/**
{
	"mac_string": "f0_4_8a_96_2b_c8",
	"IP":"352387594",
"NAME":"",
	"HWTYPE":"IRA2020",
	"HWREV":"Rev.01",
	"EXTMODE":"15",
	"MODE":"11",
	"VERSION":"7",

	"pixel_length": 100,
	"fx": 255,
	"fx_speed": 255,
	"fx_xfade": 255,
	"fx_fgnd_r": 255,
	"fx_fgnd_g": 255,
	"fx_fgnd_b": 255,
	"fx_bgnd_r": 255,
	"fx_bgnd_g": 255,
	"fx_bgnd_b": 255
}
*/

const (
	MacVarName     = "mac_string"
	NameVarName    = "NAME"
	IpVarName      = "IP"
	HwTypeVarName  = "HWTYPE"
	HwRevVarName   = "HWREV"
	ExtModeVarName = "EXTMODE"
	ModeVarName    = "MODE"
	VersionVarName = "VERSION"

	PixelLengthVar       = "pixel_length"
	FxVar                = "fx"
	FxSpeedVar           = "fx_speed"
	FxXfadeVar           = "fx_xfade"
	FxForegroundRedVar   = "fx_fgnd_r"
	FxForegroundGreenVar = "fx_fgnd_g"
	FxForegroundBlueVar  = "fx_fgnd_b"
	FxBackgroundRedVar   = "fx_bgnd_r"
	FxBackgroundGreenVar = "fx_bgnd_g"
	FxBackgroundBlueVar  = "fx_bgnd_b"
)

type Device struct {
	MAC          string         `json:"mac"`
	Name         string         `json:"name"`
	IP           string         `json:"ip"`
	Hardware     Hardware       `json:"hardware"`
	Mode         int            `json:"mode"`
	ExternalMode int            `json:"external_mode"`
	LastBeat     time.Time      `json:"last_beat"`
	Version      int            `json:"version"`
	Config       map[string]int `json:"config"`
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
		MAC:          core.MapString(data, MacVarName, ""),
		Name:         core.MapString(data, NameVarName, ""),
		IP:           core.MapString(data, IpVarName, ""),
		ExternalMode: core.MapInt(data, ExtModeVarName, -1),
		Mode:         core.MapInt(data, ModeVarName, -1),
		Hardware: Hardware{
			Kind:    core.MapString(data, HwTypeVarName, ""),
			Version: core.MapString(data, HwRevVarName, ""),
		},
		LastBeat: time.Now(),
		Version:  core.MapInt(data, VersionVarName, -1),
		Config: map[string]int{
			PixelLengthVar:       core.MapInt(data, PixelLengthVar, -1),
			FxVar:                core.MapInt(data, FxVar, -1),
			FxXfadeVar:           core.MapInt(data, FxXfadeVar, -1),
			FxSpeedVar:           core.MapInt(data, FxSpeedVar, -1),
			FxForegroundRedVar:   core.MapInt(data, FxForegroundRedVar, -1),
			FxForegroundGreenVar: core.MapInt(data, FxForegroundGreenVar, -1),
			FxForegroundBlueVar:  core.MapInt(data, FxForegroundBlueVar, -1),
			FxBackgroundRedVar:   core.MapInt(data, FxBackgroundRedVar, -1),
			FxBackgroundGreenVar: core.MapInt(data, FxBackgroundGreenVar, -1),
			FxBackgroundBlueVar:  core.MapInt(data, FxBackgroundBlueVar, -1),
		},
	}, nil
}
