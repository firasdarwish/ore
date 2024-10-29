package ore

import (
	"context"
	"fmt"
	"sync"
)

var (
	lock      = &sync.RWMutex{}
	isBuilt   = false
	container = map[typeID][]serviceResolver{}
)

type Creator[T any] interface {
	New(ctx context.Context) (T, context.Context)
}

// Generates a unique identifier for a service resolver based on type and key(s)
func getTypeId(pointerTypeName pointerTypeName, key []KeyStringer) typeID {
	for _, stringer := range key {
		if stringer == nil {
			panic(nilKey)
		}
	}
	customKey := oreKey(key)
	tt := fmt.Sprintf("%s:%v", pointerTypeName, customKey)
	return typeID(tt)
}

// Generates a unique identifier for a service resolver based on type and key(s)
func typeIdentifier[T any](key []KeyStringer) typeID {
	return getTypeId(getPointerTypeName[T](), key)
}

// Appends a service resolver to the container with type and key
func appendToContainer[T any](resolver serviceResolver, key []KeyStringer) {
	if isBuilt {
		panic(alreadyBuiltCannotAdd)
	}

	typeId := typeIdentifier[T](key)

	lock.Lock()
	container[typeId] = append(container[typeId], resolver)
	lock.Unlock()
}

func replaceServiceResolver(typeId typeID, index int, resolver serviceResolver) {
	lock.Lock()
	container[typeId][index] = resolver
	lock.Unlock()
}

func Build() {
	if isBuilt {
		panic(alreadyBuilt)
	}

	isBuilt = true
}
