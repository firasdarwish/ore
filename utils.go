package ore

import (
	"fmt"
	"strings"
)

type specialContextKey string
type specialOreKey int

type contextKey struct {
	typeID
	containerID int32
	resolverID  int
}
type typeID struct {
	pointerTypeName pointerTypeName
	oreKey          any //comparable
}
type pointerTypeName string

func (this *Container) clearAll() {
	this.resolvers = make(map[typeID][]serviceResolver)
	this.aliases = make(map[pointerTypeName][]pointerTypeName)
	this.isSealed = false
	this.DisableValidation = false
	this.name = "DEFAULT"
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
	return strings.TrimLeft(string(ptn), "*")
}

func (this typeID) String() string {
	return fmt.Sprintf("(name={%s}, key='%s')", getUnderlyingTypeName(this.pointerTypeName), this.oreKey)
}
