package ore

import (
	"context"
	"github.com/firasdarwish/ore/internal/interfaces"
	"github.com/firasdarwish/ore/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegisterCreator(t *testing.T) {
	for _, registrationType := range types {
		clearAll()

		RegisterCreator[interfaces.SomeCounter](registrationType, &models.SimpleCounter{})

		c, _ := Get[interfaces.SomeCounter](context.Background())

		c.AddOne()
		c.AddOne()

		if got := c.GetCount(); got != 2 {
			t.Errorf("got %v, expected %v", got, 2)
		}
	}
}

func TestRegisterCreatorNilFuncTransient(t *testing.T) {
	clearAll()
	assert.Panics(t, func() {
		RegisterCreator[interfaces.SomeCounter](Transient, nil)
	})
}

func TestRegisterCreatorNilFuncScoped(t *testing.T) {
	clearAll()
	assert.Panics(t, func() {
		RegisterCreator[interfaces.SomeCounter](Scoped, nil)
	})
}

func TestRegisterCreatorNilFuncSingleton(t *testing.T) {
	clearAll()
	assert.Panics(t, func() {
		RegisterCreator[interfaces.SomeCounter](Singleton, nil)
	})
}

func TestRegisterCreatorMultipleImplementations(t *testing.T) {
	for _, registrationType := range types {
		clearAll()

		RegisterCreator[interfaces.SomeCounter](registrationType, &models.SimpleCounter{})

		RegisterCreator[interfaces.SomeCounter](registrationType, &models.SimpleCounter{})

		RegisterCreator[interfaces.SomeCounter](registrationType, &models.SimpleCounter{})

		counters, _ := GetList[interfaces.SomeCounter](context.Background())

		if got := len(counters); got != 3 {
			t.Errorf("got %v, expected %v", got, 3)
		}
	}
}

func TestRegisterCreatorMultipleImplementationsKeyed(t *testing.T) {
	for _, registrationType := range types {
		clearAll()

		RegisterCreator[interfaces.SomeCounter](registrationType, &models.SimpleCounter{}, "firas")

		RegisterCreator[interfaces.SomeCounter](registrationType, &models.SimpleCounter{}, "firas")

		RegisterCreator[interfaces.SomeCounter](registrationType, &models.SimpleCounter{})

		counters, _ := GetList[interfaces.SomeCounter](context.Background(), "firas")

		if got := len(counters); got != 2 {
			t.Errorf("got %v, expected %v", got, 2)
		}
	}
}

func TestRegisterCreatorSingletonState(t *testing.T) {
	var registrationType Lifetime = Singleton

	clearAll()

	RegisterCreator[interfaces.SomeCounter](registrationType, &models.SimpleCounter{})

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

func TestRegisterCreatorScopedState(t *testing.T) {
	var registrationType Lifetime = Scoped

	clearAll()

	RegisterCreator[interfaces.SomeCounter](registrationType, &models.SimpleCounter{})

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

func TestRegisterCreatorTransientState(t *testing.T) {
	var registrationType Lifetime = Transient

	clearAll()

	RegisterCreator[interfaces.SomeCounter](registrationType, &models.SimpleCounter{})

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

func TestRegisterCreatorNilKeyOnRegistering(t *testing.T) {
	clearAll()
	assert.Panics(t, func() {
		RegisterCreator[interfaces.SomeCounter](Scoped, &models.SimpleCounter{}, "", nil)
	})
}

func TestRegisterCreatorNilKeyOnGetting(t *testing.T) {
	clearAll()
	RegisterCreator[interfaces.SomeCounter](Scoped, &models.SimpleCounter{}, "firas")
	assert.Panics(t, func() {
		Get[interfaces.SomeCounter](context.Background(), nil)
	})
}

func TestRegisterCreatorGeneric(t *testing.T) {
	for _, registrationType := range types {
		clearAll()

		RegisterCreator[interfaces.SomeCounterGeneric[uint]](registrationType, &models.CounterGeneric[uint]{})

		c, _ := Get[interfaces.SomeCounterGeneric[uint]](context.Background())

		c.Add(5)
		c.Add(5)

		if got := c.GetCount(); got != 10 {
			t.Errorf("got %v, expected %v", got, 10)
		}
	}
}

func TestRegisterCreatorMultipleGenericImplementations(t *testing.T) {
	for _, registrationType := range types {
		clearAll()

		RegisterCreator[interfaces.SomeCounterGeneric[uint]](registrationType, &models.CounterGeneric[uint]{})

		RegisterCreator[interfaces.SomeCounterGeneric[uint]](registrationType, &models.CounterGeneric[uint]{})

		RegisterCreator[interfaces.SomeCounterGeneric[uint]](registrationType, &models.CounterGeneric[uint]{})

		counters, _ := GetList[interfaces.SomeCounterGeneric[uint]](context.Background())

		if got := len(counters); got != 3 {
			t.Errorf("got %v, expected %v", got, 3)
		}
	}
}
