package ore

import (
	"reflect"
)

func validateOreKeyType(key KeyStringer) {
	switch key.(type) {
	case nil:
	case string:
	case int, int8, int16, int32, int64:
	case uint, uint8, uint16, uint32, uint64:
	case float32, float64:
		break
	default:
		panic(invalidKeyType(reflect.TypeOf(key)))
	}
}
