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

// Generates a unique identifier for an entry based on type and key(s)
func typeId[T any](key []KeyStringer) string {
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

// RegisterLazyFunc Registers a lazily initialized value using an `Initializer[T]` function signature
func RegisterLazyFunc[T any](entryType RegistrationType, initializer Initializer[T], key ...KeyStringer) {
	e := entry[T]{
		registrationType:     entryType,
		anonymousInitializer: &initializer,
	}
	appendToContainer[T](e, key)
}

// RegisterLazyCreator Registers a lazily initialized value using a `Creator[T]` interface
func RegisterLazyCreator[T any](entryType RegistrationType, creator Creator[T], key ...KeyStringer) {
	e := entry[T]{
		registrationType: entryType,
		creatorInstance:  creator,
	}
	appendToContainer[T](e, key)
}

// RegisterEagerSingleton Registers an eagerly instantiated singleton value
func RegisterEagerSingleton[T any](impl T, key ...KeyStringer) {
	e := entry[T]{
		registrationType: Singleton,
		concrete:         &impl,
	}
	appendToContainer[T](e, key)
}

// Get Retrieves an instance based on type and key (throws panic for missing or invalid implementations)
func Get[T any](ctx context.Context, key ...KeyStringer) (T, context.Context) {
	tId := typeId[T](key)

	lock.RLock()
	o, ok := container[tId]
	lock.RUnlock()

	if !ok {
		panic(noValidImplementation[T]())
	}

	oLen := len(o)

	if oLen == 0 {
		panic(noValidImplementation[T]())
	}

	index := oLen - 1

	i, okInst := o[index].(entry[T])
	if !okInst {
		panic(noValidImplementation[T]())
	}

	ctxValId := fmt.Sprintln(tId, index)

	con, ctx := i.load(ctx, ctxValId)
	if i.registrationType == Singleton {
		replace[T](tId, index, i)
	}

	return con, ctx
}

// GetList Retrieves a list of instances based on type and key (throws panic for missing or invalid implementations)
func GetList[T any](ctx context.Context, key ...KeyStringer) ([]T, context.Context) {
	tId := typeId[T](key)

	lock.RLock()
	o, ok := container[tId]
	lock.RUnlock()

	if !ok {
		panic(noValidImplementation[T]())
	}

	oLen := len(o)

	if oLen == 0 {
		panic(noValidImplementation[T]())
	}

	arrr := make([]T, oLen)

	for index := 0; index < oLen; index++ {
		i, okInst := o[index].(entry[T])
		if !okInst {
			panic(noValidImplementation[T]())
		}

		ctxValId := fmt.Sprintln(tId, index)
		connn, ctxxx := i.load(ctx, ctxValId)

		if i.registrationType == Singleton {
			replace[T](tId, index, i)
		}

		arrr[index] = connn
		ctx = ctxxx
	}

	return arrr, ctx
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
