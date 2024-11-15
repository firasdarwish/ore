package internal

import (
	"context"

	"github.com/firasdarwish/ore"
)

func BuildContainerOre(disableValidation bool) *ore.Container {
	container := ore.NewContainer()
	RegisterToOreContainer(container)
	container.DisableValidation = disableValidation
	return container
}

func RegisterToOreContainer(container *ore.Container) {
	ore.RegisterFuncToContainer(container, ore.Transient, func(ctx context.Context) (*A, context.Context) {
		b, ctx := ore.GetFromContainer[*B](container, ctx)
		c, ctx := ore.GetFromContainer[*C](container, ctx)
		return NewA(b, c), ctx
	})
	ore.RegisterFuncToContainer(container, ore.Transient, func(ctx context.Context) (*B, context.Context) {
		d, ctx := ore.GetFromContainer[*D](container, ctx)
		e, ctx := ore.GetFromContainer[*E](container, ctx)
		return NewB(d, e), ctx
	})
	ore.RegisterFuncToContainer(container, ore.Transient, func(ctx context.Context) (*C, context.Context) {
		return NewC(), ctx
	})
	ore.RegisterFuncToContainer(container, ore.Transient, func(ctx context.Context) (*D, context.Context) {
		f, ctx := ore.GetFromContainer[*F](container, ctx)
		h, ctx := ore.GetFromContainer[H](container, ctx)
		return NewD(f, h), ctx
	})
	ore.RegisterFuncToContainer(container, ore.Transient, func(ctx context.Context) (*E, context.Context) {
		gs, ctx := ore.GetListFromContainer[G](container, ctx)
		return NewE(gs), ctx
	})
	ore.RegisterFuncToContainer(container, ore.Singleton, func(ctx context.Context) (*F, context.Context) {
		return NewF(), ctx
	})
	ore.RegisterFuncToContainer(container, ore.Transient, func(ctx context.Context) (*Ga, context.Context) {
		return NewGa(), ctx
	})
	ore.RegisterFuncToContainer(container, ore.Singleton, func(ctx context.Context) (G, context.Context) {
		return NewGb(), ctx
	})
	ore.RegisterFuncToContainer(container, ore.Transient, func(ctx context.Context) (G, context.Context) {
		return NewGc(), ctx
	})
	ore.RegisterFuncToContainer(container, ore.Transient, func(ctx context.Context) (G, context.Context) {
		ga, ctx := ore.GetFromContainer[*Ga](container, ctx)
		return NewDGa(ga), ctx
	})
	ore.RegisterFuncToContainer(container, ore.Transient, func(ctx context.Context) (H, context.Context) {
		return NewHr(), ctx
	})
}
