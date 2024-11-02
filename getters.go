package ore

import (
	"context"
)

func getLastRegisteredResolver(typeId typeID) (serviceResolver, int) {
	// try to get service resolver from container
	lock.RLock()
	resolvers, resolverExists := container[typeId]
	lock.RUnlock()

	if !resolverExists {
		return nil, -1
	}

	count := len(resolvers)

	if count == 0 {
		return nil, -1
	}

	// index of the last implementation
	lastIndex := count - 1
	return resolvers[lastIndex], lastIndex
}

// Get Retrieves an instance based on type and key (panics if no valid implementations)
func Get[T any](ctx context.Context, key ...KeyStringer) (T, context.Context) {
	pointerTypeName := getPointerTypeName[T]()
	typeID := getTypeID(pointerTypeName, key)
	lastRegisteredResolver, lastIndex := getLastRegisteredResolver(typeID)
	if lastRegisteredResolver == nil { //not found, T is an alias

		lock.RLock()
		implementations, implExists := aliases[pointerTypeName]
		lock.RUnlock()

		if !implExists {
			panic(noValidImplementation[T]())
		}
		count := len(implementations)
		if count == 0 {
			panic(noValidImplementation[T]())
		}
		for i := count - 1; i >= 0; i-- {
			impl := implementations[i]
			typeID = getTypeID(impl, key)
			lastRegisteredResolver, lastIndex = getLastRegisteredResolver(typeID)
			if lastRegisteredResolver != nil {
				break
			}
		}
	}
	if lastRegisteredResolver == nil {
		panic(noValidImplementation[T]())
	}
	service, ctx := lastRegisteredResolver.resolveService(ctx, typeID, lastIndex)
	return service.(T), ctx
}

// GetList Retrieves a list of instances based on type and key
func GetList[T any](ctx context.Context, key ...KeyStringer) ([]T, context.Context) {
	inputPointerTypeName := getPointerTypeName[T]()

	lock.RLock()
	pointerTypeNames, implExists := aliases[inputPointerTypeName]
	lock.RUnlock()

	if implExists {
		pointerTypeNames = append(pointerTypeNames, inputPointerTypeName)
	} else {
		pointerTypeNames = []pointerTypeName{inputPointerTypeName}
	}

	servicesArray := []T{}

	for i := 0; i < len(pointerTypeNames); i++ {
		pointerTypeName := pointerTypeNames[i]
		// generate type identifier
		typeID := getTypeID(pointerTypeName, key)

		// try to get service resolver from container
		lock.RLock()
		resolvers, resolverExists := container[typeID]
		lock.RUnlock()

		if !resolverExists {
			continue
		}

		for index := 0; index < len(resolvers); index++ {
			resolver := resolvers[index]
			service, newCtx := resolver.resolveService(ctx, typeID, index)
			servicesArray = append(servicesArray, service.(T))
			ctx = newCtx
		}
	}

	return servicesArray, ctx
}

// GetResolvedSingletons retrieves a list of Singleton instances that implement the [TInterface].
// It would return only the instances which had been resolved. Other lazy implementations which have never been invoked will not be returned.
// This function is useful for cleaning operations.
//
// Example:
//
//	     disposableSingletons := ore.GetResolvedSingletons[Disposer]()
//		 for _, disposable := range disposableSingletons {
//		   disposable.Dispose()
//		 }
func GetResolvedSingletons[TInterface any]() []TInterface {
	lock.RLock()
	defer lock.RUnlock()

	result := []TInterface{}
	for _, resolvers := range container {
		for _, resolver := range resolvers {
			invokedValue, isInvokedSingleton := resolver.getInvokedSingleton()
			if isInvokedSingleton {
				if instance, ok := invokedValue.(TInterface); ok {
					result = append(result, instance)
				}
			}
		}
	}
	return result
}

// GetResolvedScopedInstances retrieves a list of Scoped instances that implement the [TInterface].
// It would return only the instances which had been resolved. Other lazy implementations which have never been invoked will not be returned.
// This function is useful for cleaning operations.
//
// Example:
//
//	     disposableInstances := ore.GetResolvedScopedInstances[Disposer](ctx)
//		 for _, disposable := range disposableInstances {
//		   disposable.Dispose()
//		 }
func GetResolvedScopedInstances[TInterface any](ctx context.Context) []TInterface {
	contextKeyRepository, ok := ctx.Value(contextKeysRepositoryID).(contextKeysRepository)
	if !ok {
		return []TInterface{}
	}
	result := []TInterface{}
	for _, contextKey := range contextKeyRepository {
		invokedValue := ctx.Value(contextKey)
		if instance, ok := invokedValue.(TInterface); ok {
			result = append(result, instance)
		}
	}
	return result
}
