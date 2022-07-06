package core

import (
	"image/color"
	"reflect"
	"testing"
)

func TestEffect_ToBytes(t *testing.T) {
	type fields struct {
		Kind       *EffectKind
		Speed      uint8
		Crossfade  uint8
		Foreground color.RGBA
		Background color.RGBA
	}
	tests := []struct {
		name   string
		fields fields
		want   []byte
	}{
		{
			"simple effect",
			fields{ForegroundBackgroundSwitchFx, 150, 0, color.RGBA{255, 0, 0, 0}, color.RGBA{0, 0, 150, 0}},
			[]byte{3, 150, 0, 255, 0, 0, 0, 0, 150}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Effect{
				Kind:       tt.fields.Kind,
				Speed:      tt.fields.Speed,
				Crossfade:  tt.fields.Crossfade,
				Foreground: tt.fields.Foreground,
				Background: tt.fields.Background,
			}
			if got := e.ToBytes(); !reflect.DeepEqual(got, tt.want) {
				/*
					First Byte = FX selection
					Second Byte = FX Speed
					Third Byte = FX Crossfade
					4Th Byte = FGND R
					5Th = FGND G
					6TH = FGND B
					7TH = BGND R
					8TH = BGND G
					9TH = BGND B
				*/

				t.Errorf("ToBytes() = %x, want %x", got, tt.want)
			}
		})
	}
}
