package ore

import (
	"context"
	"github.com/firasdarwish/ore/internal/interfaces"
	"testing"

	m "github.com/firasdarwish/ore/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestAliasResolverConflict(t *testing.T) {
	clearAll()
	RegisterFunc(Singleton, func(ctx context.Context) (m.IPerson, context.Context) {
		return &m.Trader{Name: "Peter Singleton"}, ctx
	})
	RegisterFunc(Transient, func(ctx context.Context) (*m.Broker, context.Context) {
		return &m.Broker{Name: "Mary Transient"}, ctx
	})

	RegisterAlias[m.IPerson, *m.Trader]()
	RegisterAlias[m.IPerson, *m.Broker]()

	ctx := context.Background()

	//The last registered IPerson is "Mary Transient", it would normally takes precedence.
	//However, we registered a direct resolver for IPerson which is "Peter Singleton".
	//So Ore won't treat IPerson as an alias and will resolve IPerson directly as "Peter Singleton"
	person, ctx := Get[m.IPerson](ctx)
	assert.Equal(t, person.(*m.Trader).Name, "Peter Singleton")

	//GetList will return all possible IPerson whatever alias or from direct resolver.
	personList, _ := GetList[m.IPerson](ctx)
	assert.Equal(t, len(personList), 2)
}

func TestAliasOfAliasIsNotAllow(t *testing.T) {
	clearAll()
	RegisterFunc(Singleton, func(ctx context.Context) (*m.Trader, context.Context) {
		return &m.Trader{Name: "Peter Singleton"}, ctx
	})
	RegisterFunc(Transient, func(ctx context.Context) (*m.Broker, context.Context) {
		return &m.Broker{Name: "Mary Transient"}, ctx
	})

	RegisterAlias[m.IPerson, *m.Trader]()
	RegisterAlias[m.IPerson, *m.Broker]()
	RegisterAlias[m.IHuman, m.IPerson]() //alias of alias

	assert.Panics(t, func() {
		_, _ = Get[m.IHuman](context.Background())
	}, "implementation not found for type: IHuman")

	humans, _ := GetList[m.IHuman](context.Background())
	assert.Empty(t, humans)
}

func TestAliasWithDifferentScope(t *testing.T) {
	clearAll()
	module := "TestGetInterfaceAliasWithDifferentScope"
	RegisterKeyedFunc(Transient, func(ctx context.Context) (*m.Broker, context.Context) {
		return &m.Broker{Name: "Transient"}, ctx
	}, module)
	RegisterKeyedFunc(Singleton, func(ctx context.Context) (*m.Broker, context.Context) {
		return &m.Broker{Name: "Singleton"}, ctx
	}, module)
	RegisterKeyedFunc(Scoped, func(ctx context.Context) (*m.Broker, context.Context) {
		return &m.Broker{Name: "Scoped"}, ctx
	}, module)
	RegisterAlias[m.IPerson, *m.Broker]() //link m.IPerson to *m.Broker

	ctx := context.Background()

	person, ctx := GetKeyed[m.IPerson](ctx, module)
	assert.Equal(t, person.(*m.Broker).Name, "Scoped")

	personList, _ := GetKeyedList[m.IPerson](ctx, module)
	assert.Equal(t, len(personList), 3)
}

func TestAliasIsScopedByKeys(t *testing.T) {
	clearAll()
	RegisterKeyedFunc(Scoped, func(ctx context.Context) (*m.Broker, context.Context) {
		return &m.Broker{Name: "Peter1"}, ctx
	}, "module1")
	RegisterKeyedFunc(Scoped, func(ctx context.Context) (*m.Broker, context.Context) {
		return &m.Broker{Name: "John1"}, ctx
	}, "module1")
	RegisterKeyedFunc(Scoped, func(ctx context.Context) (*m.Trader, context.Context) {
		return &m.Trader{Name: "Mary1"}, ctx
	}, "module1")

	RegisterKeyedFunc(Scoped, func(ctx context.Context) (*m.Broker, context.Context) {
		return &m.Broker{Name: "John2"}, ctx
	}, "module2")
	RegisterKeyedFunc(Scoped, func(ctx context.Context) (*m.Trader, context.Context) {
		return &m.Trader{Name: "Mary2"}, ctx
	}, "module2")

	RegisterKeyedFunc(Scoped, func(ctx context.Context) (*m.Trader, context.Context) {
		return &m.Trader{Name: "Mary3"}, ctx
	}, "module3")

	RegisterAlias[m.IPerson, *m.Trader]() //link m.IPerson to *m.Trader
	RegisterAlias[m.IPerson, *m.Broker]() //link m.IPerson to *m.Broker

	ctx := context.Background()

	person1, ctx := GetKeyed[m.IPerson](ctx, "module1") // will return the m.Broker John
	assert.Equal(t, person1.(*m.Broker).Name, "John1")

	personList1, ctx := GetKeyedList[m.IPerson](ctx, "module1") // will return all registered m.Broker and m.Trader
	assert.Equal(t, len(personList1), 3)

	person2, ctx := GetKeyed[m.IPerson](ctx, "module2") // will return the m.Broker John
	assert.Equal(t, person2.(*m.Broker).Name, "John2")

	personList2, ctx := GetKeyedList[m.IPerson](ctx, "module2") // will return all registered m.Broker and m.Trader
	assert.Equal(t, len(personList2), 2)

	person3, ctx := GetKeyed[m.IPerson](ctx, "module3") // will return the m.Trader Mary
	assert.Equal(t, person3.(*m.Trader).Name, "Mary3")

	personList3, ctx := GetKeyedList[m.IPerson](ctx, "module3") // will return all registered m.Broker and m.Trader
	assert.Equal(t, len(personList3), 1)

	personListNoModule, _ := GetList[m.IPerson](ctx) // will return all registered m.Broker and m.Trader without keys
	assert.Empty(t, personListNoModule)
}

func TestInvalidAlias(t *testing.T) {
	assert.Panics(t, func() {
		RegisterAlias[error, *m.Broker]() //register a struct (Broker) that does not implement interface (error)
	}, "Broker does not implements error")
}

func TestGetGenericAlias(t *testing.T) {
	for _, registrationType := range types {
		container := NewContainer()

		RegisterFuncToContainer(container, registrationType, func(ctx context.Context) (*m.SimpleCounterUint, context.Context) {
			return &m.SimpleCounterUint{}, ctx
		})
		RegisterAliasToContainer[interfaces.SomeCounterGeneric[uint], *m.SimpleCounterUint](container)

		c, _ := GetFromContainer[interfaces.SomeCounterGeneric[uint]](container, context.Background())

		c.Add(1)
		c.Add(1)

		assert.Equal(t, uint(2), c.GetCount())
	}
}

func TestGetListGenericAlias(t *testing.T) {
	for _, registrationType := range types {
		container := NewContainer()

		for i := 0; i < 3; i++ {
			RegisterFuncToContainer(container, registrationType, func(ctx context.Context) (*m.SimpleCounterUint, context.Context) {
				return &m.SimpleCounterUint{}, ctx
			})
		}

		RegisterAliasToContainer[interfaces.SomeCounterGeneric[uint], *m.SimpleCounterUint](container)

		counters, _ := GetListFromContainer[interfaces.SomeCounterGeneric[uint]](container, context.Background())
		assert.Equal(t, len(counters), 3)

		c := counters[1]
		c.Add(1)
		c.Add(1)

		assert.Equal(t, uint(2), c.GetCount())
	}
}
