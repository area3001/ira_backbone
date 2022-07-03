package sdk

import (
	"bytes"
	"image/color"
	"reflect"
	"testing"
)

func Test_toRgbBuffer(t *testing.T) {
	type args struct {
		offset uint16
		colors []color.RGBA
	}
	tests := []struct {
		name string
		args args
		want *bytes.Buffer
	}{
		{name: "single pixel", args: struct {
			offset uint16
			colors []color.RGBA
		}{offset: 2, colors: []color.RGBA{{255, 0, 0, 0}}},
			want: bytes.NewBuffer([]byte{0, 2, 0, 1, 255, 0, 0})},

		{name: "two pixels", args: struct {
			offset uint16
			colors []color.RGBA
		}{offset: 2, colors: []color.RGBA{{255, 0, 0, 0}, {0, 255, 0, 0}}},
			want: bytes.NewBuffer([]byte{0, 2, 0, 2, 255, 0, 0, 0, 255, 0})},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := toRgbBuffer(tt.args.offset, tt.args.colors); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("toRgbBuffer() = %x, want %x", got, tt.want)
			}
		})
	}
}
