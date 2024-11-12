package ore

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegisterLazyFunc(t *testing.T) {
	for _, registrationType := range types {
		clearAll()

		RegisterLazyFunc[someCounter](registrationType, func(ctx context.Context) (someCounter, context.Context) {
			return &simpleCounter{}, ctx
		})

		c, _ := Get[someCounter](context.Background())

		c.AddOne()
		c.AddOne()

		if got := c.GetCount(); got != 2 {
			t.Errorf("got %v, expected %v", got, 2)
		}
	}
}

func TestRegisterLazyFuncNilFuncTransient(t *testing.T) {
	clearAll()
	assert.Panics(t, func() {
		RegisterLazyFunc[someCounter](Transient, nil)
	})
}

func TestRegisterLazyFuncNilFuncScoped(t *testing.T) {
	clearAll()
	assert.Panics(t, func() {
		RegisterLazyFunc[someCounter](Scoped, nil)
	})
}

func TestRegisterLazyFuncNilFuncSingleton(t *testing.T) {
	clearAll()
	assert.Panics(t, func() {
		RegisterLazyFunc[someCounter](Singleton, nil)
	})
}

func TestRegisterLazyFuncMultipleImplementations(t *testing.T) {
	for _, registrationType := range types {
		clearAll()

		RegisterLazyFunc[someCounter](registrationType, func(ctx context.Context) (someCounter, context.Context) {
			return &simpleCounter{}, ctx
		})

		RegisterLazyFunc[someCounter](registrationType, func(ctx context.Context) (someCounter, context.Context) {
			return &simpleCounter{}, ctx
		})

		RegisterLazyFunc[someCounter](registrationType, func(ctx context.Context) (someCounter, context.Context) {
			return &simpleCounter{}, ctx
		})

		counters, _ := GetList[someCounter](context.Background())

		if got := len(counters); got != 3 {
			t.Errorf("got %v, expected %v", got, 3)
		}
	}
}

func TestRegisterLazyFuncMultipleImplementationsKeyed(t *testing.T) {
	for _, registrationType := range types {
		clearAll()

		RegisterLazyFunc[someCounter](registrationType, func(ctx context.Context) (someCounter, context.Context) {
			return &simpleCounter{}, ctx
		}, "firas")

		RegisterLazyFunc[someCounter](registrationType, func(ctx context.Context) (someCounter, context.Context) {
			return &simpleCounter{}, ctx
		}, "firas")

		RegisterLazyFunc[someCounter](registrationType, func(ctx context.Context) (someCounter, context.Context) {
			return &simpleCounter{}, ctx
		})

		counters, _ := GetList[someCounter](context.Background(), "firas")

		if got := len(counters); got != 2 {
			t.Errorf("got %v, expected %v", got, 2)
		}
	}
}

func TestRegisterLazyFuncSingletonState(t *testing.T) {
	var registrationType Lifetime = Singleton

	clearAll()

	RegisterLazyFunc[someCounter](registrationType, func(ctx context.Context) (someCounter, context.Context) {
		return &simpleCounter{}, ctx
	})

	c, _ := Get[someCounter](context.Background())
	c.AddOne()
	c.AddOne()

	c, _ = Get[someCounter](context.Background())
	c.AddOne()

	c, _ = Get[someCounter](context.Background())
	c.AddOne()
	c.AddOne()
	c.AddOne()

	if got := c.GetCount(); got != 6 {
		t.Errorf("got %v, expected %v", got, 6)
	}
}

func TestRegisterLazyFuncScopedState(t *testing.T) {
	var registrationType Lifetime = Scoped

	clearAll()

	RegisterLazyFunc[someCounter](registrationType, func(ctx context.Context) (someCounter, context.Context) {
		return &simpleCounter{}, ctx
	})

	ctx := context.Background()

	c, ctx := Get[someCounter](ctx)
	c.AddOne()
	c.AddOne()

	c, ctx = Get[someCounter](ctx)
	c.AddOne()

	c, _ = Get[someCounter](ctx)
	c.AddOne()
	c.AddOne()
	c.AddOne()

	if got := c.GetCount(); got != 6 {
		t.Errorf("got %v, expected %v", got, 6)
	}
}

