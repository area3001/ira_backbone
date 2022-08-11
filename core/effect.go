package core

import (
	"encoding/json"
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
	Rainbow       = &EffectKind{Code: 5, Name: "Rainbow", Description: "Rainbow", AllowedParams: []*EffectParam{}}
	RainbowSpread = &EffectKind{Code: 6, Name: "Rainbow Spread", Description: "Rainbow Spread", AllowedParams: []*EffectParam{}}
)

var Effects = []*EffectKind{
	PixelLoopFx, RandomPixelLoopFx, ForegroundBackgroundLoopFx, ForegroundBackgroundSwitchFx, Fire2021Fx, Rainbow, RainbowSpread,
}

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
			result.Background = bg

		case SpeedFxParam:
			val, err := ParseStringToUint8(param)
			if err != nil {
				return nil, fmt.Errorf("invalid speed: %w", err)
			}
			result.Speed = 255 - val

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
	Foreground color.RGBA  `json:"foreground" swaggertype:"string"`
	Background color.RGBA  `json:"background" swaggertype:"string"`
}

func (e *Effect) ToBytes() []byte {
	return []byte{
		e.Kind.Code, e.Speed, e.Crossfade,
		e.Foreground.R, e.Foreground.G, e.Foreground.B,
		e.Background.R, e.Background.G, e.Background.B,
	}
}

func (e Effect) MarshalJSON() ([]byte, error) {
	we := wireEffect{
		Kind:       e.Kind,
		Speed:      e.Speed,
		Crossfade:  e.Crossfade,
		Foreground: fmt.Sprintf("#%02x%02x%02x", e.Foreground.R, e.Foreground.G, e.Foreground.B),
		Background: fmt.Sprintf("#%02x%02x%02x", e.Background.R, e.Background.G, e.Background.B),
	}
	return json.Marshal(we)
}

func (e *Effect) UnmarshalJSON(data []byte) error {
	var we wireEffect
	if err := json.Unmarshal(data, &we); err != nil {
		return err
	}

	fg, err := ParseHexColor(we.Foreground)
	if err != nil {
		return fmt.Errorf("unable to parse foreground color: %w", err)
	}

	bg, err := ParseHexColor(we.Background)
	if err != nil {
		return fmt.Errorf("unable to parse background color: %w", err)
	}

	*e = Effect{Kind: we.Kind, Speed: we.Speed, Crossfade: we.Crossfade, Foreground: fg, Background: bg}

	return nil
}

type wireEffect struct {
	Kind       *EffectKind `json:"kind"`
	Speed      uint8       `json:"speed"`
	Crossfade  uint8       `json:"crossfade"`
	Foreground string      `json:"foreground"`
	Background string      `json:"background"`
}
