package ore

import (
	"context"
	"github.com/firasdarwish/ore/internal/interfaces"
	"github.com/firasdarwish/ore/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegisterEagerSingleton(t *testing.T) {
	clearAll()

	RegisterEagerSingleton[interfaces.SomeCounter](&models.SimpleCounter{})

	c, _ := Get[interfaces.SomeCounter](context.Background())

	c.AddOne()
	c.AddOne()

	if got := c.GetCount(); got != 2 {
		t.Errorf("got %v, expected %v", got, 2)
	}
}

func TestRegisterEagerSingletonNilImplementation(t *testing.T) {
	clearAll()
	assert.Panics(t, func() {
		RegisterEagerSingleton[interfaces.SomeCounter](nil)
	})
}

func TestRegisterEagerSingletonMultipleImplementations(t *testing.T) {
	clearAll()

	RegisterEagerSingleton[interfaces.SomeCounter](&models.SimpleCounter{})
	RegisterEagerSingleton[interfaces.SomeCounter](&models.SimpleCounter{})
	RegisterEagerSingleton[interfaces.SomeCounter](&models.SimpleCounter{})

	counters, _ := GetList[interfaces.SomeCounter](context.Background())

	if got := len(counters); got != 3 {
		t.Errorf("got %v, expected %v", got, 3)
	}
}

func TestRegisterEagerSingletonMultipleImplementationsKeyed(t *testing.T) {
	clearAll()

	RegisterEagerSingleton[interfaces.SomeCounter](&models.SimpleCounter{}, "firas")
	RegisterEagerSingleton[interfaces.SomeCounter](&models.SimpleCounter{}, "firas")

	RegisterEagerSingleton[interfaces.SomeCounter](&models.SimpleCounter{})

	counters, _ := GetList[interfaces.SomeCounter](context.Background(), "firas")

	if got := len(counters); got != 2 {
		t.Errorf("got %v, expected %v", got, 2)
	}
}

func TestRegisterEagerSingletonSingletonState(t *testing.T) {
	clearAll()

	RegisterEagerSingleton[interfaces.SomeCounter](&models.SimpleCounter{})

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

func TestRegisterEagerSingletonNilKeyOnRegistering(t *testing.T) {
	clearAll()
	assert.Panics(t, func() {
		RegisterEagerSingleton[interfaces.SomeCounter](&models.SimpleCounter{}, nil, "")
	})
}

func TestRegisterEagerSingletonNilKeyOnGetting(t *testing.T) {
	clearAll()
	RegisterEagerSingleton[interfaces.SomeCounter](&models.SimpleCounter{}, "firas")
	assert.Panics(t, func() {
		Get[interfaces.SomeCounter](context.Background(), nil, "")
	})
}

func TestRegisterEagerSingletonGeneric(t *testing.T) {
	clearAll()

	RegisterEagerSingleton[interfaces.SomeCounterGeneric[uint]](&models.CounterGeneric[uint]{})

	c, _ := Get[interfaces.SomeCounterGeneric[uint]](context.Background())

	c.Add(5)
	c.Add(5)

	if got := c.GetCount(); got != 10 {
		t.Errorf("got %v, expected %v", got, 10)
	}
}

func TestRegisterEagerSingletonMultipleGenericImplementations(t *testing.T) {
	clearAll()

	RegisterEagerSingleton[interfaces.SomeCounterGeneric[uint]](&models.CounterGeneric[uint]{})
	RegisterEagerSingleton[interfaces.SomeCounterGeneric[uint]](&models.CounterGeneric[uint]{})
	RegisterEagerSingleton[interfaces.SomeCounterGeneric[uint]](&models.CounterGeneric[uint]{})

	counters, _ := GetList[interfaces.SomeCounterGeneric[uint]](context.Background())

	if got := len(counters); got != 3 {
		t.Errorf("got %v, expected %v", got, 3)
	}
}