func TestRegisterLazyFuncTransientState(t *testing.T) {
	var registrationType Lifetime = Transient

	clearAll()

	RegisterLazyFunc[someCounter](registrationType, func(ctx context.Context) (someCounter, context.Context) {
		return &simpleCounter{}, ctx
	})

	ctx := context.Background()

	c, ctx := Get[someCounter](ctx)
	c.AddOne()
	c.AddOne()

	c, ctx = Get[someCounter](ctx)
	c.AddOne()

	c, _ = Get[someCounter](ctx)
	c.AddOne()
	c.AddOne()
	c.AddOne()

	if got := c.GetCount(); got != 3 {
		t.Errorf("got %v, expected %v", got, 3)
	}
}

func TestRegisterLazyFuncNilKeyOnRegistering(t *testing.T) {
	clearAll()
	assert.Panics(t, func() {
		RegisterLazyFunc[someCounter](Scoped, func(ctx context.Context) (someCounter, context.Context) {
			return &simpleCounter{}, ctx
		}, "", nil)
	})
}

func TestRegisterLazyFuncNilKeyOnGetting(t *testing.T) {
	clearAll()
	RegisterLazyFunc[someCounter](Scoped, func(ctx context.Context) (someCounter, context.Context) {
		return &simpleCounter{}, ctx
	}, "firas")

	assert.Panics(t, func() {
		Get[someCounter](context.Background(), "", nil)
	})
}

func TestRegisterLazyFuncGeneric(t *testing.T) {
	for _, registrationType := range types {
		clearAll()

		RegisterLazyFunc[someCounterGeneric[uint]](registrationType, func(ctx context.Context) (someCounterGeneric[uint], context.Context) {
			return &counterGeneric[uint]{}, ctx
		})

		c, _ := Get[someCounterGeneric[uint]](context.Background())

		c.Add(5)
		c.Add(5)

		if got := c.GetCount(); got != 10 {
			t.Errorf("got %v, expected %v", got, 10)
		}
	}
}

func TestRegisterLazyFuncMultipleGenericImplementations(t *testing.T) {
	for _, registrationType := range types {
		clearAll()

		RegisterLazyFunc[someCounterGeneric[uint]](registrationType, func(ctx context.Context) (someCounterGeneric[uint], context.Context) {
			return &counterGeneric[uint]{}, ctx
		})

		RegisterLazyFunc[someCounterGeneric[uint]](registrationType, func(ctx context.Context) (someCounterGeneric[uint], context.Context) {
			return &counterGeneric[uint]{}, ctx
		})

		RegisterLazyFunc[someCounterGeneric[uint]](registrationType, func(ctx context.Context) (someCounterGeneric[uint], context.Context) {
			return &counterGeneric[uint]{}, ctx
		})

		counters, _ := GetList[someCounterGeneric[uint]](context.Background())

		if got := len(counters); got != 3 {
			t.Errorf("got %v, expected %v", got, 3)
		}
	}
}

func TestRegisterLazyFuncScopedNested(t *testing.T) {
	clearAll()

	RegisterLazyFunc[*a](Transient, func(ctx context.Context) (*a, context.Context) {
		cc, ctx := Get[*c](ctx)
		return &a{
			C: cc,
		}, ctx
	})

	RegisterLazyFunc[*c](Scoped, func(ctx context.Context) (*c, context.Context) {
		return &c{}, ctx
	})

	ctx := context.Background()

	a1, ctx := Get[*a](ctx)
	a1.C.Counter += 1

	a2, ctx := Get[*a](ctx)
	a2.C.Counter += 1

	a3, _ := Get[*a](ctx)
	a3.C.Counter += 1

	if got := a2.C.Counter; got != 3 {
		t.Errorf("got %v, expected %v", got, 3)
	}
}
