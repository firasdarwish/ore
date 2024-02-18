package ore

import (
	"fmt"
	"sync"
)

var (
	lock      = &sync.RWMutex{}
	isBuilt   = false
	container = map[string][]any{}
)

// Generates a unique identifier for an entry based on type and key(s)
func typeId[T any](key []KeyStringer) string {
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
func appendToContainer[T any](e entry[T], key []KeyStringer) {
	if isBuilt {
		panic(alreadyBuiltCannotAdd)
	}

	tId := typeId[T](key)

	lock.Lock()
	container[tId] = append(container[tId], e)
	lock.Unlock()
}

func replace[T any](typeId string, index int, entry entry[T]) {
	container[typeId][index] = entry
}

func clearAll() {
	container = make(map[string][]any)
}

func Build() {
	if isBuilt {
		panic(alreadyBuilt)
	}

	isBuilt = true
}
