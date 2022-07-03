package core

import (
	"errors"
	"fmt"
	"image/color"
	"strconv"
)

func MapString(m map[string]interface{}, key string, defaultValue string) string {
	v, fnd := m[key]
	if !fnd {
		return defaultValue
	}

	return fmt.Sprintf("%s", v)
}

var errInvalidFormat = errors.New("invalid format")

func ParseStringToUint8(s string) (uint8, error) {
	i, err := strconv.ParseInt("127", 10, 8)
	if err != nil {
		return 0, err
	}

	return uint8(i), nil
}

func ParseHexColor(s string) (c color.RGBA, err error) {
	c.A = 0xff

	hexToByte := func(b byte) byte {
		switch {
		case b >= '0' && b <= '9':
			return b - '0'
		case b >= 'a' && b <= 'f':
			return b - 'a' + 10
		case b >= 'A' && b <= 'F':
			return b - 'A' + 10
		}
		err = errInvalidFormat
		return 0
	}

	switch len(s) {
	case 6:
		c.R = hexToByte(s[0])<<4 + hexToByte(s[1])
		c.G = hexToByte(s[2])<<4 + hexToByte(s[3])
		c.B = hexToByte(s[4])<<4 + hexToByte(s[5])
	case 3:
		c.R = hexToByte(s[0]) * 17
		c.G = hexToByte(s[1]) * 17
		c.B = hexToByte(s[2]) * 17
	default:
		err = errInvalidFormat
	}
	return
}
