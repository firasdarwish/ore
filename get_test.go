package ore

import (
	"context"
	"fmt"
	"github.com/firasdarwish/ore/internal/interfaces"
	"testing"
	"time"

	m "github.com/firasdarwish/ore/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	for _, registrationType := range types {
		clearAll()

		RegisterCreator[interfaces.SomeCounter](registrationType, &m.SimpleCounter{})

		c, _ := Get[interfaces.SomeCounter](context.Background())

		c.AddOne()
		c.AddOne()

		if got := c.GetCount(); got != 2 {
			t.Errorf("got %v, expected %v", got, 2)
		}
	}
}

func TestGetLatestByDefault(t *testing.T) {
	for _, registrationType := range types {
		clearAll()

		RegisterCreator[interfaces.SomeCounter](registrationType, &m.SimpleCounter{})
		c, _ := Get[interfaces.SomeCounter](context.Background())
		c.AddOne()
		c.AddOne()

		RegisterCreator[interfaces.SomeCounter](registrationType, &m.SimpleCounter2{})
		c, _ = Get[interfaces.SomeCounter](context.Background())
		c.AddOne()
		c.AddOne()
		c.AddOne()
		c.AddOne()

		if got := c.GetCount(); got != 4 {
			t.Errorf("got %v, expected %v", got, 4)
		}
	}
}

func TestGetPanicIfNoImplementations(t *testing.T) {
	clearAll()
	assert.Panics(t, func() {
		Get[interfaces.SomeCounter](context.Background())
	})
}

func TestGetKeyedSingleton(t *testing.T) {
	clearAll()

	//ctx := context.Background()
	key := fmt.Sprintf("keynum: %v", Singleton)

	RegisterKeyedCreator[interfaces.SomeCounter](Singleton, &m.SimpleCounter{}, key)

	c, _ := GetKeyed[interfaces.SomeCounter](context.Background(), key)
	c.AddOne()
	c.AddOne()

	c, _ = GetKeyedFromContainer[interfaces.SomeCounter](DefaultContainer, context.Background(), key)
	c.AddOne()
	c.AddOne()

	if got := c.GetCount(); got != 4 {
		t.Errorf("got %v, expected %v", got, 4)
	}
}

func TestGetKeyedScoped(t *testing.T) {
	clearAll()

	ctx := context.Background()
	key := fmt.Sprintf("keynum: %v", Scoped)

	RegisterKeyedCreator[interfaces.SomeCounter](Scoped, &m.SimpleCounter{}, key)

	c, ctx := GetKeyed[interfaces.SomeCounter](ctx, key)
	c.AddOne()
	c.AddOne()

	c, ctx = GetKeyedFromContainer[interfaces.SomeCounter](DefaultContainer, ctx, key)
	c.AddOne()
	c.AddOne()

	if got := c.GetCount(); got != 4 {
		t.Errorf("got %v, expected %v", got, 4)
	}
}

func TestGetKeyedTransient(t *testing.T) {
	clearAll()

	ctx := context.Background()
	key := fmt.Sprintf("keynum: %v", Transient)

	RegisterKeyedCreator[interfaces.SomeCounter](Transient, &m.SimpleCounter{}, key)

	c, ctx := GetKeyed[interfaces.SomeCounter](ctx, key)
	c.AddOne()
	c.AddOne()

	c, ctx = GetKeyedFromContainer[interfaces.SomeCounter](DefaultContainer, ctx, key)
	c.AddOne()
	c.AddOne()

	if got := c.GetCount(); got != 2 {
		t.Errorf("got %v, expected %v", got, 2)
	}
}

