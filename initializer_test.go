package ore

import (
	"context"
	"testing"

	"github.com/firasdarwish/ore/internal/interfaces"
	"github.com/firasdarwish/ore/internal/models"

	"github.com/stretchr/testify/assert"
)

func TestRegisterFuncSingleton(t *testing.T) {
	clearAll()

	RegisterFunc[interfaces.SomeCounter](Singleton, func(ctx context.Context) (interfaces.SomeCounter, context.Context) {
		return &models.SimpleCounter{}, ctx
	})

	ctx := context.Background()

	c, ctx := Get[interfaces.SomeCounter](ctx)
	c.AddOne()
	c.AddOne()
	c.AddOne()
	c.AddOne()
	c.AddOne()
	c.AddOne()

	RegisterFuncToContainer[interfaces.SomeCounter](DefaultContainer, Singleton, func(ctx context.Context) (interfaces.SomeCounter, context.Context) {
		return &models.SimpleCounter{}, ctx
	})

	c, ctx = Get[interfaces.SomeCounter](ctx)

	c.AddOne()
	c.AddOne()
	c.AddOne()

	RegisterKeyedFuncToContainer[interfaces.SomeCounter](DefaultContainer, Singleton, func(ctx context.Context) (interfaces.SomeCounter, context.Context) {
		return &models.SimpleCounter{}, ctx
	}, "firas")

	c, ctx = Get[interfaces.SomeCounter](ctx)

	if got := c.GetCount(); got != 3 {
		t.Errorf("got %v, expected %v", got, 3)
	}
}

func TestRegisterFuncScoped(t *testing.T) {
	clearAll()

	RegisterFunc[interfaces.SomeCounter](Scoped, func(ctx context.Context) (interfaces.SomeCounter, context.Context) {
		return &models.SimpleCounter{}, ctx
	})

	ctx := context.Background()

	c, ctx := Get[interfaces.SomeCounter](ctx)
	c.AddOne()
	c.AddOne()
	c.AddOne()
	c.AddOne()
	c.AddOne()
	c.AddOne()

	RegisterFuncToContainer[interfaces.SomeCounter](DefaultContainer, Scoped, func(ctx context.Context) (interfaces.SomeCounter, context.Context) {
		return &models.SimpleCounter{}, ctx
	})

	c, ctx = Get[interfaces.SomeCounter](ctx)

	c.AddOne()
	c.AddOne()
	c.AddOne()

	RegisterKeyedFuncToContainer[interfaces.SomeCounter](DefaultContainer, Scoped, func(ctx context.Context) (interfaces.SomeCounter, context.Context) {
		return &models.SimpleCounter{}, ctx
	}, "firas")

	c, ctx = Get[interfaces.SomeCounter](ctx)

	if got := c.GetCount(); got != 3 {
		t.Errorf("got %v, expected %v", got, 3)
	}
}

func TestRegisterFuncTransient(t *testing.T) {
	clearAll()

	RegisterFunc[interfaces.SomeCounter](Transient, func(ctx context.Context) (interfaces.SomeCounter, context.Context) {
		return &models.SimpleCounter{}, ctx
	})

	ctx := context.Background()

	c, ctx := Get[interfaces.SomeCounter](ctx)
	c.AddOne()
	c.AddOne()
	c.AddOne()
	c.AddOne()
	c.AddOne()
	c.AddOne()

	RegisterFuncToContainer[interfaces.SomeCounter](DefaultContainer, Transient, func(ctx context.Context) (interfaces.SomeCounter, context.Context) {
		return &models.SimpleCounter{}, ctx
	})

	c, ctx = Get[interfaces.SomeCounter](ctx)

	c.AddOne()

	RegisterKeyedFuncToContainer[interfaces.SomeCounter](DefaultContainer, Transient, func(ctx context.Context) (interfaces.SomeCounter, context.Context) {
		return &models.SimpleCounter{}, ctx
	}, "firas")

	c, ctx = Get[interfaces.SomeCounter](ctx)

	if got := c.GetCount(); got != 0 {
		t.Errorf("got %v, expected %v", got, 0)
	}
}

func TestRegisterFuncNilFuncTransient(t *testing.T) {
	clearAll()
	assert.Panics(t, func() {
		RegisterFunc[interfaces.SomeCounter](Transient, nil)
	})
}

func TestRegisterFuncNilFuncScoped(t *testing.T) {
	clearAll()
	assert.Panics(t, func() {
		RegisterFunc[interfaces.SomeCounter](Scoped, nil)
	})
}

func TestRegisterFuncNilFuncSingleton(t *testing.T) {
	clearAll()
	assert.Panics(t, func() {
		RegisterFunc[interfaces.SomeCounter](Singleton, nil)
	})
}

