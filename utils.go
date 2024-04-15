package ore

import "fmt"

func isNil[T comparable](impl T) bool {
	var mock T
	return impl == mock
}

func clearAll() {
	container = make(map[string][]any)
	isBuilt = false
}

func contextValueId(typeId string, index int) string {
	return fmt.Sprintln(typeId, index)
}
