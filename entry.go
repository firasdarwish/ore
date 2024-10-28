package ore

import "context"

type (
	Initializer[T any] func(ctx context.Context) (T, context.Context)
)

type entry[T any] struct {
	anonymousInitializer *Initializer[T]
	creatorInstance      Creator[T]
	concrete             *T
	lifetime             Lifetime
}

func (i *entry[T]) load(ctx context.Context, ctxTidVal string) (T, context.Context, bool) {
	// try get concrete implementation
	if i.lifetime == Singleton && i.concrete != nil {
		return *i.concrete, ctx, false
	}

	// try get concrete from context scope
	if i.lifetime == Scoped {
		fromCtx, ok := ctx.Value(ctxTidVal).(T)
		if ok {
			return fromCtx, ctx, false
		}
	}

	var con T

	// first, try make concrete implementation from `anonymousInitializer`
	// if nil, try the concrete implementation `Creator`
	if i.anonymousInitializer != nil {
		con, ctx = (*i.anonymousInitializer)(ctx)
	} else {
		con, ctx = i.creatorInstance.New(ctx)
	}

	// if scoped, attach to the current context
	if i.lifetime == Scoped {
		ctx = context.WithValue(ctx, ctxTidVal, con)
	}

	// if was lazily-created, then attach the newly-created concrete implementation
	// to the entry
	if i.lifetime == Singleton {
		i.concrete = &con

		return con, ctx, true
	}

	return con, ctx, false
}
