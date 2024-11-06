package ore

import (
	"context"
	"testing"

	m "github.com/firasdarwish/ore/internal/models"
	"github.com/firasdarwish/ore/internal/testtools/assert2"
	"github.com/stretchr/testify/assert"
)

func TestValidate_CircularDepsUniformLifetype(t *testing.T) {
	for _, lt := range types {
		t.Run("Direct circular "+lt.String()+" (1 calls 1)", func(t *testing.T) {
			clearAll()
			RegisterLazyFunc(lt, func(ctx context.Context) (*m.DisposableService1, context.Context) {
				_, ctx = Get[*m.DisposableService1](ctx) //1 calls 1
				return &m.DisposableService1{Name: "1"}, ctx
			})
			assert2.PanicsWithError(t, assert2.ErrorStartsWith("detect cyclic dependency"), Validate)
		})
		t.Run("Indirect circular "+lt.String()+" (1 calls 2 calls 3 calls 1)", func(t *testing.T) {
			clearAll()
			RegisterLazyFunc(lt, func(ctx context.Context) (*m.DisposableService1, context.Context) {
				_, ctx = Get[*m.DisposableService2](ctx) //1 calls 2
				return &m.DisposableService1{Name: "1"}, ctx
			})
			RegisterLazyFunc(lt, func(ctx context.Context) (*m.DisposableService2, context.Context) {
				_, ctx = Get[*m.DisposableService3](ctx) //2 calls 3
				return &m.DisposableService2{Name: "2"}, ctx
			})
			RegisterLazyFunc(lt, func(ctx context.Context) (*m.DisposableService3, context.Context) {
				_, ctx = Get[*m.DisposableService1](ctx) //3 calls 1
				return &m.DisposableService3{Name: "3"}, ctx
			})
			assert2.PanicsWithError(t, assert2.ErrorStartsWith("detect cyclic dependency"), Validate)
		})
		t.Run("Middle circular "+lt.String()+" (1 calls 2 calls 3 calls 4 calls 2)", func(t *testing.T) {
			clearAll()
			RegisterLazyFunc(lt, func(ctx context.Context) (*m.DisposableService1, context.Context) {
				_, ctx = Get[*m.DisposableService2](ctx) //1 calls 2
				return &m.DisposableService1{Name: "1"}, ctx
			})
			RegisterLazyFunc(lt, func(ctx context.Context) (*m.DisposableService2, context.Context) {
				_, ctx = Get[*m.DisposableService3](ctx) //2 calls 3
				return &m.DisposableService2{Name: "2"}, ctx
			})
			RegisterLazyFunc(lt, func(ctx context.Context) (*m.DisposableService3, context.Context) {
				_, ctx = Get[*m.DisposableService4](ctx) //3 calls 4
				return &m.DisposableService3{Name: "3"}, ctx
			})
			RegisterLazyFunc(lt, func(ctx context.Context) (*m.DisposableService4, context.Context) {
				_, ctx = Get[*m.DisposableService2](ctx) //4 calls 2
				return &m.DisposableService4{Name: "4"}, ctx
			})
			assert2.PanicsWithError(t, assert2.ErrorStartsWith("detect cyclic dependency"), Validate)
		})
		t.Run("circular on complex tree "+lt.String()+"", func(t *testing.T) {
			clearAll()
			RegisterLazyFunc(lt, func(ctx context.Context) (*m.DisposableService1, context.Context) {
				_, ctx = Get[*m.DisposableService2](ctx) //1 calls 2
				_, ctx = Get[*m.DisposableService3](ctx) //1 calls 3
				return &m.DisposableService1{Name: "1"}, ctx
			})
			RegisterLazyFunc(lt, func(ctx context.Context) (*m.DisposableService2, context.Context) {
				_, ctx = Get[*m.DisposableService4](ctx) //2 calls 4
				_, ctx = Get[*m.DisposableService5](ctx) //2 calls 5
				return &m.DisposableService2{Name: "2"}, ctx
			})
			RegisterLazyFunc(lt, func(ctx context.Context) (*m.DisposableService3, context.Context) {
				_, ctx = Get[*m.DisposableService4](ctx) //3 calls 4
				return &m.DisposableService3{Name: "3"}, ctx
			})
			RegisterLazyFunc(lt, func(ctx context.Context) (*m.DisposableService4, context.Context) {
				_, ctx = Get[*m.DisposableService5](ctx) //4 calls 5
				return &m.DisposableService4{Name: "4"}, ctx
			})
			RegisterLazyFunc(lt, func(ctx context.Context) (*m.DisposableService5, context.Context) {
				_, ctx = Get[*m.DisposableService3](ctx) //5 calls 3 => circular here: 5->3->4->5
				return &m.DisposableService5{Name: "5"}, ctx
			})
			assert2.PanicsWithError(t, assert2.ErrorStartsWith("detect cyclic dependency"), Validate)
		})
		t.Run("fake circular top down "+lt.String()+": (1 calls 2 (x2) calls 3 calls 4, 2 calls 4)", func(t *testing.T) {
			clearAll()
			RegisterLazyFunc(lt, func(ctx context.Context) (*m.DisposableService1, context.Context) {
				_, ctx = Get[*m.DisposableService2](ctx) //1 calls 2
				_, ctx = Get[*m.DisposableService2](ctx) //1 calls 2 again
				return &m.DisposableService1{Name: "1"}, ctx
			})
			RegisterLazyFunc(lt, func(ctx context.Context) (*m.DisposableService2, context.Context) {
				_, ctx = Get[*m.DisposableService3](ctx) //2 calls 3
				_, ctx = Get[*m.DisposableService4](ctx) //2 calls 4
				return &m.DisposableService2{Name: "2"}, ctx
			})
			RegisterLazyFunc(lt, func(ctx context.Context) (*m.DisposableService3, context.Context) {
				_, ctx = Get[*m.DisposableService4](ctx) //3 calls 4
				_, ctx = Get[*m.DisposableService4](ctx) //3 calls 4
				return &m.DisposableService3{Name: "3"}, ctx
			})
			RegisterLazyFunc(lt, func(ctx context.Context) (*m.DisposableService4, context.Context) {
				return &m.DisposableService4{Name: "4"}, ctx
			})
			assert.NotPanics(t, Validate)
		})
		t.Run("fake circular sibling "+lt.String()+": 1 calls 2 & 3;  2 calls 3)", func(t *testing.T) {
			clearAll()
			RegisterLazyFunc(lt, func(ctx context.Context) (*m.DisposableService1, context.Context) {
				_, ctx = Get[*m.DisposableService2](ctx) //1 calls 2
				_, ctx = Get[*m.DisposableService3](ctx) //1 calls 3
				return &m.DisposableService1{Name: "1"}, ctx
			})
			RegisterLazyFunc(lt, func(ctx context.Context) (*m.DisposableService2, context.Context) {
				_, ctx = Get[*m.DisposableService3](ctx) //2 calls 3
				return &m.DisposableService2{Name: "2"}, ctx
			})
			RegisterLazyFunc(lt, func(ctx context.Context) (*m.DisposableService3, context.Context) {
				return &m.DisposableService3{Name: "3"}, ctx
			})
			assert.NotPanics(t, Validate)
		})
	}
}

