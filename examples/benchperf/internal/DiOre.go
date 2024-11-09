package internal

import (
	"context"

	"github.com/firasdarwish/ore"
)

func BuildContainerOre() bool {
	ore.RegisterLazyFunc(ore.Transient, func(ctx context.Context) (*A, context.Context) {
		b, ctx := ore.Get[*B](ctx)
		c, ctx := ore.Get[*C](ctx)
		return NewA(b, c), ctx
	})
	ore.RegisterLazyFunc(ore.Transient, func(ctx context.Context) (*B, context.Context) {
		d, ctx := ore.Get[*D](ctx)
		e, ctx := ore.Get[*E](ctx)
		return NewB(d, e), ctx
	})
	ore.RegisterLazyFunc(ore.Transient, func(ctx context.Context) (*C, context.Context) {
		return NewC(), ctx
	})
	ore.RegisterLazyFunc(ore.Transient, func(ctx context.Context) (*D, context.Context) {
		f, ctx := ore.Get[*F](ctx)
		h, ctx := ore.Get[H](ctx)
		return NewD(f, h), ctx
	})
	ore.RegisterLazyFunc(ore.Transient, func(ctx context.Context) (*E, context.Context) {
		gs, ctx := ore.GetList[G](ctx)
		return NewE(gs), ctx
	})
	ore.RegisterLazyFunc(ore.Singleton, func(ctx context.Context) (*F, context.Context) {
		return NewF(), ctx
	})
	ore.RegisterLazyFunc(ore.Transient, func(ctx context.Context) (*Ga, context.Context) {
		return NewGa(), ctx
	})
	ore.RegisterLazyFunc(ore.Singleton, func(ctx context.Context) (G, context.Context) {
		return NewGb(), ctx
	})
	ore.RegisterLazyFunc(ore.Transient, func(ctx context.Context) (G, context.Context) {
		return NewGc(), ctx
	})
	ore.RegisterLazyFunc(ore.Transient, func(ctx context.Context) (G, context.Context) {
		ga, ctx := ore.Get[*Ga](ctx)
		return NewDGa(ga), ctx
	})
	ore.RegisterLazyFunc(ore.Transient, func(ctx context.Context) (H, context.Context) {
		return NewHr(), ctx
	})
	return true
}
