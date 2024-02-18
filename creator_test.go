package ore

import (
	"context"
	"testing"
)

func TestRegisterLazyCreator(t *testing.T) {
	for _, registrationType := range types {
		clearAll()

		RegisterLazyCreator[Counter](registrationType, &simpleCounter{})

		c, _ := Get[Counter](context.Background())

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
	RegisterLazyCreator[Counter](Transient, nil)
}

func TestRegisterLazyCreatorNilFuncScoped(t *testing.T) {
	clearAll()
	defer mustHavePanicked(t)
	RegisterLazyCreator[Counter](Scoped, nil)
}

func TestRegisterLazyCreatorNilFuncSingleton(t *testing.T) {
	clearAll()
	defer mustHavePanicked(t)
	RegisterLazyCreator[Counter](Singleton, nil)
}

func TestRegisterLazyCreatorMultipleImplementations(t *testing.T) {
	for _, registrationType := range types {
		clearAll()

		RegisterLazyCreator[Counter](registrationType, &simpleCounter{})

		RegisterLazyCreator[Counter](registrationType, &simpleCounter{})

		RegisterLazyCreator[Counter](registrationType, &simpleCounter{})

		counters, _ := GetList[Counter](context.Background())

		if got := len(counters); got != 3 {
			t.Errorf("got %v, expected %v", got, 3)
		}
	}
}

func TestRegisterLazyCreatorMultipleImplementationsKeyed(t *testing.T) {
	for _, registrationType := range types {
		clearAll()

		RegisterLazyCreator[Counter](registrationType, &simpleCounter{}, "firas")

		RegisterLazyCreator[Counter](registrationType, &simpleCounter{}, "firas")

		RegisterLazyCreator[Counter](registrationType, &simpleCounter{})

		counters, _ := GetList[Counter](context.Background(), "firas")

		if got := len(counters); got != 2 {
			t.Errorf("got %v, expected %v", got, 2)
		}
	}
}

func TestRegisterLazyCreatorSingletonState(t *testing.T) {
	var registrationType RegistrationType = Singleton

	clearAll()

	RegisterLazyCreator[Counter](registrationType, &simpleCounter{})

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

func TestRegisterLazyCreatorScopedState(t *testing.T) {
	var registrationType RegistrationType = Scoped

	clearAll()

	RegisterLazyCreator[Counter](registrationType, &simpleCounter{})

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

func TestRegisterLazyCreatorTransientState(t *testing.T) {
	var registrationType RegistrationType = Transient

	clearAll()

	RegisterLazyCreator[Counter](registrationType, &simpleCounter{})

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

func TestRegisterLazyCreatorNilKeyOnRegistering(t *testing.T) {
	clearAll()
	defer mustHavePanicked(t)
	RegisterLazyCreator[Counter](Scoped, &simpleCounter{}, nil)
}

func TestRegisterLazyCreatorNilKeyOnGetting(t *testing.T) {
	clearAll()
	defer mustHavePanicked(t)
	RegisterLazyCreator[Counter](Scoped, &simpleCounter{}, "firas")

	Get[Counter](context.Background(), nil)
}

func TestRegisterLazyCreatorGeneric(t *testing.T) {
	for _, registrationType := range types {
		clearAll()

		RegisterLazyCreator[CounterGeneric[uint]](registrationType, &counterGeneric[uint]{})

		c, _ := Get[CounterGeneric[uint]](context.Background())

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

		RegisterLazyCreator[CounterGeneric[uint]](registrationType, &counterGeneric[uint]{})

		RegisterLazyCreator[CounterGeneric[uint]](registrationType, &counterGeneric[uint]{})

		RegisterLazyCreator[CounterGeneric[uint]](registrationType, &counterGeneric[uint]{})

		counters, _ := GetList[CounterGeneric[uint]](context.Background())

		if got := len(counters); got != 3 {
			t.Errorf("got %v, expected %v", got, 3)
		}
	}
}
