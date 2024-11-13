package ore

import (
	"reflect"
	"strconv"
	"strings"
)

type KeyStringer any

func oreKey(key ...KeyStringer) string {
	if key == nil {
		return ""
	}

	l := len(key)

	if l == 1 {
		keyT, kV := stringifyOreKey(key[0])
		return keyT + kV
	}

	var sb strings.Builder

	for _, s := range key {
		keyT, keyV := stringifyOreKey(s)
		sb.WriteString(keyT)
		sb.WriteString(keyV)
	}

	return sb.String()
}

func stringifyOreKey(key KeyStringer) (string, string) {
	switch key.(type) {
	case nil:
		return "n", ""
	case string:
		return "s", key.(string)
	case int:
		return "i", strconv.Itoa(key.(int))
	case int8:
		return "i8", strconv.FormatInt(int64(key.(int8)), 36)
	case int16:
		return "i16", strconv.FormatInt(int64(key.(int16)), 36)
	case int32:
		return "i32", strconv.FormatInt(int64(key.(int32)), 36)
	case int64:
		return "i64", strconv.FormatInt(key.(int64), 36)

	case uint:
		return "ui", strconv.FormatUint(uint64(key.(uint)), 36)
	case uint8:
		return "ui8", strconv.FormatUint(uint64(key.(uint8)), 36)
	case uint16:
		return "ui16", strconv.FormatUint(uint64(key.(uint16)), 36)
	case uint32:
		return "ui32", strconv.FormatUint(uint64(key.(uint32)), 36)
	case uint64:
		return "ui64", strconv.FormatUint(key.(uint64), 36)
	case float32:
		return "f32", strconv.FormatFloat(float64(key.(float32)), 'x', -1, 32)
	case float64:
		return "f64", strconv.FormatFloat(key.(float64), 'x', -1, 64)
	default:
		panic(invalidKeyType(reflect.TypeOf(key)))
	}
}
