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

	if s[0] != '#' {
		return c, errInvalidFormat
	}

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
	case 7:
		c.R = hexToByte(s[1])<<4 + hexToByte(s[2])
		c.G = hexToByte(s[3])<<4 + hexToByte(s[4])
		c.B = hexToByte(s[5])<<4 + hexToByte(s[6])
	case 4:
		c.R = hexToByte(s[1]) * 17
		c.G = hexToByte(s[2]) * 17
		c.B = hexToByte(s[3]) * 17
	default:
		err = errInvalidFormat
	}
	return
}
