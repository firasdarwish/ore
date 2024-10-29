package ore

import (
	"context"
)

// Get Retrieves an instance based on type and key (panics if no valid implementations)
func Get[T any](ctx context.Context, key ...KeyStringer) (T, context.Context) {
	// generate type identifier
	typeId := typeIdentifier[T](key)

	// try to get service resolver from container
	lock.RLock()
	resolvers, resolverExists := container[typeId]
	lock.RUnlock()

	if !resolverExists {
		panic(noValidImplementation[T]())
	}

	count := len(resolvers)

	if count == 0 {
		panic(noValidImplementation[T]())
	}

	// lastIndex of the last implementation
	lastIndex := count - 1
	lastRegisteredResolver := resolvers[lastIndex]
	service, ctx := lastRegisteredResolver.resolveService(ctx, typeId, lastIndex)
	return service.(T), ctx
}

// GetList Retrieves a list of instances based on type and key
func GetList[T any](ctx context.Context, key ...KeyStringer) ([]T, context.Context) {
	// generate type identifier
	typeId := typeIdentifier[T](key)

	// try to get service resolver from container
	lock.RLock()
	resolvers, resolverExists := container[typeId]
	lock.RUnlock()

	if !resolverExists {
		return make([]T, 0), nil
	}

	count := len(resolvers)

	if count == 0 {
		return make([]T, 0), nil
	}

	servicesArray := make([]T, count)

	for index := 0; index < count; index++ {
		resolver := resolvers[index]
		service, newCtx := resolver.resolveService(ctx, typeId, index)
		servicesArray[index] = service.(T)
		ctx = newCtx
	}

	return servicesArray, ctx
}
