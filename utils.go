package ore

import "fmt"

type contextKey struct {
	typeID
	index int
}
type typeID struct {
	pointerTypeName pointerTypeName
	oreKey          string
}
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

// Get type name of *T.
// it allocates less memory and is faster than `reflect.TypeFor[*T]().String()`
func getPointerTypeName[T any]() pointerTypeName {
	var mockValue *T
	return pointerTypeName(fmt.Sprintf("%T", mockValue))
}