func TestGetResolvedSingletons(t *testing.T) {
	t.Run("When multiple lifetimes and keys are registered", func(t *testing.T) {
		//Arrange
		clearAll()
		RegisterFunc(Singleton, func(ctx context.Context) (*m.DisposableService1, context.Context) {
			return &m.DisposableService1{Name: "A1"}, ctx
		})
		RegisterFunc(Singleton, func(ctx context.Context) (*m.DisposableService1, context.Context) {
			return &m.DisposableService1{Name: "A2"}, ctx
		})
		RegisterSingleton(&m.DisposableService2{Name: "E1"})
		RegisterSingleton(&m.DisposableService2{Name: "E2"})
		RegisterFunc(Scoped, func(ctx context.Context) (*m.DisposableService3, context.Context) {
			return &m.DisposableService3{Name: "S1"}, ctx
		})
		RegisterFunc(Transient, func(ctx context.Context) (*m.DisposableService3, context.Context) {
			return &m.DisposableService3{Name: "S2"}, ctx
		})
		RegisterFunc(Singleton, func(ctx context.Context) (*m.DisposableService4, context.Context) {
			return &m.DisposableService4{Name: "X1"}, ctx
		})
		RegisterKeyedFunc(Singleton, func(ctx context.Context) (*m.DisposableService4, context.Context) {
			return &m.DisposableService4{Name: "X2"}, ctx
		}, "somekey")

		ctx := context.Background()
		//Act
		disposables := GetResolvedSingletons[m.Disposer]() //E1, E2
		assert.Equal(t, 2, len(disposables))

		//invoke A1, A2
		_, ctx = GetList[*m.DisposableService1](ctx) //A1, A2

		//Act
		disposables = GetResolvedSingletons[m.Disposer]() //E1, E2, A1, A2
		assert.Equal(t, 4, len(disposables))

		//invoke S1, S2, X1
		RegisterAlias[fmt.Stringer, *m.DisposableService3]()
		RegisterAlias[fmt.Stringer, *m.DisposableService4]()
		_, ctx = GetList[fmt.Stringer](ctx) //S1, S2, X1

		//Act
		//because S1, S2 are not singleton, so they won't be returned, only X1 will be returned in addition
		disposables = GetResolvedSingletons[m.Disposer]() //E1, E2, A1, A2, X1
		assert.Equal(t, 5, len(disposables))

		//invoke X2 in "somekey" scope
		_, _ = GetKeyedList[fmt.Stringer](ctx, "somekey")

		//Act
		//all invoked singleton would be returned whatever keys they are registered with
		disposables = GetResolvedSingletons[m.Disposer]() //E1, E2, A1, A2, X1, X2
		assert.Equal(t, 6, len(disposables))
	})

	t.Run("respect invocation chronological time order", func(t *testing.T) {
		//Arrange
		clearAll()
		RegisterFunc(Singleton, func(ctx context.Context) (*m.DisposableService1, context.Context) {
			return &m.DisposableService1{Name: "A"}, ctx
		})
		RegisterFunc(Singleton, func(ctx context.Context) (*m.DisposableService2, context.Context) {
			return &m.DisposableService2{Name: "B"}, ctx
		})
		RegisterFunc(Singleton, func(ctx context.Context) (*m.DisposableService3, context.Context) {
			return &m.DisposableService3{Name: "C"}, ctx
		})

		ctx := context.Background()

		//invocation order: [A,C,B]
		_, ctx = Get[*m.DisposableService1](ctx)
		time.Sleep(1 * time.Microsecond)
		_, ctx = Get[*m.DisposableService3](ctx)
		time.Sleep(1 * time.Microsecond)
		_, _ = Get[*m.DisposableService2](ctx)

		//Act
		disposables := GetResolvedSingletons[m.Disposer]() //B, A

		//Assert that the order is [B,C,A], the most recent invocation would be returned first
		assert.Equal(t, 3, len(disposables))
		assert.Equal(t, "B", disposables[0].String())
		assert.Equal(t, "C", disposables[1].String())
		assert.Equal(t, "A", disposables[2].String())
	})
	t.Run("deeper invocation level is returned first", func(t *testing.T) {
		//Arrange
		clearAll()
		RegisterFunc(Singleton, func(ctx context.Context) (*m.DisposableService1, context.Context) {
			_, ctx = Get[*m.DisposableService2](ctx) //1 calls 2
			return &m.DisposableService1{Name: "1"}, ctx
		})
		RegisterFunc(Singleton, func(ctx context.Context) (*m.DisposableService2, context.Context) {
			_, ctx = Get[*m.DisposableService3](ctx) //2 calls 3
			return &m.DisposableService2{Name: "2"}, ctx
		})
		RegisterFunc(Singleton, func(ctx context.Context) (*m.DisposableService3, context.Context) {
			return &m.DisposableService3{Name: "3"}, ctx
		})

		//invocation order: [1,2,3]
		_, _ = Get[*m.DisposableService1](context.Background())

		//Act
		disposables := GetResolvedSingletons[m.Disposer]()

		//Assert that the order is [B,C,A], the deepest invocation level would be returned first
		assert.Equal(t, 3, len(disposables))
		assert.Equal(t, "3", disposables[0].String())
		assert.Equal(t, "2", disposables[1].String())
		assert.Equal(t, "1", disposables[2].String())
	})
}