func TestValidate_CircularMixedLifetype(t *testing.T) {
	clearAll()

	RegisterLazyFunc(Scoped, func(ctx context.Context) (*m.DisposableService2, context.Context) {
		_, ctx = Get[*m.DisposableService4](ctx) //2 calls 4
		_, ctx = Get[*m.DisposableService5](ctx) //2 calls 5
		return &m.DisposableService2{Name: "2"}, ctx
	})
	RegisterLazyFunc(Singleton, func(ctx context.Context) (*m.DisposableService3, context.Context) {
		_, ctx = Get[*m.DisposableService4](ctx) //3 calls 4
		return &m.DisposableService3{Name: "3"}, ctx
	})
	RegisterLazyFunc(Singleton, func(ctx context.Context) (*m.DisposableService4, context.Context) {
		_, ctx = Get[*m.DisposableService5](ctx) //4 calls 5
		return &m.DisposableService4{Name: "4"}, ctx
	})
	RegisterLazyFunc(Singleton, func(ctx context.Context) (*m.DisposableService5, context.Context) {
		_, ctx = Get[*m.DisposableService3](ctx) //5 calls 3 => circular here: 5->3->4->5
		return &m.DisposableService5{Name: "5"}, ctx
	})
	RegisterLazyFunc(Transient, func(ctx context.Context) (*m.DisposableService1, context.Context) {
		_, ctx = Get[*m.DisposableService2](ctx) //1 calls 2
		_, ctx = Get[*m.DisposableService3](ctx) //1 calls 3
		return &m.DisposableService1{Name: "1"}, ctx
	})
	assert2.PanicsWithError(t, assert2.ErrorStartsWith("detect cyclic dependency"), Validate)
	assert2.PanicsWithError(t, assert2.ErrorStartsWith("detect cyclic dependency"), func() {
		_, _ = Get[*m.DisposableService1](context.Background())
	})
}

