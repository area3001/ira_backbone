package core

import (
	"fmt"
	"image/color"
)

var (
	PixelLoopFx = &EffectKind{Code: 0, Name: "Pixel Loop", Description: "Loop a single pixel", AllowedParams: []*EffectParam{
		ForegroundFxParam, BackgroundFxParam, SpeedFxParam,
	}}
	RandomPixelLoopFx = &EffectKind{Code: 1, Name: "Random Pixel Loop", Description: "Loop a single pixel", AllowedParams: []*EffectParam{
		ForegroundFxParam, BackgroundFxParam, SpeedFxParam,
	}}
	ForegroundBackgroundLoopFx = &EffectKind{Code: 2, Name: "Foreground/Background Loop", Description: "Loop the foreground and the background", AllowedParams: []*EffectParam{
		ForegroundFxParam, BackgroundFxParam, SpeedFxParam,
	}}
	ForegroundBackgroundSwitchFx = &EffectKind{Code: 3, Name: "Foreground/Background Switch", Description: "Switch the foreground and the background", AllowedParams: []*EffectParam{
		ForegroundFxParam, BackgroundFxParam, SpeedFxParam,
	}}
	Fire2021Fx = &EffectKind{Code: 4, Name: "Fire 2021", Description: "Fire !!!", AllowedParams: []*EffectParam{
		SpeedFxParam, CrossfadeFxParam,
	}}
)

var (
	ForegroundFxParam = &EffectParam{"fg", "The foreground color"}
	BackgroundFxParam = &EffectParam{"bg", "The background color"}
	SpeedFxParam      = &EffectParam{"speed", "The speed of the effect"}
	CrossfadeFxParam  = &EffectParam{"crossover", "The amount of crossover between effect cycles"}
)

type EffectParam struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type EffectKind struct {
	Name          string         `json:"name"`
	Description   string         `json:"description"`
	Code          uint8          `json:"code"`
	AllowedParams []*EffectParam `json:"allowedParams"`
}

func NewEffect(kind *EffectKind, params map[string]string) (*Effect, error) {
	result := &Effect{
		Kind:       kind,
		Speed:      0,
		Crossfade:  0,
		Foreground: color.RGBA{},
		Background: color.RGBA{},
	}

	for _, p := range kind.AllowedParams {
		param, fnd := params[p.Name]
		if !fnd {
			continue
		}

		switch p {
		case ForegroundFxParam:
			fg, err := ParseHexColor(param)
			if err != nil {
				return nil, fmt.Errorf("invalid foreground color: %w", err)
			}
			result.Foreground = fg

		case BackgroundFxParam:
			bg, err := ParseHexColor(param)
			if err != nil {
				return nil, fmt.Errorf("invalid background color: %w", err)
			}
			result.Foreground = bg

		case SpeedFxParam:
			val, err := ParseStringToUint8(param)
			if err != nil {
				return nil, fmt.Errorf("invalid speed: %w", err)
			}
			result.Speed = val

		case CrossfadeFxParam:
			val, err := ParseStringToUint8(param)
			if err != nil {
				return nil, fmt.Errorf("invalid crossfade: %w", err)
			}
			result.Crossfade = val

		}
	}

	return result, nil
}

type Effect struct {
	Kind       *EffectKind `json:"kind"`
	Speed      uint8       `json:"speed"`
	Crossfade  uint8       `json:"crossfade"`
	Foreground color.RGBA  `json:"foreground"`
	Background color.RGBA  `json:"background"`
}

func (e *Effect) ToBytes() []byte {
	return []byte{
		e.Kind.Code, e.Speed, e.Crossfade,
		e.Foreground.R, e.Foreground.G, e.Foreground.B,
		e.Background.R, e.Background.G, e.Background.B,
	}
}