func TestGetResolvedScopedInstances(t *testing.T) {
	t.Run("When multiple lifetimes and keys are registered", func(t *testing.T) {
		clearAll()
		RegisterSingleton(&m.DisposableService1{Name: "S1"})
		RegisterFunc(Scoped, func(ctx context.Context) (*m.DisposableService1, context.Context) {
			return &m.DisposableService1{Name: "S2"}, ctx
		})
		RegisterKeyedFunc(Scoped, func(ctx context.Context) (*m.DisposableService2, context.Context) {
			return &m.DisposableService2{Name: "T1"}, ctx
		}, "module1")

		ctx := context.Background()

		//Act
		disposables := GetResolvedScopedInstances[m.Disposer](ctx) //empty
		assert.Empty(t, disposables)

		//invoke S2
		_, ctx = GetList[*m.DisposableService1](ctx)

		//Act
		disposables = GetResolvedScopedInstances[m.Disposer](ctx) //S2
		assert.Equal(t, 1, len(disposables))
		assert.Equal(t, "S2", disposables[0].String())

		//invoke the keyed service T1
		_, ctx = GetKeyedList[*m.DisposableService2](ctx, "module1")

		//Act
		disposables = GetResolvedScopedInstances[m.Disposer](ctx) //S2, T1
		assert.Equal(t, 2, len(disposables))
	})

	t.Run("respect invocation chronological time order", func(t *testing.T) {
		//Arrange
		clearAll()
		RegisterFunc(Scoped, func(ctx context.Context) (*m.DisposableService1, context.Context) {
			return &m.DisposableService1{Name: "A"}, ctx
		})
		RegisterFunc(Scoped, func(ctx context.Context) (*m.DisposableService2, context.Context) {
			return &m.DisposableService2{Name: "B"}, ctx
		})
		RegisterFunc(Scoped, func(ctx context.Context) (*m.DisposableService3, context.Context) {
			return &m.DisposableService3{Name: "C"}, ctx
		})

		ctx := context.Background()

		//invocation order: [A,C,B]
		_, ctx = Get[*m.DisposableService1](ctx)
		time.Sleep(1 * time.Microsecond)
		_, ctx = Get[*m.DisposableService3](ctx)
		time.Sleep(1 * time.Microsecond)
		_, ctx = Get[*m.DisposableService2](ctx)

		//Act
		disposables := GetResolvedScopedInstances[m.Disposer](ctx) //B, A

		//Assert that the order is [B,C,A], the most recent invocation would be returned first
		assert.Equal(t, 3, len(disposables))
		assert.Equal(t, "B", disposables[0].String())
		assert.Equal(t, "C", disposables[1].String())
		assert.Equal(t, "A", disposables[2].String())
	})

	t.Run("respect invocation deep level", func(t *testing.T) {
		//Arrange
		clearAll()
		RegisterFunc(Scoped, func(ctx context.Context) (*m.DisposableService1, context.Context) {
			//1 calls 3
			_, ctx = Get[*m.DisposableService3](ctx)
			return &m.DisposableService1{Name: "1"}, ctx
		})
		RegisterFunc(Scoped, func(ctx context.Context) (*m.DisposableService2, context.Context) {
			return &m.DisposableService2{Name: "2"}, ctx
		})
		RegisterFunc(Scoped, func(ctx context.Context) (*m.DisposableService3, context.Context) {
			return &m.DisposableService3{Name: "3"}, ctx
		})
		RegisterFunc(Scoped, func(ctx context.Context) (*m.DisposableService4, context.Context) {
			//4 calls 1, 2
			_, ctx = Get[*m.DisposableService1](ctx)
			_, ctx = Get[*m.DisposableService2](ctx)
			return &m.DisposableService4{Name: "4"}, ctx
		})

		ctx := context.Background()

		//invocation order: [4,1,3,2].
		_, ctx = Get[*m.DisposableService4](ctx)

		//Act
		disposables := GetResolvedScopedInstances[m.Disposer](ctx)

		assert.Equal(t, 4, len(disposables))

		//find the position of the disposables
		index1 := m.FindIndexOf(disposables, "1")
		index2 := m.FindIndexOf(disposables, "2")
		index3 := m.FindIndexOf(disposables, "3")
		index4 := m.FindIndexOf(disposables, "4")

		//Assert that 4 should be disposed after 1 and 2  (because 4 calls 1 and 2)
		assert.Greater(t, index4, index1)
		assert.Greater(t, index4, index2)

		//Assert that 1 should be disposed after 3 (because 1 calls 3)
		assert.Greater(t, index1, index3)
	})
}
