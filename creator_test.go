package ore

import (
	"context"
	"github.com/firasdarwish/ore/internal/interfaces"
	"github.com/firasdarwish/ore/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegisterLazyCreator(t *testing.T) {
	for _, registrationType := range types {
		clearAll()

		RegisterLazyCreator[interfaces.SomeCounter](registrationType, &models.SimpleCounter{})

		c, _ := Get[interfaces.SomeCounter](context.Background())

		c.AddOne()
		c.AddOne()

		if got := c.GetCount(); got != 2 {
			t.Errorf("got %v, expected %v", got, 2)
		}
	}
}

func TestRegisterLazyCreatorNilFuncTransient(t *testing.T) {
	clearAll()
	assert.Panics(t, func() {
		RegisterLazyCreator[interfaces.SomeCounter](Transient, nil)
	})
}

func TestRegisterLazyCreatorNilFuncScoped(t *testing.T) {
	clearAll()
	assert.Panics(t, func() {
		RegisterLazyCreator[interfaces.SomeCounter](Scoped, nil)
	})
}

func TestRegisterLazyCreatorNilFuncSingleton(t *testing.T) {
	clearAll()
	assert.Panics(t, func() {
		RegisterLazyCreator[interfaces.SomeCounter](Singleton, nil)
	})
}

func TestRegisterLazyCreatorMultipleImplementations(t *testing.T) {
	for _, registrationType := range types {
		clearAll()

		RegisterLazyCreator[interfaces.SomeCounter](registrationType, &models.SimpleCounter{})

		RegisterLazyCreator[interfaces.SomeCounter](registrationType, &models.SimpleCounter{})

		RegisterLazyCreator[interfaces.SomeCounter](registrationType, &models.SimpleCounter{})

		counters, _ := GetList[interfaces.SomeCounter](context.Background())

		if got := len(counters); got != 3 {
			t.Errorf("got %v, expected %v", got, 3)
		}
	}
}

func TestRegisterLazyCreatorMultipleImplementationsKeyed(t *testing.T) {
	for _, registrationType := range types {
		clearAll()

		RegisterLazyCreator[interfaces.SomeCounter](registrationType, &models.SimpleCounter{}, "firas")

		RegisterLazyCreator[interfaces.SomeCounter](registrationType, &models.SimpleCounter{}, "firas")

		RegisterLazyCreator[interfaces.SomeCounter](registrationType, &models.SimpleCounter{})

		counters, _ := GetList[interfaces.SomeCounter](context.Background(), "firas")

		if got := len(counters); got != 2 {
			t.Errorf("got %v, expected %v", got, 2)
		}
	}
}

func TestRegisterLazyCreatorSingletonState(t *testing.T) {
	var registrationType Lifetime = Singleton

	clearAll()

	RegisterLazyCreator[interfaces.SomeCounter](registrationType, &models.SimpleCounter{})

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

func TestRegisterLazyCreatorScopedState(t *testing.T) {
	var registrationType Lifetime = Scoped

	clearAll()

	RegisterLazyCreator[interfaces.SomeCounter](registrationType, &models.SimpleCounter{})

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

func TestRegisterLazyCreatorTransientState(t *testing.T) {
	var registrationType Lifetime = Transient

	clearAll()

	RegisterLazyCreator[interfaces.SomeCounter](registrationType, &models.SimpleCounter{})

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

func TestRegisterLazyCreatorNilKeyOnRegistering(t *testing.T) {
	clearAll()
	assert.Panics(t, func() {
		RegisterLazyCreator[interfaces.SomeCounter](Scoped, &models.SimpleCounter{}, "", nil)
	})
}

func TestRegisterLazyCreatorNilKeyOnGetting(t *testing.T) {
	clearAll()
	RegisterLazyCreator[interfaces.SomeCounter](Scoped, &models.SimpleCounter{}, "firas")
	assert.Panics(t, func() {
		Get[interfaces.SomeCounter](context.Background(), nil)
	})
}

func TestRegisterLazyCreatorGeneric(t *testing.T) {
	for _, registrationType := range types {
		clearAll()

		RegisterLazyCreator[interfaces.SomeCounterGeneric[uint]](registrationType, &models.CounterGeneric[uint]{})

		c, _ := Get[interfaces.SomeCounterGeneric[uint]](context.Background())

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

		RegisterLazyCreator[interfaces.SomeCounterGeneric[uint]](registrationType, &models.CounterGeneric[uint]{})

		RegisterLazyCreator[interfaces.SomeCounterGeneric[uint]](registrationType, &models.CounterGeneric[uint]{})

		RegisterLazyCreator[interfaces.SomeCounterGeneric[uint]](registrationType, &models.CounterGeneric[uint]{})

		counters, _ := GetList[interfaces.SomeCounterGeneric[uint]](context.Background())

		if got := len(counters); got != 3 {
			t.Errorf("got %v, expected %v", got, 3)
		}
	}
}
