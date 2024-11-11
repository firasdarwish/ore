package ore

import (
	"fmt"
	"strconv"
	"strings"
)

type KeyStringer any

type stringer interface {
	String() string
}

func oreKey(key []KeyStringer) string {
	if key == nil {
		return ""
	}

	l := len(key)

	if l == 1 {
		return stringifyOreKey(key[0])
	}

	var sb strings.Builder

	for _, s := range key {
		sb.WriteString(stringifyOreKey(s))
	}

	return sb.String()
}

func stringifyOreKey(key KeyStringer) string {
	switch key.(type) {
	case nil:
		return ""
	case string:
		return key.(string)
	case bool:
		return strconv.FormatBool(key.(bool))
	case int:
		return strconv.Itoa(key.(int))
	case int8:
		return strconv.FormatInt(int64(key.(int8)), 36)
	case int16:
		return strconv.FormatInt(int64(key.(int16)), 36)
	case int32:
		return strconv.FormatInt(int64(key.(int32)), 36)
	case int64:
		return strconv.FormatInt(key.(int64), 36)

	case uint:
		return strconv.FormatUint(uint64(key.(uint)), 36)
	case uint8:
		return strconv.FormatUint(uint64(key.(uint8)), 36)
	case uint16:
		return strconv.FormatUint(uint64(key.(uint16)), 36)
	case uint32:
		return strconv.FormatUint(uint64(key.(uint32)), 36)
	case uint64:
		return strconv.FormatUint(key.(uint64), 36)
	case float32:
		return strconv.FormatFloat(float64(key.(float32)), 'x', -1, 32)
	case float64:
		return strconv.FormatFloat(key.(float64), 'x', -1, 64)
	case stringer:
		return key.(stringer).String()

	default:
		return stringifyOreKeyUnknown(key)
	}
}

func stringifyOreKeyUnknown(key KeyStringer) string {
	return fmt.Sprintf("%v", key)
}
