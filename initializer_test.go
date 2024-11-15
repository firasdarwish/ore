package ore

import (
	"context"
	"github.com/firasdarwish/ore/internal/interfaces"
	"github.com/firasdarwish/ore/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegisterLazyFunc(t *testing.T) {
	for _, registrationType := range types {
		clearAll()

		RegisterLazyFunc[interfaces.SomeCounter](registrationType, func(ctx context.Context) (interfaces.SomeCounter, context.Context) {
			return &models.SimpleCounter{}, ctx
		})

		c, _ := Get[interfaces.SomeCounter](context.Background())

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
		RegisterLazyFunc[interfaces.SomeCounter](Transient, nil)
	})
}

func TestRegisterLazyFuncNilFuncScoped(t *testing.T) {
	clearAll()
	assert.Panics(t, func() {
		RegisterLazyFunc[interfaces.SomeCounter](Scoped, nil)
	})
}

func TestRegisterLazyFuncNilFuncSingleton(t *testing.T) {
	clearAll()
	assert.Panics(t, func() {
		RegisterLazyFunc[interfaces.SomeCounter](Singleton, nil)
	})
}

func TestRegisterLazyFuncMultipleImplementations(t *testing.T) {
	for _, registrationType := range types {
		clearAll()

		RegisterLazyFunc[interfaces.SomeCounter](registrationType, func(ctx context.Context) (interfaces.SomeCounter, context.Context) {
			return &models.SimpleCounter{}, ctx
		})

		RegisterLazyFunc[interfaces.SomeCounter](registrationType, func(ctx context.Context) (interfaces.SomeCounter, context.Context) {
			return &models.SimpleCounter{}, ctx
		})

		RegisterLazyFunc[interfaces.SomeCounter](registrationType, func(ctx context.Context) (interfaces.SomeCounter, context.Context) {
			return &models.SimpleCounter{}, ctx
		})

		counters, _ := GetList[interfaces.SomeCounter](context.Background())

		if got := len(counters); got != 3 {
			t.Errorf("got %v, expected %v", got, 3)
		}
	}
}

func TestRegisterLazyFuncMultipleImplementationsKeyed(t *testing.T) {
	for _, registrationType := range types {
		clearAll()

		RegisterLazyFunc[interfaces.SomeCounter](registrationType, func(ctx context.Context) (interfaces.SomeCounter, context.Context) {
			return &models.SimpleCounter{}, ctx
		}, "firas")

		RegisterLazyFunc[interfaces.SomeCounter](registrationType, func(ctx context.Context) (interfaces.SomeCounter, context.Context) {
			return &models.SimpleCounter{}, ctx
		}, "firas")

		RegisterLazyFunc[interfaces.SomeCounter](registrationType, func(ctx context.Context) (interfaces.SomeCounter, context.Context) {
			return &models.SimpleCounter{}, ctx
		})

		counters, _ := GetList[interfaces.SomeCounter](context.Background(), "firas")

		if got := len(counters); got != 2 {
			t.Errorf("got %v, expected %v", got, 2)
		}
	}
}

func TestRegisterLazyFuncSingletonState(t *testing.T) {
	var registrationType Lifetime = Singleton

	clearAll()

	RegisterLazyFunc[interfaces.SomeCounter](registrationType, func(ctx context.Context) (interfaces.SomeCounter, context.Context) {
		return &models.SimpleCounter{}, ctx
	})

	c, _ := Get[interfaces.SomeCounter](context.Background())
	c.AddOne()
	c.AddOne()

	c, _ = Get[interfaces.SomeCounter](context.Background())
	c.AddOne()

	c, _ = Get[interfaces.SomeCounter](context.Background())
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

	RegisterLazyFunc[interfaces.SomeCounter](registrationType, func(ctx context.Context) (interfaces.SomeCounter, context.Context) {
		return &models.SimpleCounter{}, ctx
	})

	ctx := context.Background()

	c, ctx := Get[interfaces.SomeCounter](ctx)
	c.AddOne()
	c.AddOne()

	c, ctx = Get[interfaces.SomeCounter](ctx)
	c.AddOne()

	c, _ = Get[interfaces.SomeCounter](ctx)
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

	RegisterLazyFunc[interfaces.SomeCounter](registrationType, func(ctx context.Context) (interfaces.SomeCounter, context.Context) {
		return &models.SimpleCounter{}, ctx
	})

	ctx := context.Background()

	c, ctx := Get[interfaces.SomeCounter](ctx)
	c.AddOne()
	c.AddOne()

	c, ctx = Get[interfaces.SomeCounter](ctx)
	c.AddOne()

	c, _ = Get[interfaces.SomeCounter](ctx)
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
		RegisterLazyFunc[interfaces.SomeCounter](Scoped, func(ctx context.Context) (interfaces.SomeCounter, context.Context) {
			return &models.SimpleCounter{}, ctx
		}, "", nil)
	})
}

func TestRegisterLazyFuncNilKeyOnGetting(t *testing.T) {
	clearAll()
	RegisterLazyFunc[interfaces.SomeCounter](Scoped, func(ctx context.Context) (interfaces.SomeCounter, context.Context) {
		return &models.SimpleCounter{}, ctx
	}, "firas")

	assert.Panics(t, func() {
		Get[interfaces.SomeCounter](context.Background(), "", nil)
	})
}

func TestRegisterLazyFuncGeneric(t *testing.T) {
	for _, registrationType := range types {
		clearAll()

		RegisterLazyFunc[interfaces.SomeCounterGeneric[uint]](registrationType, func(ctx context.Context) (interfaces.SomeCounterGeneric[uint], context.Context) {
			return &models.CounterGeneric[uint]{}, ctx
		})

		c, _ := Get[interfaces.SomeCounterGeneric[uint]](context.Background())

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

		RegisterLazyFunc[interfaces.SomeCounterGeneric[uint]](registrationType, func(ctx context.Context) (interfaces.SomeCounterGeneric[uint], context.Context) {
			return &models.CounterGeneric[uint]{}, ctx
		})

		RegisterLazyFunc[interfaces.SomeCounterGeneric[uint]](registrationType, func(ctx context.Context) (interfaces.SomeCounterGeneric[uint], context.Context) {
			return &models.CounterGeneric[uint]{}, ctx
		})

		RegisterLazyFunc[interfaces.SomeCounterGeneric[uint]](registrationType, func(ctx context.Context) (interfaces.SomeCounterGeneric[uint], context.Context) {
			return &models.CounterGeneric[uint]{}, ctx
		})

		counters, _ := GetList[interfaces.SomeCounterGeneric[uint]](context.Background())

		if got := len(counters); got != 3 {
			t.Errorf("got %v, expected %v", got, 3)
		}
	}
}

func TestRegisterLazyFuncScopedNested(t *testing.T) {
	clearAll()

	RegisterLazyFunc[*models.A](Transient, func(ctx context.Context) (*models.A, context.Context) {
		cc, ctx := Get[*models.C](ctx)
		return &models.A{
			C: cc,
		}, ctx
	})

	RegisterLazyFunc[*models.C](Scoped, func(ctx context.Context) (*models.C, context.Context) {
		return &models.C{}, ctx
	})

	ctx := context.Background()

	a1, ctx := Get[*models.A](ctx)
	a1.C.Counter += 1

	a2, ctx := Get[*models.A](ctx)
	a2.C.Counter += 1

	a3, _ := Get[*models.A](ctx)
	a3.C.Counter += 1

	if got := a2.C.Counter; got != 3 {
		t.Errorf("got %v, expected %v", got, 3)
	}
}
