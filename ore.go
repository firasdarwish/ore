package ore

import (
	"context"
	"sync"
)

var (
	lock      = &sync.RWMutex{}
	isBuilt   = false
	container = map[typeID][]serviceResolver{}

	//map the alias type (usually an interface) to the original types (usually implementations of the interface)
	aliases = map[pointerTypeName][]pointerTypeName{}
)

type Creator[T any] interface {
	New(ctx context.Context) (T, context.Context)
}

// Generates a unique identifier for a service resolver based on type and key(s)
func getTypeID(pointerTypeName pointerTypeName, key []KeyStringer) typeID {
	for _, stringer := range key {
		if stringer == nil {
			panic(nilKey)
		}
	}
	return typeID{pointerTypeName, oreKey(key)}
}

// Generates a unique identifier for a service resolver based on type and key(s)
func typeIdentifier[T any](key []KeyStringer) typeID {
	return getTypeID(getPointerTypeName[T](), key)
}

// Appends a service resolver to the container with type and key
func appendToContainer[T any](resolver serviceResolver, key []KeyStringer) {
	if isBuilt {
		panic(alreadyBuiltCannotAdd)
	}

	typeID := typeIdentifier[T](key)

	lock.Lock()
	container[typeID] = append(container[typeID], resolver)
	lock.Unlock()
}

func replaceServiceResolver(typeId typeID, index int, resolver serviceResolver) {
	lock.Lock()
	container[typeId][index] = resolver
	lock.Unlock()
}

func appendToAliases[TInterface, TImpl any]() {
	originalType := getPointerTypeName[TImpl]()
	aliasType := getPointerTypeName[TInterface]()
	if originalType == aliasType {
		return
	}
	lock.Lock()
	for _, ot := range aliases[aliasType] {
		if ot == originalType {
			return //already registered
		}
	}
	aliases[aliasType] = append(aliases[aliasType], originalType)
	lock.Unlock()
}

func Build() {
	if isBuilt {
		panic(alreadyBuilt)
	}

	isBuilt = true
}
