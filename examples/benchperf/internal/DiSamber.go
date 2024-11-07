package internal

import (
	"github.com/samber/do/v2"
)

func BuildContainerDo() do.Injector {
	injector := do.New()
	do.ProvideTransient(injector, func(inj do.Injector) (*A, error) {
		return NewA(do.MustInvoke[*B](inj), do.MustInvoke[*C](inj)), nil
	})
	do.ProvideTransient(injector, func(inj do.Injector) (*B, error) {
		return NewB(do.MustInvoke[*D](inj), do.MustInvoke[*E](inj)), nil
	})
	do.ProvideTransient(injector, func(inj do.Injector) (*C, error) {
		return NewC(), nil
	})
	do.ProvideTransient(injector, func(inj do.Injector) (*D, error) {
		return NewD(do.MustInvoke[*F](inj), do.MustInvoke[H](inj)), nil
	})
	do.ProvideTransient(injector, func(inj do.Injector) (*E, error) {
		gs := []G{
			do.MustInvoke[*DGa](inj),
			do.MustInvoke[*Gb](inj),
			do.MustInvoke[*Gc](inj),
		}
		return NewE(gs), nil
	})
	do.Provide(injector, func(inj do.Injector) (*F, error) {
		return NewF(), nil
	})
	do.ProvideTransient(injector, func(inj do.Injector) (*Ga, error) {
		return NewGa(), nil
	})
	do.Provide(injector, func(inj do.Injector) (*Gb, error) {
		return NewGb(), nil
	})
	do.ProvideTransient(injector, func(inj do.Injector) (*Gc, error) {
		return NewGc(), nil
	})
	do.ProvideTransient(injector, func(inj do.Injector) (*DGa, error) {
		return NewDGa(do.MustInvoke[*Ga](inj)), nil
	})
	do.ProvideTransient(injector, func(inj do.Injector) (H, error) {
		return NewHr(), nil
	})
	return injector
}
