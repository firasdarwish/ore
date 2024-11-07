package ore

import (
	"container/list"
	"context"
	"fmt"
	"time"
)

type (
	Initializer[T any] func(ctx context.Context) (T, context.Context)
)

type serviceResolver interface {
	fmt.Stringer
	resolveService(ctx context.Context) (*concrete, context.Context)
	//return the invoked singleton value, or false if the resolver is not a singleton or has not been invoked
	getInvokedSingleton() (con *concrete, isInvokedSingleton bool)
}

type resolverMetadata struct {
	id       contextKey
	lifetime Lifetime
}

type serviceResolverImpl[T any] struct {
	resolverMetadata
	anonymousInitializer *Initializer[T]
	creatorInstance      Creator[T]
	singletonConcrete    *concrete
}

// resolversChain is a linkedList[resolverMetadata], describing a dependencies chain which a resolver has to invoke other resolvers to resolve its dependencies.
// Before a resolver creates a new concrete value it would be registered to the resolversChain.
// Once the concrete is resolved (with help of other resolvers), then it would be removed from the chain.
//
// While a Resolver forms a tree with other dependent resolvers.
//
// Example:
//
//	A calls B and C; B calls D; C calls E.
//
// then resolversChain is a "path" in the tree from the root to one of the bottom.
//
// Example:
//
//	A -> B -> D or A -> C -> E
//
// The resolversChain is stored in the context. Analyze the chain will help to
//
//   - (1) detect the invocation level
//   - (2) detect cyclic dependencies
//   - (3) detect lifetime misalignment (when a service of longer lifetime depends on a service of shorter lifetime)
type resolversChain = *list.List

// make sure that the `serviceResolverImpl` struct implements the `serviceResolver` interface
var _ serviceResolver = serviceResolverImpl[any]{}

func (this serviceResolverImpl[T]) resolveService(ctx context.Context) (*concrete, context.Context) {
	// try get concrete implementation
	if this.lifetime == Singleton && this.singletonConcrete != nil {
		return this.singletonConcrete, ctx
	}

	// try get concrete from context scope
	if this.lifetime == Scoped {
		scopedConcrete, ok := ctx.Value(this.id).(*concrete)
		if ok {
			return scopedConcrete, ctx
		}
	}

	// this resolver is about to create a new concrete value, we have to put it to the resolversChain until the creation done

	// get the current currentChain from the context
	var currentChain resolversChain
	var marker *list.Element
	if !DisableValidation {
		untypedCurrentChain := ctx.Value(contextKeyResolversChain)
		if untypedCurrentChain == nil {
			currentChain = list.New()
			ctx = context.WithValue(ctx, contextKeyResolversChain, currentChain)
		} else {
			currentChain = untypedCurrentChain.(resolversChain)
		}

		// push this newest resolver to the resolversChain
		marker = appendResolver(currentChain, this.resolverMetadata)
	}
	var concreteValue T
	createdAt := time.Now()

	// first, try make concrete implementation from `anonymousInitializer`
	// if nil, try the concrete implementation `Creator`
	if this.anonymousInitializer != nil {
		concreteValue, ctx = (*this.anonymousInitializer)(ctx)
	} else {
		concreteValue, ctx = this.creatorInstance.New(ctx)
	}

	invocationLevel := 0
	if !DisableValidation {
		invocationLevel = currentChain.Len()

		// the concreteValue is created, we must to remove it from the resolversChain so that downstream resolvers (meaning the future resolvers) won't link to it
		currentChain.Remove(marker)
	}

	con := &concrete{
		value:           concreteValue,
		lifetime:        this.lifetime,
		createdAt:       createdAt,
		invocationLevel: invocationLevel,
	}

	// if scoped, attach to the current context
	if this.lifetime == Scoped {
		ctx = context.WithValue(ctx, this.id, con)
		ctx = addToContextKeysRepository(ctx, this.id)
	}

	// if was lazily-created, then attach the newly-created concrete implementation
	// to the service resolver
	if this.lifetime == Singleton {
		this.singletonConcrete = con
		replaceServiceResolver(this)
		return con, ctx
	}

	return con, ctx
}

// appendToResolversChain push the given resolver to the Back of the ResolversChain.
// `marker.previous` refers to the calling (parent) resolver
func appendResolver(chain resolversChain, currentResolver resolverMetadata) (marker *list.Element) {
	if chain.Len() != 0 {
		//detect lifetime misalignment
		lastElem := chain.Back()
		lastResolver := lastElem.Value.(resolverMetadata)
		if lastResolver.lifetime > currentResolver.lifetime {
			panic(lifetimeMisalignment(lastResolver, currentResolver))
		}

		//detect cyclic dependencies
		for e := chain.Back(); e != nil; e = e.Prev() {
			if e.Value.(resolverMetadata).id == currentResolver.id {
				panic(cyclicDependency(currentResolver))
			}
		}
	}
	marker = chain.PushBack(currentResolver) // `marker.previous` refers to the calling (parent) resolver
	return marker
}

func (this serviceResolverImpl[T]) getInvokedSingleton() (con *concrete, isInvokedSingleton bool) {
	if this.lifetime == Singleton && this.singletonConcrete != nil {
		return this.singletonConcrete, true
	}
	return nil, false
}

func addToContextKeysRepository(ctx context.Context, newContextKey contextKey) context.Context {
	repository, ok := ctx.Value(contextKeysRepositoryID).(contextKeysRepository)
	if ok {
		repository = append(repository, newContextKey)
	} else {
		repository = contextKeysRepository{newContextKey}
	}
	return context.WithValue(ctx, contextKeysRepositoryID, repository)
}

func (this resolverMetadata) String() string {
	return fmt.Sprintf("Resolver(%s, type={%s}, key='%s')", this.lifetime, getUnderlyingTypeName(this.id.pointerTypeName), this.id.oreKey)
}

// func toString(resolversChain resolversChain) string {
// 	var sb string
// 	for e := resolversChain.Front(); e != nil; e = e.Next() {
// 		sb = fmt.Sprintf("%s%s\n", sb, e.Value.(resolverMetadata).String())
// 	}
// 	return sb
// }
