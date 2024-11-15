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
			RegisterFunc(lt, func(ctx context.Context) (*m.DisposableService1, context.Context) {
				_, ctx = Get[*m.DisposableService1](ctx) //1 calls 1
				return &m.DisposableService1{Name: "1"}, ctx
			})
			assert2.PanicsWithError(t, assert2.ErrorStartsWith("detected cyclic dependency"), Validate)
		})
		t.Run("Indirect circular "+lt.String()+" (1 calls 2 calls 3 calls 1)", func(t *testing.T) {
			clearAll()
			RegisterFunc(lt, func(ctx context.Context) (*m.DisposableService1, context.Context) {
				_, ctx = Get[*m.DisposableService2](ctx) //1 calls 2
				return &m.DisposableService1{Name: "1"}, ctx
			})
			RegisterFunc(lt, func(ctx context.Context) (*m.DisposableService2, context.Context) {
				_, ctx = Get[*m.DisposableService3](ctx) //2 calls 3
				return &m.DisposableService2{Name: "2"}, ctx
			})
			RegisterFunc(lt, func(ctx context.Context) (*m.DisposableService3, context.Context) {
				_, ctx = Get[*m.DisposableService1](ctx) //3 calls 1
				return &m.DisposableService3{Name: "3"}, ctx
			})
			assert2.PanicsWithError(t, assert2.ErrorStartsWith("detected cyclic dependency"), Validate)
		})
		t.Run("Middle circular "+lt.String()+" (1 calls 2 calls 3 calls 4 calls 2)", func(t *testing.T) {
			clearAll()
			RegisterFunc(lt, func(ctx context.Context) (*m.DisposableService1, context.Context) {
				_, ctx = Get[*m.DisposableService2](ctx) //1 calls 2
				return &m.DisposableService1{Name: "1"}, ctx
			})
			RegisterFunc(lt, func(ctx context.Context) (*m.DisposableService2, context.Context) {
				_, ctx = Get[*m.DisposableService3](ctx) //2 calls 3
				return &m.DisposableService2{Name: "2"}, ctx
			})
			RegisterFunc(lt, func(ctx context.Context) (*m.DisposableService3, context.Context) {
				_, ctx = Get[*m.DisposableService4](ctx) //3 calls 4
				return &m.DisposableService3{Name: "3"}, ctx
			})
			RegisterFunc(lt, func(ctx context.Context) (*m.DisposableService4, context.Context) {
				_, ctx = Get[*m.DisposableService2](ctx) //4 calls 2
				return &m.DisposableService4{Name: "4"}, ctx
			})
			assert2.PanicsWithError(t, assert2.ErrorStartsWith("detected cyclic dependency"), Validate)
		})
		t.Run("circular on complex tree "+lt.String()+"", func(t *testing.T) {
			clearAll()
			RegisterFunc(lt, func(ctx context.Context) (*m.DisposableService1, context.Context) {
				_, ctx = Get[*m.DisposableService2](ctx) //1 calls 2
				_, ctx = Get[*m.DisposableService3](ctx) //1 calls 3
				return &m.DisposableService1{Name: "1"}, ctx
			})
			RegisterFunc(lt, func(ctx context.Context) (*m.DisposableService2, context.Context) {
				_, ctx = Get[*m.DisposableService4](ctx) //2 calls 4
				_, ctx = Get[*m.DisposableService5](ctx) //2 calls 5
				return &m.DisposableService2{Name: "2"}, ctx
			})
			RegisterFunc(lt, func(ctx context.Context) (*m.DisposableService3, context.Context) {
				_, ctx = Get[*m.DisposableService4](ctx) //3 calls 4
				return &m.DisposableService3{Name: "3"}, ctx
			})
			RegisterFunc(lt, func(ctx context.Context) (*m.DisposableService4, context.Context) {
				_, ctx = Get[*m.DisposableService5](ctx) //4 calls 5
				return &m.DisposableService4{Name: "4"}, ctx
			})
			RegisterFunc(lt, func(ctx context.Context) (*m.DisposableService5, context.Context) {
				_, ctx = Get[*m.DisposableService3](ctx) //5 calls 3 => circular here: 5->3->4->5
				return &m.DisposableService5{Name: "5"}, ctx
			})
			assert2.PanicsWithError(t, assert2.ErrorStartsWith("detected cyclic dependency"), Validate)
		})
		t.Run("fake circular top down "+lt.String()+": (1 calls 2 (x2) calls 3 calls 4, 2 calls 4)", func(t *testing.T) {
			clearAll()
			RegisterFunc(lt, func(ctx context.Context) (*m.DisposableService1, context.Context) {
				_, ctx = Get[*m.DisposableService2](ctx) //1 calls 2
				_, ctx = Get[*m.DisposableService2](ctx) //1 calls 2 again
				return &m.DisposableService1{Name: "1"}, ctx
			})
			RegisterFunc(lt, func(ctx context.Context) (*m.DisposableService2, context.Context) {
				_, ctx = Get[*m.DisposableService3](ctx) //2 calls 3
				_, ctx = Get[*m.DisposableService4](ctx) //2 calls 4
				return &m.DisposableService2{Name: "2"}, ctx
			})
			RegisterFunc(lt, func(ctx context.Context) (*m.DisposableService3, context.Context) {
				_, ctx = Get[*m.DisposableService4](ctx) //3 calls 4
				_, ctx = Get[*m.DisposableService4](ctx) //3 calls 4
				return &m.DisposableService3{Name: "3"}, ctx
			})
			RegisterFunc(lt, func(ctx context.Context) (*m.DisposableService4, context.Context) {
				return &m.DisposableService4{Name: "4"}, ctx
			})
			assert.NotPanics(t, Validate)
		})
		t.Run("fake circular sibling "+lt.String()+": 1 calls 2 & 3;  2 calls 3)", func(t *testing.T) {
			clearAll()
			RegisterFunc(lt, func(ctx context.Context) (*m.DisposableService1, context.Context) {
				_, ctx = Get[*m.DisposableService2](ctx) //1 calls 2
				_, ctx = Get[*m.DisposableService3](ctx) //1 calls 3
				return &m.DisposableService1{Name: "1"}, ctx
			})
			RegisterFunc(lt, func(ctx context.Context) (*m.DisposableService2, context.Context) {
				_, ctx = Get[*m.DisposableService3](ctx) //2 calls 3
				return &m.DisposableService2{Name: "2"}, ctx
			})
			RegisterFunc(lt, func(ctx context.Context) (*m.DisposableService3, context.Context) {
				return &m.DisposableService3{Name: "3"}, ctx
			})
			assert.NotPanics(t, Validate)
		})
	}
}

