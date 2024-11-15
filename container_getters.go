package ore

import (
	"context"
)

// GetFromContainer Retrieves an instance from the given container based on type and key (panics if no valid implementations)
func GetFromContainer[T any](con *Container, ctx context.Context, key ...KeyStringer) (T, context.Context) {
	pointerTypeName := getPointerTypeName[T]()
	typeID := getTypeID(pointerTypeName, key...)
	lastRegisteredResolver := con.getLastRegisteredResolver(typeID)
	if lastRegisteredResolver == nil { //not found, T is an alias

		con.lock.RLock()
		implementations, implExists := con.aliases[pointerTypeName]
		con.lock.RUnlock()

		if !implExists {
			panic(noValidImplementation[T]())
		}
		count := len(implementations)
		if count == 0 {
			panic(noValidImplementation[T]())
		}
		for i := count - 1; i >= 0; i-- {
			impl := implementations[i]
			typeID = getTypeID(impl, key...)
			lastRegisteredResolver = con.getLastRegisteredResolver(typeID)
			if lastRegisteredResolver != nil {
				break
			}
		}
	}
	if lastRegisteredResolver == nil {
		panic(noValidImplementation[T]())
	}
	concrete, ctx := lastRegisteredResolver.resolveService(con, ctx)
	return concrete.value.(T), ctx
}

// GetListFromContainer Retrieves a list of instances from the given container based on type and key
func GetListFromContainer[T any](con *Container, ctx context.Context, key ...KeyStringer) ([]T, context.Context) {
	inputPointerTypeName := getPointerTypeName[T]()

	con.lock.RLock()
	pointerTypeNames, implExists := con.aliases[inputPointerTypeName]
	con.lock.RUnlock()

	if implExists {
		pointerTypeNames = append(pointerTypeNames, inputPointerTypeName)
	} else {
		pointerTypeNames = []pointerTypeName{inputPointerTypeName}
	}

	servicesArray := []T{}

	for i := 0; i < len(pointerTypeNames); i++ {
		pointerTypeName := pointerTypeNames[i]
		// generate type identifier
		typeID := getTypeID(pointerTypeName, key...)

		// try to get service resolver from container
		con.lock.RLock()
		resolvers, resolverExists := con.resolvers[typeID]
		con.lock.RUnlock()

		if !resolverExists {
			continue
		}

		for index := 0; index < len(resolvers); index++ {
			resolver := resolvers[index]
			if resolver.isPlaceHolder() && !resolver.isScopedValueResolved(ctx) {
				//the resolver is a placeHolder and the placeHolder's value has not been provided
				//don't panic, just skip (don't add anything to the list)
				continue
			}
			con, newCtx := resolver.resolveService(con, ctx)
			servicesArray = append(servicesArray, con.value.(T))
			ctx = newCtx
		}
	}

	return servicesArray, ctx
}

// GetResolvedSingletonsFromContainer retrieves a list of Singleton instances that implement the [TInterface] from the given container.
// See [GetResolvedSingletons] for more information.
func GetResolvedSingletonsFromContainer[TInterface any](con *Container) []TInterface {
	con.lock.RLock()
	defer con.lock.RUnlock()

	list := []*concrete{}

	//filtering
	for _, resolvers := range con.resolvers {
		for _, resolver := range resolvers {
			con, isInvokedSingleton := resolver.getInvokedSingleton()
			if isInvokedSingleton {
				if _, ok := con.value.(TInterface); ok {
					list = append(list, con)
				}
			}
		}
	}

	return sortAndSelect[TInterface](list)
}

// TODO separate to a file
// GetResolvedSingletons retrieves a list of Singleton instances that implement the [TInterface].
// The returned instances are sorted by creation time (a.k.a the invocation order), the first one being the "most recently" created one.
// If an instance "A" depends on certain instances "B" and "C" then this function guarantee to return "B" and "C" before "A" in the list.
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
	return GetResolvedSingletonsFromContainer[TInterface](DefaultContainer)
}

// TODO separate to a file
// GetResolvedScopedInstances retrieves a list of Scoped instances that implement the [TInterface].
// The returned instances are sorted by creation time (a.k.a the invocation order), the first one being the most recently created one.
// If an instance "A" depends on certain instances "B" and "C" then this function guarantee to return "B" and "C" before "A" in the list.
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

	list := []*concrete{}

	//filtering
	for _, contextKey := range contextKeyRepository {
		con := ctx.Value(contextKey).(*concrete)
		if _, ok := con.value.(TInterface); ok {
			list = append(list, con)
		}
	}

	return sortAndSelect[TInterface](list)
}
