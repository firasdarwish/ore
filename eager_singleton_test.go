package ore

import (
	"context"
	"testing"
)

func TestRegisterEagerSingleton(t *testing.T) {
	clearAll()

	RegisterEagerSingleton[Counter](&simpleCounter{})

	c, _ := Get[Counter](context.Background())

	c.AddOne()
	c.AddOne()

	if got := c.GetCount(); got != 2 {
		t.Errorf("got %v, expected %v", got, 2)
	}
}

func TestRegisterEagerSingletonNilImplementation(t *testing.T) {
	clearAll()
	defer mustHavePanicked(t)
	RegisterEagerSingleton[Counter](nil)
}

func TestRegisterEagerSingletonMultipleImplementations(t *testing.T) {
	clearAll()

	RegisterEagerSingleton[Counter](&simpleCounter{})
	RegisterEagerSingleton[Counter](&simpleCounter{})
	RegisterEagerSingleton[Counter](&simpleCounter{})

	counters, _ := GetList[Counter](context.Background())

	if got := len(counters); got != 3 {
		t.Errorf("got %v, expected %v", got, 3)
	}
}

func TestRegisterEagerSingletonMultipleImplementationsKeyed(t *testing.T) {
	clearAll()

	RegisterEagerSingleton[Counter](&simpleCounter{}, "firas")
	RegisterEagerSingleton[Counter](&simpleCounter{}, "firas")

	RegisterEagerSingleton[Counter](&simpleCounter{})

	counters, _ := GetList[Counter](context.Background(), "firas")

	if got := len(counters); got != 2 {
		t.Errorf("got %v, expected %v", got, 2)
	}
}

func TestRegisterEagerSingletonSingletonState(t *testing.T) {
	clearAll()

	RegisterEagerSingleton[Counter](&simpleCounter{})

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

func TestRegisterEagerSingletonNilKeyOnRegistering(t *testing.T) {
	clearAll()
	defer mustHavePanicked(t)
	RegisterEagerSingleton[Counter](&simpleCounter{}, nil)
}

func TestRegisterEagerSingletonNilKeyOnGetting(t *testing.T) {
	clearAll()
	defer mustHavePanicked(t)
	RegisterEagerSingleton[Counter](&simpleCounter{}, "firas")

	Get[Counter](context.Background(), nil)
}

func TestRegisterEagerSingletonGeneric(t *testing.T) {
	clearAll()

	RegisterEagerSingleton[CounterGeneric[uint]](&counterGeneric[uint]{})

	c, _ := Get[CounterGeneric[uint]](context.Background())

	c.Add(5)
	c.Add(5)

	if got := c.GetCount(); got != 10 {
		t.Errorf("got %v, expected %v", got, 10)
	}
}

func TestRegisterEagerSingletonMultipleGenericImplementations(t *testing.T) {
	clearAll()

	RegisterEagerSingleton[CounterGeneric[uint]](&counterGeneric[uint]{})
	RegisterEagerSingleton[CounterGeneric[uint]](&counterGeneric[uint]{})
	RegisterEagerSingleton[CounterGeneric[uint]](&counterGeneric[uint]{})

	counters, _ := GetList[CounterGeneric[uint]](context.Background())

	if got := len(counters); got != 3 {
		t.Errorf("got %v, expected %v", got, 3)
	}
}