func TestValidate_CircularMixedLifetype(t *testing.T) {
	clearAll()

	RegisterFunc(Scoped, func(ctx context.Context) (*m.DisposableService2, context.Context) {
		_, ctx = Get[*m.DisposableService4](ctx) //2 calls 4
		_, ctx = Get[*m.DisposableService5](ctx) //2 calls 5
		return &m.DisposableService2{Name: "2"}, ctx
	})
	RegisterFunc(Singleton, func(ctx context.Context) (*m.DisposableService3, context.Context) {
		_, ctx = Get[*m.DisposableService4](ctx) //3 calls 4
		return &m.DisposableService3{Name: "3"}, ctx
	})
	RegisterFunc(Singleton, func(ctx context.Context) (*m.DisposableService4, context.Context) {
		_, ctx = Get[*m.DisposableService5](ctx) //4 calls 5
		return &m.DisposableService4{Name: "4"}, ctx
	})
	RegisterFunc(Singleton, func(ctx context.Context) (*m.DisposableService5, context.Context) {
		_, ctx = Get[*m.DisposableService3](ctx) //5 calls 3 => circular here: 5->3->4->5
		return &m.DisposableService5{Name: "5"}, ctx
	})
	RegisterFunc(Transient, func(ctx context.Context) (*m.DisposableService1, context.Context) {
		_, ctx = Get[*m.DisposableService2](ctx) //1 calls 2
		_, ctx = Get[*m.DisposableService3](ctx) //1 calls 3
		return &m.DisposableService1{Name: "1"}, ctx
	})
	assert2.PanicsWithError(t, assert2.ErrorStartsWith("detected cyclic dependency"), Validate)
	assert2.PanicsWithError(t, assert2.ErrorStartsWith("detected cyclic dependency"), func() {
		_, _ = Get[*m.DisposableService1](context.Background())
	})
}

