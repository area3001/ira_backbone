package core

const (
	PixelLoopFx                  = byte(0)
	RandomPixelLoopFx            = byte(1)
	ForegroundBackgroundLoopFx   = byte(2)
	ForegroundBackgroundSwitchFx = byte(3)
	Fire2021Fx                   = byte(4)
)

type Effect struct {
	Kind            byte
	Speed           byte
	Crossfade       byte
	ForegroundRed   byte
	ForegroundGreen byte
	ForegroundBlue  byte
	BackgroundRed   byte
	BackgroundGreen byte
	BackgroundBlue  byte
}

func (e *Effect) toBytes() []byte {
	return []byte{
		e.Kind, e.Speed, e.Crossfade,
		e.ForegroundRed, e.ForegroundGreen, e.ForegroundBlue,
		e.BackgroundRed, e.BackgroundGreen, e.BackgroundBlue,
	}
}