func TestValidate_LifetimeAlignment(t *testing.T) {
	t.Run("Singleton depends on Scoped", func(t *testing.T) {
		clearAll()
		RegisterLazyFunc(Singleton, func(ctx context.Context) (*m.DisposableService1, context.Context) {
			_, ctx = Get[*m.DisposableService2](ctx) //1 depends on 2
			return &m.DisposableService1{Name: "1"}, ctx
		})
		RegisterLazyFunc(Scoped, func(ctx context.Context) (*m.DisposableService2, context.Context) {
			return &m.DisposableService2{Name: "2"}, ctx
		})
		assert2.PanicsWithError(t, assert2.ErrorStartsWith("detect lifetime misalignment"), Validate)
	})
	t.Run("Scoped depends on Transient", func(t *testing.T) {
		clearAll()
		RegisterLazyFunc(Scoped, func(ctx context.Context) (*m.DisposableService1, context.Context) {
			_, ctx = Get[*m.DisposableService2](ctx) //1 depends on 2
			return &m.DisposableService1{Name: "1"}, ctx
		})
		RegisterLazyFunc(Transient, func(ctx context.Context) (*m.DisposableService2, context.Context) {
			return &m.DisposableService2{Name: "2"}, ctx
		})
		assert2.PanicsWithError(t, assert2.ErrorStartsWith("detect lifetime misalignment"), Validate)
	})
	t.Run("Singleton depends on Transient", func(t *testing.T) {
		clearAll()
		RegisterLazyFunc(Singleton, func(ctx context.Context) (*m.DisposableService1, context.Context) {
			_, ctx = Get[*m.DisposableService2](ctx) //1 depends on 2
			return &m.DisposableService1{Name: "1"}, ctx
		})
		RegisterLazyFunc(Singleton, func(ctx context.Context) (*m.DisposableService2, context.Context) {
			_, ctx = Get[*m.DisposableService3](ctx) //2 depends on 3
			return &m.DisposableService2{Name: "2"}, ctx
		})
		RegisterLazyFunc(Transient, func(ctx context.Context) (*m.DisposableService3, context.Context) {
			return &m.DisposableService3{Name: "3"}, ctx
		})
		assert2.PanicsWithError(t, assert2.ErrorStartsWith("detect lifetime misalignment"), Validate)
	})
}

func TestValidate_MissingDependency(t *testing.T) {
	clearAll()
	RegisterLazyFunc(Transient, func(ctx context.Context) (*m.DisposableService1, context.Context) {
		_, ctx = Get[*m.DisposableService2](ctx) //1 depends on 2
		return &m.DisposableService1{Name: "1"}, ctx
	})
	RegisterLazyFunc(Scoped, func(ctx context.Context) (*m.DisposableService2, context.Context) {
		_, ctx = Get[*m.DisposableService3](ctx) //2 depends on 3
		return &m.DisposableService2{Name: "2"}, ctx
	})
	RegisterLazyFunc(Singleton, func(ctx context.Context) (*m.DisposableService3, context.Context) {
		_, ctx = Get[*m.DisposableService4](ctx) //3 depends on 4
		return &m.DisposableService3{Name: "3"}, ctx
	})
	//forget to register 4
	assert2.PanicsWithError(t, assert2.ErrorStartsWith("implementation not found for type"), Validate)
}
