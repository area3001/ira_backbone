package core

import "fmt"

func MapString(m map[string]interface{}, key string, defaultValue string) string {
	v, fnd := m[key]
	if !fnd {
		return defaultValue
	}

	return fmt.Sprintf("%s", v)
}
