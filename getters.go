package ore

import (
	"context"
)

// Get Retrieves an instance based on type and key (panics if no valid implementations)
func Get[T any](ctx context.Context, key ...KeyStringer) (T, context.Context) {
	// generate type identifier
	typeId := typeIdentifier[T](key)

	// try to get entry from container
	lock.RLock()
	entries, entryExists := container[typeId]
	lock.RUnlock()

	if !entryExists {
		panic(noValidImplementation[T]())
	}

	entriesCount := len(entries)

	if entriesCount == 0 {
		panic(noValidImplementation[T]())
	}

	// index of the last implementation
	index := entriesCount - 1

	implementation := entries[index].(entry[T])

	service, ctx, updateEntry := implementation.load(ctx, contextValueId(typeId, index))
	if updateEntry {
		replaceEntry[T](typeId, index, implementation)
	}

	return service, ctx
}

// GetList Retrieves a list of instances based on type and key
func GetList[T any](ctx context.Context, key ...KeyStringer) ([]T, context.Context) {
	// generate type identifier
	typeId := typeIdentifier[T](key)

	// try to get entry from container
	lock.RLock()
	entries, entryExists := container[typeId]
	lock.RUnlock()

	if !entryExists {
		return make([]T, 0), nil
	}

	entriesCount := len(entries)

	if entriesCount == 0 {
		return make([]T, 0), nil
	}

	servicesArray := make([]T, entriesCount)

	for index := 0; index < entriesCount; index++ {
		e := entries[index].(entry[T])

		service, newCtx, updateEntry := e.load(ctx, contextValueId(typeId, index))

		if updateEntry {
			replaceEntry[T](typeId, index, e)
		}

		servicesArray[index] = service
		ctx = newCtx
	}

	return servicesArray, ctx
}
