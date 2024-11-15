package ore

import "context"

// GetResolvedSingletonsFromContainer retrieves a list of Singleton instances that implement the [TInterface] from the given container.
// See [GetResolvedSingletons] for more information.
func GetResolvedSingletonsFromContainer[TInterface any](con *Container) []TInterface {
	return getResolvedSingletonsFromContainer[TInterface](con)
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
	return getResolvedSingletonsFromContainer[TInterface](DefaultContainer)
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
