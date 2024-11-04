package ore

import (
	"context"
	"time"
)

type (
	Initializer[T any] func(ctx context.Context) (T, context.Context)
)

type serviceResolver interface {
	resolveService(ctx context.Context, typeId typeID, index int) (*concrete, context.Context)
	//return the invoked singleton value, or false if the resolver is not a singleton or has not been invoked
	getInvokedSingleton() (con *concrete, isInvokedSingleton bool)
}

type serviceResolverImpl[T any] struct {
	anonymousInitializer *Initializer[T]
	creatorInstance      Creator[T]
	singletonConcrete    *concrete
	lifetime             Lifetime
}

// make sure that the `serviceResolverImpl` struct implements the `serviceResolver` interface
var _ serviceResolver = serviceResolverImpl[any]{}

func (this serviceResolverImpl[T]) resolveService(ctx context.Context, typeID typeID, index int) (*concrete, context.Context) {
	// try get concrete implementation
	if this.lifetime == Singleton && this.singletonConcrete != nil {
		return this.singletonConcrete, ctx
	}

	ctxKey := contextKey{typeID, index}

	// try get concrete from context scope
	if this.lifetime == Scoped {
		scopedConcrete, ok := ctx.Value(ctxKey).(*concrete)
		if ok {
			return scopedConcrete, ctx
		}
	}

	var concreteValue T

	// first, try make concrete implementation from `anonymousInitializer`
	// if nil, try the concrete implementation `Creator`
	if this.anonymousInitializer != nil {
		concreteValue, ctx = (*this.anonymousInitializer)(ctx)
	} else {
		concreteValue, ctx = this.creatorInstance.New(ctx)
	}

	con := &concrete{
		value:     concreteValue,
		lifetime:  this.lifetime,
		createdAt: time.Now(),
	}

	// if scoped, attach to the current context
	if this.lifetime == Scoped {
		ctx = context.WithValue(ctx, ctxKey, con)
		ctx = addToContextKeysRepository(ctx, ctxKey)
	}

	// if was lazily-created, then attach the newly-created concrete implementation
	// to the service resolver
	if this.lifetime == Singleton {
		this.singletonConcrete = con
		replaceServiceResolver(typeID, index, this)
		return con, ctx
	}

	return con, ctx
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
