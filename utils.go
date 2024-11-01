package ore

import "fmt"

type contextValueID string
type typeID string
type pointerTypeName string

func isNil[T comparable](impl T) bool {
	var mock T
	return impl == mock
}

func clearAll() {
	container = make(map[typeID][]serviceResolver)
	aliases = make(map[pointerTypeName][]pointerTypeName)
	isBuilt = false
}

func getContextValueID(typeId typeID, index int) contextValueID {
	return contextValueID(fmt.Sprintln(typeId, index))
}

// Get type name of *T.
// it allocates less memory and is faster than `reflect.TypeFor[*T]().String()`
func getPointerTypeName[T any]() pointerTypeName {
	var mockValue *T
	return pointerTypeName(fmt.Sprintf("%T", mockValue))
}
