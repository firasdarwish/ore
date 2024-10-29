package ore

import (
	"context"
	"testing"
)

func TestRegisterLazyCreator(t *testing.T) {
	for _, registrationType := range types {
		clearAll()

		RegisterLazyCreator[someCounter](registrationType, &simpleCounter{})

		c, _ := Get[someCounter](context.Background())

		c.AddOne()
		c.AddOne()

		if got := c.GetCount(); got != 2 {
			t.Errorf("got %v, expected %v", got, 2)
		}
	}
}

func TestRegisterLazyCreatorNilFuncTransient(t *testing.T) {
	clearAll()
	defer mustHavePanicked(t)
	RegisterLazyCreator[someCounter](Transient, nil)
}

func TestRegisterLazyCreatorNilFuncScoped(t *testing.T) {
	clearAll()
	defer mustHavePanicked(t)
	RegisterLazyCreator[someCounter](Scoped, nil)
}

func TestRegisterLazyCreatorNilFuncSingleton(t *testing.T) {
	clearAll()
	defer mustHavePanicked(t)
	RegisterLazyCreator[someCounter](Singleton, nil)
}

func TestRegisterLazyCreatorMultipleImplementations(t *testing.T) {
	for _, registrationType := range types {
		clearAll()

		RegisterLazyCreator[someCounter](registrationType, &simpleCounter{})

		RegisterLazyCreator[someCounter](registrationType, &simpleCounter{})

		RegisterLazyCreator[someCounter](registrationType, &simpleCounter{})

		counters, _ := GetList[someCounter](context.Background())

		if got := len(counters); got != 3 {
			t.Errorf("got %v, expected %v", got, 3)
		}
	}
}

func TestRegisterLazyCreatorMultipleImplementationsKeyed(t *testing.T) {
	for _, registrationType := range types {
		clearAll()

		RegisterLazyCreator[someCounter](registrationType, &simpleCounter{}, "firas")

		RegisterLazyCreator[someCounter](registrationType, &simpleCounter{}, "firas")

		RegisterLazyCreator[someCounter](registrationType, &simpleCounter{})

		counters, _ := GetList[someCounter](context.Background(), "firas")

		if got := len(counters); got != 2 {
			t.Errorf("got %v, expected %v", got, 2)
		}
	}
}

func TestRegisterLazyCreatorSingletonState(t *testing.T) {
	var registrationType Lifetime = Singleton

	clearAll()

	RegisterLazyCreator[someCounter](registrationType, &simpleCounter{})

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

func TestRegisterLazyCreatorScopedState(t *testing.T) {
	var registrationType Lifetime = Scoped

	clearAll()

	RegisterLazyCreator[someCounter](registrationType, &simpleCounter{})

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

func TestRegisterLazyCreatorTransientState(t *testing.T) {
	var registrationType Lifetime = Transient

	clearAll()

	RegisterLazyCreator[someCounter](registrationType, &simpleCounter{})

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

func TestRegisterLazyCreatorNilKeyOnRegistering(t *testing.T) {
	clearAll()
	defer mustHavePanicked(t)
	RegisterLazyCreator[someCounter](Scoped, &simpleCounter{}, nil)
}

func TestRegisterLazyCreatorNilKeyOnGetting(t *testing.T) {
	clearAll()
	defer mustHavePanicked(t)
	RegisterLazyCreator[someCounter](Scoped, &simpleCounter{}, "firas")

	Get[someCounter](context.Background(), nil)
}

func TestRegisterLazyCreatorGeneric(t *testing.T) {
	for _, registrationType := range types {
		clearAll()

		RegisterLazyCreator[someCounterGeneric[uint]](registrationType, &counterGeneric[uint]{})

		c, _ := Get[someCounterGeneric[uint]](context.Background())

		c.Add(5)
		c.Add(5)

		if got := c.GetCount(); got != 10 {
			t.Errorf("got %v, expected %v", got, 10)
		}
	}
}

func TestRegisterLazyCreatorMultipleGenericImplementations(t *testing.T) {
	for _, registrationType := range types {
		clearAll()

		RegisterLazyCreator[someCounterGeneric[uint]](registrationType, &counterGeneric[uint]{})

		RegisterLazyCreator[someCounterGeneric[uint]](registrationType, &counterGeneric[uint]{})

		RegisterLazyCreator[someCounterGeneric[uint]](registrationType, &counterGeneric[uint]{})

		counters, _ := GetList[someCounterGeneric[uint]](context.Background())

		if got := len(counters); got != 3 {
			t.Errorf("got %v, expected %v", got, 3)
		}
	}
}

func TestRegisterLazyCreatorScopedNested(t *testing.T) {
	clearAll()

	RegisterLazyCreator[*a](Transient, &a{})

	RegisterLazyCreator[*c](Scoped, &c{})

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
