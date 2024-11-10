package ore

import (
	"context"
	"sort"
)

func getLastRegisteredResolver(typeID typeID) serviceResolver {
	// try to get service resolver from container
	lock.RLock()
	resolvers, resolverExists := container[typeID]
	lock.RUnlock()

	if !resolverExists {
		return nil
	}

	count := len(resolvers)

	if count == 0 {
		return nil
	}

	// index of the last implementation
	lastIndex := count - 1
	return resolvers[lastIndex]
}

// Get Retrieves an instance based on type and key (panics if no valid implementations)
func Get[T any](ctx context.Context, key ...KeyStringer) (T, context.Context) {
	pointerTypeName := getPointerTypeName[T]()
	typeID := getTypeID(pointerTypeName, key)
	lastRegisteredResolver := getLastRegisteredResolver(typeID)
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
			lastRegisteredResolver = getLastRegisteredResolver(typeID)
			if lastRegisteredResolver != nil {
				break
			}
		}
	}
	if lastRegisteredResolver == nil {
		panic(noValidImplementation[T]())
	}
	con, ctx := lastRegisteredResolver.resolveService(ctx)
	return con.value.(T), ctx
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

	var servicesArray []T

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
			con, newCtx := resolver.resolveService(ctx)
			servicesArray = append(servicesArray, con.value.(T))
			ctx = newCtx
		}
	}

	return servicesArray, ctx
}

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
	lock.RLock()
	defer lock.RUnlock()

	var list []*concrete

	//filtering
	for _, resolvers := range container {
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

	var list []*concrete

	//filtering
	for _, contextKey := range contextKeyRepository {
		con := ctx.Value(contextKey).(*concrete)
		if _, ok := con.value.(TInterface); ok {
			list = append(list, con)
		}
	}

	return sortAndSelect[TInterface](list)
}

// sortAndSelect sorts concretes by invocation order and return its value.
func sortAndSelect[TInterface any](list []*concrete) []TInterface {
	//sorting
	sort.Slice(list, func(i, j int) bool {
		return list[i].invocationOrder > list[j].invocationOrder ||
			(list[i].invocationOrder == list[j].invocationOrder && list[i].invocationLevel > list[j].invocationLevel)
	})

	//selecting
	result := make([]TInterface, len(list))
	for i := 0; i < len(list); i++ {
		result[i] = list[i].value.(TInterface)
	}
	return result
}
