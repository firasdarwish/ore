package ore

import (
	"context"
	"sort"
)

func (this *Container) getLastRegisteredResolver(typeID typeID) serviceResolver {
	// try to get service resolver from container
	this.lock.RLock()
	resolvers, resolverExists := this.resolvers[typeID]
	this.lock.RUnlock()

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

// sortAndSelect sorts concretes by invocation order and return its value.
func sortAndSelect[TInterface any](list []*concrete) []TInterface {
	//sorting
	sort.Slice(list, func(i, j int) bool {
		return list[i].invocationTime.After(list[j].invocationTime) ||
			(list[i].invocationTime == list[j].invocationTime &&
				list[i].invocationLevel > list[j].invocationLevel)
	})

	//selecting
	result := make([]TInterface, len(list))
	for i := 0; i < len(list); i++ {
		result[i] = list[i].value.(TInterface)
	}
	return result
}

func getFromContainer[T any](con *Container, ctx context.Context, key KeyStringer) (T, context.Context) {
	pointerTypeName := getPointerTypeName[T]()
	typeID := getTypeID(pointerTypeName, key)
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
			typeID = getTypeID(impl, key)
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

func getListFromContainer[T any](con *Container, ctx context.Context, key KeyStringer) ([]T, context.Context) {
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
		typeID := getTypeID(pointerTypeName, key)

		// try to get service resolver from container
		con.lock.RLock()
		resolvers, resolverExists := con.resolvers[typeID]
		con.lock.RUnlock()

		if !resolverExists {
			continue
		}

		for index := 0; index < len(resolvers); index++ {
			resolver := resolvers[index]
			if resolver.isPlaceholder() && !resolver.isScopedValueResolved(ctx) {
				//the resolver is a placeholder and the placeholder's value has not been provided
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

func getResolvedSingletonsFromContainer[TInterface any](con *Container) []TInterface {
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