func TestRegisterFuncMultipleImplementations(t *testing.T) {
	for _, registrationType := range types {
		clearAll()

		RegisterFunc[interfaces.SomeCounter](registrationType, func(ctx context.Context) (interfaces.SomeCounter, context.Context) {
			return &models.SimpleCounter{}, ctx
		})

		RegisterFunc[interfaces.SomeCounter](registrationType, func(ctx context.Context) (interfaces.SomeCounter, context.Context) {
			return &models.SimpleCounter{}, ctx
		})

		RegisterFunc[interfaces.SomeCounter](registrationType, func(ctx context.Context) (interfaces.SomeCounter, context.Context) {
			return &models.SimpleCounter{}, ctx
		})

		counters, _ := GetList[interfaces.SomeCounter](context.Background())

		if got := len(counters); got != 3 {
			t.Errorf("got %v, expected %v", got, 3)
		}
	}
}

func TestRegisterFuncMultipleImplementationsKeyed(t *testing.T) {
	for _, registrationType := range types {
		clearAll()

		RegisterKeyedFunc[interfaces.SomeCounter](registrationType, func(ctx context.Context) (interfaces.SomeCounter, context.Context) {
			return &models.SimpleCounter{}, ctx
		}, "firas")

		RegisterKeyedFunc[interfaces.SomeCounter](registrationType, func(ctx context.Context) (interfaces.SomeCounter, context.Context) {
			return &models.SimpleCounter{}, ctx
		}, "firas")

		RegisterFunc[interfaces.SomeCounter](registrationType, func(ctx context.Context) (interfaces.SomeCounter, context.Context) {
			return &models.SimpleCounter{}, ctx
		})

		counters, _ := GetKeyedList[interfaces.SomeCounter](context.Background(), "firas")

		if got := len(counters); got != 2 {
			t.Errorf("got %v, expected %v", got, 2)
		}
	}
}

func TestRegisterFuncSingletonState(t *testing.T) {
	var registrationType Lifetime = Singleton

	clearAll()

	RegisterFunc[interfaces.SomeCounter](registrationType, func(ctx context.Context) (interfaces.SomeCounter, context.Context) {
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

func TestRegisterFuncScopedState(t *testing.T) {
	var registrationType Lifetime = Scoped

	clearAll()

	RegisterFunc[interfaces.SomeCounter](registrationType, func(ctx context.Context) (interfaces.SomeCounter, context.Context) {
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

func TestRegisterFuncTransientState(t *testing.T) {
	var registrationType Lifetime = Transient

	clearAll()

	RegisterFunc[interfaces.SomeCounter](registrationType, func(ctx context.Context) (interfaces.SomeCounter, context.Context) {
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

func TestRegisterFuncNilKeyOnGetting(t *testing.T) {
	clearAll()
	RegisterKeyedFunc[interfaces.SomeCounter](Scoped, func(ctx context.Context) (interfaces.SomeCounter, context.Context) {
		return &models.SimpleCounter{}, ctx
	}, "firas")

	assert.Panics(t, func() {
		GetKeyed[interfaces.SomeCounter](context.Background(), "")
	})
}

func TestRegisterFuncGeneric(t *testing.T) {
	for _, registrationType := range types {
		clearAll()

		RegisterFunc[interfaces.SomeCounterGeneric[uint]](registrationType, func(ctx context.Context) (interfaces.SomeCounterGeneric[uint], context.Context) {
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

func TestRegisterFuncMultipleGenericImplementations(t *testing.T) {
	for _, registrationType := range types {
		clearAll()

		RegisterFunc[interfaces.SomeCounterGeneric[uint]](registrationType, func(ctx context.Context) (interfaces.SomeCounterGeneric[uint], context.Context) {
			return &models.CounterGeneric[uint]{}, ctx
		})

		RegisterFunc[interfaces.SomeCounterGeneric[uint]](registrationType, func(ctx context.Context) (interfaces.SomeCounterGeneric[uint], context.Context) {
			return &models.CounterGeneric[uint]{}, ctx
		})

		RegisterFunc[interfaces.SomeCounterGeneric[uint]](registrationType, func(ctx context.Context) (interfaces.SomeCounterGeneric[uint], context.Context) {
			return &models.CounterGeneric[uint]{}, ctx
		})

		counters, _ := GetList[interfaces.SomeCounterGeneric[uint]](context.Background())

		if got := len(counters); got != 3 {
			t.Errorf("got %v, expected %v", got, 3)
		}
	}
}

func TestRegisterFuncScopedNested(t *testing.T) {
	clearAll()

	RegisterFunc[*models.A](Transient, func(ctx context.Context) (*models.A, context.Context) {
		cc, ctx := Get[*models.C](ctx)
		return &models.A{
			C: cc,
		}, ctx
	})

	RegisterFunc[*models.C](Scoped, func(ctx context.Context) (*models.C, context.Context) {
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
