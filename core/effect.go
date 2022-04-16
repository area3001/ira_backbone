package core

const (
	PixelLoopFx                  = byte(0)
	RandomPixelLoopFx            = byte(1)
	ForegroundBackgroundLoopFx   = byte(2)
	ForegroundBackgroundSwitchFx = byte(3)
	Fire2021Fx                   = byte(4)
)

type Effect struct {
	Kind            byte `json:"kind"`
	Speed           byte `json:"speed"`
	Crossfade       byte `json:"crossfade"`
	ForegroundRed   byte `json:"foregroundRed"`
	ForegroundGreen byte `json:"foregroundGreen"`
	ForegroundBlue  byte `json:"foregroundBlue"`
	BackgroundRed   byte `json:"backgroundRed"`
	BackgroundGreen byte `json:"backgroundGreen"`
	BackgroundBlue  byte `json:"backgroundBlue"`
}

func (e *Effect) ToBytes() []byte {
	return []byte{
		e.Kind, e.Speed, e.Crossfade,
		e.ForegroundRed, e.ForegroundGreen, e.ForegroundBlue,
		e.BackgroundRed, e.BackgroundGreen, e.BackgroundBlue,
	}
}
