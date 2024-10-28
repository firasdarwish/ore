package ore

import (
	"context"
	"fmt"
	"sync"
)

var (
	lock      = &sync.RWMutex{}
	isBuilt   = false
	container = map[string][]any{}
)

type Creator[T any] interface {
	New(ctx context.Context) (T, context.Context)
}

// Generates a unique identifier for an entry based on type and key(s)
func typeIdentifier[T any](key []KeyStringer) string {
	for _, stringer := range key {
		if stringer == nil {
			panic(nilKey)
		}
	}

	var mockType *T
	customKey := oreKey(key)
	tt := fmt.Sprintf("%c:%v", mockType, customKey)
	return tt
}

// Appends an entry to the container with type and key
func appendToContainer[T any](entry entry[T], key []KeyStringer) {
	if isBuilt {
		panic(alreadyBuiltCannotAdd)
	}

	typeId := typeIdentifier[T](key)

	lock.Lock()
	container[typeId] = append(container[typeId], entry)
	lock.Unlock()
}

func replaceEntry[T any](typeId string, index int, entry entry[T]) {
	lock.Lock()
	container[typeId][index] = entry
	lock.Unlock()
}

func Build() {
	if isBuilt {
		panic(alreadyBuilt)
	}

	isBuilt = true
}
