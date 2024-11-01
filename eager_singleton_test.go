package ore

import (
	"context"
	"testing"
)

func TestRegisterEagerSingleton(t *testing.T) {
	clearAll()

	RegisterEagerSingleton[someCounter](&simpleCounter{})

	c, _ := Get[someCounter](context.Background())

	c.AddOne()
	c.AddOne()

	if got := c.GetCount(); got != 2 {
		t.Errorf("got %v, expected %v", got, 2)
	}
}

func TestRegisterEagerSingletonNilImplementation(t *testing.T) {
	clearAll()
	defer mustHavePanicked(t)
	RegisterEagerSingleton[someCounter](nil)
}

func TestRegisterEagerSingletonMultipleImplementations(t *testing.T) {
	clearAll()

	RegisterEagerSingleton[someCounter](&simpleCounter{})
	RegisterEagerSingleton[someCounter](&simpleCounter{})
	RegisterEagerSingleton[someCounter](&simpleCounter{})

	counters, _ := GetList[someCounter](context.Background())

	if got := len(counters); got != 3 {
		t.Errorf("got %v, expected %v", got, 3)
	}
}

func TestRegisterEagerSingletonMultipleImplementationsKeyed(t *testing.T) {
	clearAll()

	RegisterEagerSingleton[someCounter](&simpleCounter{}, "firas")
	RegisterEagerSingleton[someCounter](&simpleCounter{}, "firas")

	RegisterEagerSingleton[someCounter](&simpleCounter{})

	counters, _ := GetList[someCounter](context.Background(), "firas")

	if got := len(counters); got != 2 {
		t.Errorf("got %v, expected %v", got, 2)
	}
}

func TestRegisterEagerSingletonSingletonState(t *testing.T) {
	clearAll()

	RegisterEagerSingleton[someCounter](&simpleCounter{})

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

func TestRegisterEagerSingletonNilKeyOnRegistering(t *testing.T) {
	clearAll()
	defer mustHavePanicked(t)
	RegisterEagerSingleton[someCounter](&simpleCounter{}, nil)
}

func TestRegisterEagerSingletonNilKeyOnGetting(t *testing.T) {
	clearAll()
	defer mustHavePanicked(t)
	RegisterEagerSingleton[someCounter](&simpleCounter{}, "firas")

	Get[someCounter](context.Background(), nil)
}

func TestRegisterEagerSingletonGeneric(t *testing.T) {
	clearAll()

	RegisterEagerSingleton[someCounterGeneric[uint]](&counterGeneric[uint]{})

	c, _ := Get[someCounterGeneric[uint]](context.Background())

	c.Add(5)
	c.Add(5)

	if got := c.GetCount(); got != 10 {
		t.Errorf("got %v, expected %v", got, 10)
	}
}

func TestRegisterEagerSingletonMultipleGenericImplementations(t *testing.T) {
	clearAll()

	RegisterEagerSingleton[someCounterGeneric[uint]](&counterGeneric[uint]{})
	RegisterEagerSingleton[someCounterGeneric[uint]](&counterGeneric[uint]{})
	RegisterEagerSingleton[someCounterGeneric[uint]](&counterGeneric[uint]{})

	counters, _ := GetList[someCounterGeneric[uint]](context.Background())

	if got := len(counters); got != 3 {
		t.Errorf("got %v, expected %v", got, 3)
	}
}