func TestValidate_LifetimeAlignment_SingletonCallsScoped(t *testing.T) {
	con := NewContainer()
	RegisterFuncToContainer(con, Scoped, func(ctx context.Context) (*m.DisposableService2, context.Context) {
		return &m.DisposableService2{Name: "2"}, ctx
	})
	RegisterFuncToContainer(con, Singleton, func(ctx context.Context) (*m.DisposableService1, context.Context) {
		_, ctx = GetFromContainer[*m.DisposableService2](con, ctx) //1 depends on 2
		return &m.DisposableService1{Name: "1"}, ctx
	})
	assert2.PanicsWithError(t, assert2.ErrorStartsWith("detected lifetime misalignment"), con.Validate)
}
func TestValidate_LifetimeAlignment_ScopedCallsTransient(t *testing.T) {
	con := NewContainer()
	RegisterFuncToContainer(con, Scoped, func(ctx context.Context) (*m.DisposableService1, context.Context) {
		_, ctx = GetFromContainer[*m.DisposableService2](con, ctx) //1 depends on 2
		return &m.DisposableService1{Name: "1"}, ctx
	})
	RegisterFuncToContainer(con, Transient, func(ctx context.Context) (*m.DisposableService2, context.Context) {
		return &m.DisposableService2{Name: "2"}, ctx
	})
	assert2.PanicsWithError(t, assert2.ErrorStartsWith("detected lifetime misalignment"), con.Validate)
}
func TestValidate_LifetimeAlignment_SingletonCallsTransient(t *testing.T) {
	con := NewContainer()
	RegisterFuncToContainer(con, Singleton, func(ctx context.Context) (*m.DisposableService1, context.Context) {
		_, ctx = GetFromContainer[*m.DisposableService2](con, ctx) //1 depends on 2
		return &m.DisposableService1{Name: "1"}, ctx
	})
	RegisterFuncToContainer(con, Singleton, func(ctx context.Context) (*m.DisposableService2, context.Context) {
		_, ctx = GetFromContainer[*m.DisposableService3](con, ctx) //2 depends on 3
		return &m.DisposableService2{Name: "2"}, ctx
	})
	RegisterFuncToContainer(con, Transient, func(ctx context.Context) (*m.DisposableService3, context.Context) {
		return &m.DisposableService3{Name: "3"}, ctx
	})
	assert2.PanicsWithError(t, assert2.ErrorStartsWith("detected lifetime misalignment"), con.Validate)
}

func TestValidate_MissingDependency(t *testing.T) {
	clearAll()
	RegisterFunc(Transient, func(ctx context.Context) (*m.DisposableService1, context.Context) {
		_, ctx = Get[*m.DisposableService2](ctx) //1 depends on 2
		return &m.DisposableService1{Name: "1"}, ctx
	})
	RegisterFunc(Scoped, func(ctx context.Context) (*m.DisposableService2, context.Context) {
		_, ctx = Get[*m.DisposableService3](ctx) //2 depends on 3
		return &m.DisposableService2{Name: "2"}, ctx
	})
	RegisterFunc(Singleton, func(ctx context.Context) (*m.DisposableService3, context.Context) {
		_, ctx = Get[*m.DisposableService4](ctx) //3 depends on 4
		return &m.DisposableService3{Name: "3"}, ctx
	})
	//forget to register 4
	assert2.PanicsWithError(t, assert2.ErrorStartsWith("implementation not found for type"), Validate)
}

func TestValidate_WithPlaceholder(t *testing.T) {
	con := NewContainer()
	RegisterPlaceholderToContainer[*m.Trader](con)
	assert.NotPanics(t, con.Validate)
}

func TestValidate_WithPlaceholderInterface(t *testing.T) {
	con := NewContainer()
	RegisterPlaceholderToContainer[m.IPerson](con)
	assert.NotPanics(t, con.Validate)
}

func TestValidate_DisableValidation(t *testing.T) {
	con := NewContainer()
	RegisterPlaceholderToContainer[*m.Trader](con)
	RegisterFuncToContainer(con, Singleton, func(ctx context.Context) (*m.Broker, context.Context) {
		_, ctx = GetFromContainer[*m.Trader](con, ctx)
		return &m.Broker{Name: "John"}, ctx
	})

	assert2.PanicsWithError(t, assert2.ErrorStartsWith("detected lifetime misalignment"), con.Validate)
}
