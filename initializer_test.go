package ore

import (
	"context"
	"testing"
)

func TestRegisterLazyFunc(t *testing.T) {
	for _, registrationType := range types {
		clearAll()

		RegisterLazyFunc[Counter](registrationType, func(ctx context.Context) Counter {
			return &simpleCounter{}
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

		RegisterLazyFunc[Counter](registrationType, func(ctx context.Context) Counter {
			return &simpleCounter{}
		})

		RegisterLazyFunc[Counter](registrationType, func(ctx context.Context) Counter {
			return &simpleCounter{}
		})

		RegisterLazyFunc[Counter](registrationType, func(ctx context.Context) Counter {
			return &simpleCounter{}
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

		RegisterLazyFunc[Counter](registrationType, func(ctx context.Context) Counter {
			return &simpleCounter{}
		}, "firas")

		RegisterLazyFunc[Counter](registrationType, func(ctx context.Context) Counter {
			return &simpleCounter{}
		}, "firas")

		RegisterLazyFunc[Counter](registrationType, func(ctx context.Context) Counter {
			return &simpleCounter{}
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

	RegisterLazyFunc[Counter](registrationType, func(ctx context.Context) Counter {
		return &simpleCounter{}
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

	RegisterLazyFunc[Counter](registrationType, func(ctx context.Context) Counter {
		return &simpleCounter{}
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

	RegisterLazyFunc[Counter](registrationType, func(ctx context.Context) Counter {
		return &simpleCounter{}
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
	RegisterLazyFunc[Counter](Scoped, func(ctx context.Context) Counter {
		return &simpleCounter{}
	}, nil)
}

func TestRegisterLazyFuncNilKeyOnGetting(t *testing.T) {
	clearAll()
	defer mustHavePanicked(t)
	RegisterLazyFunc[Counter](Scoped, func(ctx context.Context) Counter {
		return &simpleCounter{}
	}, "firas")

	Get[Counter](context.Background(), nil)
}

func TestRegisterLazyFuncGeneric(t *testing.T) {
	for _, registrationType := range types {
		clearAll()

		RegisterLazyFunc[CounterGeneric[uint]](registrationType, func(ctx context.Context) CounterGeneric[uint] {
			return &counterGeneric[uint]{}
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

		RegisterLazyFunc[CounterGeneric[uint]](registrationType, func(ctx context.Context) CounterGeneric[uint] {
			return &counterGeneric[uint]{}
		})

		RegisterLazyFunc[CounterGeneric[uint]](registrationType, func(ctx context.Context) CounterGeneric[uint] {
			return &counterGeneric[uint]{}
		})

		RegisterLazyFunc[CounterGeneric[uint]](registrationType, func(ctx context.Context) CounterGeneric[uint] {
			return &counterGeneric[uint]{}
		})

		counters, _ := GetList[CounterGeneric[uint]](context.Background())

		if got := len(counters); got != 3 {
			t.Errorf("got %v, expected %v", got, 3)
		}
	}
}
