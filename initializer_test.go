package ore

import (
	"context"
	"testing"
)

func TestRegisterLazyFunc(t *testing.T) {
	for _, registrationType := range types {
		clearAll()

		RegisterLazyFunc[Counter](registrationType, func(ctx context.Context) (Counter, context.Context) {
			return &simpleCounter{}, ctx
		})

		c, _ := Get[Counter](context.Background())

		c.AddOne()
		c.AddOne()

		if got := c.GetCount(); got != 2 {
			t.Errorf("got %v, expected %v", got, 2)
		}
	}
}

func TestRegisterLazyFuncNilFuncTransient(t *testing.T) {
	clearAll()
	defer mustHavePanicked(t)
	RegisterLazyFunc[Counter](Transient, nil)
}

func TestRegisterLazyFuncNilFuncScoped(t *testing.T) {
	clearAll()
	defer mustHavePanicked(t)
	RegisterLazyFunc[Counter](Scoped, nil)
}

func TestRegisterLazyFuncNilFuncSingleton(t *testing.T) {
	clearAll()
	defer mustHavePanicked(t)
	RegisterLazyFunc[Counter](Singleton, nil)
}

func TestRegisterLazyFuncMultipleImplementations(t *testing.T) {
	for _, registrationType := range types {
		clearAll()

		RegisterLazyFunc[Counter](registrationType, func(ctx context.Context) (Counter, context.Context) {
			return &simpleCounter{}, ctx
		})

		RegisterLazyFunc[Counter](registrationType, func(ctx context.Context) (Counter, context.Context) {
			return &simpleCounter{}, ctx
		})

		RegisterLazyFunc[Counter](registrationType, func(ctx context.Context) (Counter, context.Context) {
			return &simpleCounter{}, ctx
		})

		counters, _ := GetList[Counter](context.Background())

		if got := len(counters); got != 3 {
			t.Errorf("got %v, expected %v", got, 3)
		}
	}
}

func TestRegisterLazyFuncMultipleImplementationsKeyed(t *testing.T) {
	for _, registrationType := range types {
		clearAll()

		RegisterLazyFunc[Counter](registrationType, func(ctx context.Context) (Counter, context.Context) {
			return &simpleCounter{}, ctx
		}, "firas")

		RegisterLazyFunc[Counter](registrationType, func(ctx context.Context) (Counter, context.Context) {
			return &simpleCounter{}, ctx
		}, "firas")

		RegisterLazyFunc[Counter](registrationType, func(ctx context.Context) (Counter, context.Context) {
			return &simpleCounter{}, ctx
		})

		counters, _ := GetList[Counter](context.Background(), "firas")

		if got := len(counters); got != 2 {
			t.Errorf("got %v, expected %v", got, 2)
		}
	}
}

func TestRegisterLazyFuncSingletonState(t *testing.T) {
	var registrationType Lifetime = Singleton

	clearAll()

	RegisterLazyFunc[Counter](registrationType, func(ctx context.Context) (Counter, context.Context) {
		return &simpleCounter{}, ctx
	})

	c, _ := Get[Counter](context.Background())
	c.AddOne()
	c.AddOne()

	c, _ = Get[Counter](context.Background())
	c.AddOne()

	c, _ = Get[Counter](context.Background())
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

	RegisterLazyFunc[Counter](registrationType, func(ctx context.Context) (Counter, context.Context) {
		return &simpleCounter{}, ctx
	})

	ctx := context.Background()

	c, ctx := Get[Counter](ctx)
	c.AddOne()
	c.AddOne()

	c, ctx = Get[Counter](ctx)
	c.AddOne()

	c, ctx = Get[Counter](ctx)
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

	RegisterLazyFunc[Counter](registrationType, func(ctx context.Context) (Counter, context.Context) {
		return &simpleCounter{}, ctx
	})

	ctx := context.Background()

	c, ctx := Get[Counter](ctx)
	c.AddOne()
	c.AddOne()

	c, ctx = Get[Counter](ctx)
	c.AddOne()

	c, ctx = Get[Counter](ctx)
	c.AddOne()
	c.AddOne()
	c.AddOne()

	if got := c.GetCount(); got != 3 {
		t.Errorf("got %v, expected %v", got, 3)
	}
}

func TestRegisterLazyFuncNilKeyOnRegistering(t *testing.T) {
	clearAll()
	defer mustHavePanicked(t)
	RegisterLazyFunc[Counter](Scoped, func(ctx context.Context) (Counter, context.Context) {
		return &simpleCounter{}, ctx
	}, nil)
}

func TestRegisterLazyFuncNilKeyOnGetting(t *testing.T) {
	clearAll()
	defer mustHavePanicked(t)
	RegisterLazyFunc[Counter](Scoped, func(ctx context.Context) (Counter, context.Context) {
		return &simpleCounter{}, ctx
	}, "firas")

	Get[Counter](context.Background(), nil)
}

func TestRegisterLazyFuncGeneric(t *testing.T) {
	for _, registrationType := range types {
		clearAll()

		RegisterLazyFunc[CounterGeneric[uint]](registrationType, func(ctx context.Context) (CounterGeneric[uint], context.Context) {
			return &counterGeneric[uint]{}, ctx
		})

		c, _ := Get[CounterGeneric[uint]](context.Background())

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

		RegisterLazyFunc[CounterGeneric[uint]](registrationType, func(ctx context.Context) (CounterGeneric[uint], context.Context) {
			return &counterGeneric[uint]{}, ctx
		})

		RegisterLazyFunc[CounterGeneric[uint]](registrationType, func(ctx context.Context) (CounterGeneric[uint], context.Context) {
			return &counterGeneric[uint]{}, ctx
		})

		RegisterLazyFunc[CounterGeneric[uint]](registrationType, func(ctx context.Context) (CounterGeneric[uint], context.Context) {
			return &counterGeneric[uint]{}, ctx
		})

		counters, _ := GetList[CounterGeneric[uint]](context.Background())

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

	a3, ctx := Get[*a](ctx)
	a3.C.Counter += 1

	if got := a2.C.Counter; got != 3 {
		t.Errorf("got %v, expected %v", got, 3)
	}
}
