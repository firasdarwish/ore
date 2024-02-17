package ore

import (
	"fmt"
	"strings"
)

type KeyStringer any

func oreKey(key []KeyStringer) string {
	l := len(key)

	if key == nil || l == 0 {
		return ""
	}

	if l == 1 {
		return fmt.Sprintf("%v", key[0])
	}

	var keys []string

	for _, stringer := range key {
		keys = append(keys, fmt.Sprintf("%v", stringer))
	}

	return strings.Join(keys, ":/#+")
}
