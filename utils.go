package ore

import (
	"fmt"
	"strings"
)

type specialContextKey string

type contextKey struct {
	typeID
	containerID int32
	resolverID  int
}
type typeID struct {
	pointerTypeName pointerTypeName
	oreKey          KeyStringer
}
type pointerTypeName string

func (this *Container) clearAll() {
	this.resolvers = make(map[typeID][]serviceResolver)
	this.aliases = make(map[pointerTypeName][]pointerTypeName)
	this.isSealed = false
	this.DisableValidation = false
}

func clearAll() {
	DefaultContainer.clearAll()
}

// Get type name of *T.
// it allocates less memory and is faster than `reflect.TypeFor[*T]().String()`
func getPointerTypeName[T any]() pointerTypeName {
	var mockValue *T
	return pointerTypeName(fmt.Sprintf("%T", mockValue))
}

func getUnderlyingTypeName(ptn pointerTypeName) string {
	s := string(ptn)
	index := strings.Index(s, "*")
	if index == -1 {
		return s // no '*' found, return the original string
	}
	return s[:index] + s[index+1:]
}

func (this typeID) String() string {
	return fmt.Sprintf("(name={%s}, key='%s')", getUnderlyingTypeName(this.pointerTypeName), this.oreKey)
}
