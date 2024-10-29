package ore

import "context"

type (
	Initializer[T any] func(ctx context.Context) (T, context.Context)
)

type serviceResolver interface {
	resolveService(ctx context.Context, typeId typeID, index int) (any, context.Context)
}

type serviceResolverImpl[T any] struct {
	anonymousInitializer *Initializer[T]
	creatorInstance      Creator[T]
	singletonConcrete    *T
	lifetime             Lifetime
}

//make sure that the `serviceResolverImpl` struct implements the `serviceResolver` interface
var _ serviceResolver = serviceResolverImpl[any]{}

func (this serviceResolverImpl[T]) resolveService(ctx context.Context, typeId typeID, index int) (any, context.Context) {

	ctxTidVal := getContextValueID(typeId, index)

	// try get concrete implementation
	if this.lifetime == Singleton && this.singletonConcrete != nil {
		return *this.singletonConcrete, ctx
	}

	// try get concrete from context scope
	if this.lifetime == Scoped {
		scopedConcrete, ok := ctx.Value(ctxTidVal).(T)
		if ok {
			return scopedConcrete, ctx
		}
	}

	var con T

	// first, try make concrete implementation from `anonymousInitializer`
	// if nil, try the concrete implementation `Creator`
	if this.anonymousInitializer != nil {
		con, ctx = (*this.anonymousInitializer)(ctx)
	} else {
		con, ctx = this.creatorInstance.New(ctx)
	}

	// if scoped, attach to the current context
	if this.lifetime == Scoped {
		ctx = context.WithValue(ctx, ctxTidVal, con)
	}

	// if was lazily-created, then attach the newly-created concrete implementation
	// to the service resolver
	if this.lifetime == Singleton {
		this.singletonConcrete = &con
		replaceServiceResolver(typeId, index, this)
		return con, ctx
	}

	return con, ctx
}
