package ore

import "sort"

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
