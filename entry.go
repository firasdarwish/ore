package ore

import "context"

type (
	Initializer[T any] func(ctx context.Context) T
)

type entry[T any] struct {
	anonymousInitializer *Initializer[T]
	creatorInstance      Creator[T]
	concrete             *T
	registrationType     RegistrationType
}

func (i *entry[T]) load(ctx context.Context, ctxTidVal string) (T, context.Context) {

	// try get concrete implementation
	if i.registrationType == Singleton {
		if i.concrete != nil {
			return *i.concrete, ctx
		}
	}

	// try get concrete from context scope
	if i.registrationType == Scoped {
		fromCtx, ok := ctx.Value(ctxTidVal).(T)
		if ok {
			return fromCtx, ctx
		}
	}

	var con T

	// first, try make concrete implementation from `anonymousInitializer`
	// if nil, try the concrete implementation `Creator`
	if i.anonymousInitializer != nil {
		con = (*i.anonymousInitializer)(ctx)
	} else {
		con = i.creatorInstance.New(ctx)
	}

	// if scoped, attach to the current context
	if i.registrationType == Scoped {
		ctx = context.WithValue(ctx, ctxTidVal, con)
	}

	// if was lazily-created, then attach the newly-created concrete implementation
	// to the entry
	if i.registrationType == Singleton {
		i.concrete = &con
	}

	return con, ctx
}
