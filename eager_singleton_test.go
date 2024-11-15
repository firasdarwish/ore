package ore

import (
	"context"
	"github.com/firasdarwish/ore/internal/interfaces"
	"github.com/firasdarwish/ore/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegisterSingleton(t *testing.T) {
	clearAll()

	RegisterSingleton[interfaces.SomeCounter](&models.SimpleCounter{})
	c, _ := Get[interfaces.SomeCounter](context.Background())
	c.AddOne()
	c.AddOne()
	c.AddOne()

	RegisterSingletonToContainer[interfaces.SomeCounter](DefaultContainer, &models.SimpleCounter{})
	c, _ = Get[interfaces.SomeCounter](context.Background())

	c.AddOne()
	c.AddOne()

	if got := c.GetCount(); got != 2 {
		t.Errorf("got %v, expected %v", got, 2)
	}
}

func TestRegisterSingletonNilImplementation(t *testing.T) {
	clearAll()
	assert.Panics(t, func() {
		RegisterSingleton[interfaces.SomeCounter](nil)
	})
}

func TestRegisterSingletonMultipleImplementations(t *testing.T) {
	clearAll()

	RegisterSingleton[interfaces.SomeCounter](&models.SimpleCounter{})
	RegisterSingleton[interfaces.SomeCounter](&models.SimpleCounter{})
	RegisterSingleton[interfaces.SomeCounter](&models.SimpleCounter{})

	counters, _ := GetList[interfaces.SomeCounter](context.Background())

	if got := len(counters); got != 3 {
		t.Errorf("got %v, expected %v", got, 3)
	}
}

func TestRegisterSingletonMultipleImplementationsKeyed(t *testing.T) {
	clearAll()

	RegisterKeyedSingleton[interfaces.SomeCounter](&models.SimpleCounter{}, "firas")
	RegisterKeyedSingletonToContainer[interfaces.SomeCounter](DefaultContainer, &models.SimpleCounter{}, "firas")
	RegisterKeyedSingleton[interfaces.SomeCounter](&models.SimpleCounter{}, "firas")

	RegisterSingleton[interfaces.SomeCounter](&models.SimpleCounter{})

	counters, _ := GetKeyedList[interfaces.SomeCounter](context.Background(), "firas")

	if got := len(counters); got != 3 {
		t.Errorf("got %v, expected %v", got, 3)
	}
}

func TestRegisterSingletonSingletonState(t *testing.T) {
	clearAll()

	RegisterSingleton[interfaces.SomeCounter](&models.SimpleCounter{})

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

func TestRegisterSingletonNilKeyOnRegistering(t *testing.T) {
	clearAll()
	assert.Panics(t, func() {
		RegisterKeyedSingleton[interfaces.SomeCounter](&models.SimpleCounter{}, nil)
	})
}

func TestRegisterSingletonNilKeyOnGetting(t *testing.T) {
	clearAll()
	RegisterKeyedSingleton[interfaces.SomeCounter](&models.SimpleCounter{}, "firas")
	assert.Panics(t, func() {
		GetKeyed[interfaces.SomeCounter](context.Background(), nil)
	})
}

func TestRegisterSingletonGeneric(t *testing.T) {
	clearAll()

	RegisterSingleton[interfaces.SomeCounterGeneric[uint]](&models.CounterGeneric[uint]{})

	c, _ := Get[interfaces.SomeCounterGeneric[uint]](context.Background())

	c.Add(5)
	c.Add(5)

	if got := c.GetCount(); got != 10 {
		t.Errorf("got %v, expected %v", got, 10)
	}
}

func TestRegisterSingletonMultipleGenericImplementations(t *testing.T) {
	clearAll()

	RegisterSingleton[interfaces.SomeCounterGeneric[uint]](&models.CounterGeneric[uint]{})
	RegisterSingleton[interfaces.SomeCounterGeneric[uint]](&models.CounterGeneric[uint]{})
	RegisterSingleton[interfaces.SomeCounterGeneric[uint]](&models.CounterGeneric[uint]{})

	counters, _ := GetList[interfaces.SomeCounterGeneric[uint]](context.Background())

	if got := len(counters); got != 3 {
		t.Errorf("got %v, expected %v", got, 3)
	}
}
