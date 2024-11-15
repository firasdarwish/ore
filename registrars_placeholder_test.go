package ore

import (
	"context"
	"testing"

	m "github.com/firasdarwish/ore/internal/models"
	"github.com/firasdarwish/ore/internal/testtools/assert2"
	"github.com/stretchr/testify/assert"
)

func TestPlaceHolder_HappyPath(t *testing.T) {
	clearAll()

	//register a placeholder
	RegisterPlaceholder[*m.Trader]()

	//get the placeholder value would fail
	assert2.PanicsWithError(t, assert2.ErrorStartsWith("no value has been provided for this placeholder"), func() {
		_, _ = Get[*m.Trader](context.Background())
	})

	//get list would return empty
	traders, _ := GetList[*m.Trader](context.Background())
	assert.Empty(t, traders)

	//provide a value to the placeholder
	ctx := ProvideScopedValue[*m.Trader](context.Background(), &m.Trader{Name: "Peter"})

	//get the placeholder value would success
	trader, _ := Get[*m.Trader](ctx)
	assert.Equal(t, "Peter", trader.Name)

	//get list will include the placeholder value
	traders, ctx = GetList[*m.Trader](ctx)
	assert.Equal(t, 1, len(traders))

	//Register alias to the placeholder
	RegisterAlias[m.IPerson, *m.Trader]()

	//get the alias value would success
	person, ctx := Get[m.IPerson](ctx)
	assert.Equal(t, "Peter", person.(*m.Trader).Name)

	//get list will include the placeholder value
	persons, _ := GetList[m.IPerson](ctx)
	assert.Equal(t, 1, len(persons))
}

func TestPlaceHolder_ProvideValueBeforeRegistering(t *testing.T) {
	clearAll()

	//provide a value to the placeholder
	ctx := ProvideScopedValue[*m.Trader](context.Background(), &m.Trader{Name: "Mary"})

	//get the placeholder value would fail because no placeholder has been registered
	assert2.PanicsWithError(t, assert2.ErrorStartsWith("implementation not found for type"), func() {
		_, _ = Get[*m.Trader](ctx)
	})

	//register a matching placeholder
	RegisterPlaceholder[*m.Trader]()

	//get the placeholder value would success
	trader, _ := Get[*m.Trader](ctx)
	assert.Equal(t, "Mary", trader.Name)
}

// can not register a placeholder to override a real resolver
func TestPlaceHolder_OverrideRealResolver(t *testing.T) {
	clearAll()

	//register a real resolver
	RegisterKeyedSingleton(&m.Trader{Name: "Mary"}, "module1")

	//register a placeholder to override the real resolver should fail
	assert2.PanicsWithError(t, assert2.ErrorContains("has already been registered"), func() {
		RegisterKeyedPlaceholder[*m.Trader]("module1")
	})

	//register 2 time a placeholder should fail
	RegisterKeyedPlaceholder[*m.Trader]("module2")
	assert2.PanicsWithError(t, assert2.ErrorContains("has already been registered"), func() {
		RegisterKeyedPlaceholder[*m.Trader]("module2")
	})
}

func TestPlaceHolder_OverridePlaceHolder(t *testing.T) {
	clearAll()
	//register a placeholder
	RegisterKeyedPlaceholder[*m.Trader]("module2")

	//Provide the value to the placeholder
	ctx := ProvideKeyedScopedValue[*m.Trader](context.Background(), &m.Trader{Name: "John"}, "module2")

	//get the placeholder value would success
	trader, ctx := GetKeyed[*m.Trader](ctx, "module2")
	assert.Equal(t, "John", trader.Name)

	//replace the placeholder value "John" with a new value "David"
	ctx = ProvideKeyedScopedValue[*m.Trader](ctx, &m.Trader{Name: "David"}, "module2")
	trader, ctx = GetKeyed[*m.Trader](ctx, "module2")
	assert.Equal(t, "David", trader.Name)

	traders, ctx := GetKeyedList[*m.Trader](ctx, "module2")
	assert.Equal(t, 1, len(traders))
	assert.Equal(t, "David", traders[0].Name)

	//Register a real resolver should override the placeholder resolver
	RegisterKeyedFunc(Singleton, func(ctx context.Context) (*m.Trader, context.Context) {
		return &m.Trader{Name: "Mary"}, ctx
	}, "module2")

	trader, ctx = GetKeyed[*m.Trader](ctx, "module2")
	assert.Equal(t, "Mary", trader.Name)

	//Get both the placeholder value ("David") and the real resolver value ("Mary")
	traders, ctx = GetKeyedList[*m.Trader](ctx, "module2")
	assert.Equal(t, 2, len(traders)) //David and Mary
	assert.True(t, tradersListContainsName(traders, "David"))
	assert.True(t, tradersListContainsName(traders, "Mary"))

	//replace the placeholder value ("David") with a new value ("Nathan")
	ctx = ProvideKeyedScopedValue[*m.Trader](ctx, &m.Trader{Name: "Nathan"}, "module2")

	//the placeholder value cannot override the real resolver value
	trader, ctx = GetKeyed[*m.Trader](ctx, "module2")
	assert.Equal(t, "Mary", trader.Name)

	//but it replaces the old placeholder value ("Nathan" will replace "David")
	traders, _ = GetKeyedList[*m.Trader](ctx, "module2")
	assert.Equal(t, 2, len(traders)) //Nathan and Mary
	assert.True(t, tradersListContainsName(traders, "Nathan"))
	assert.True(t, tradersListContainsName(traders, "Mary"))
}

// placeholder value of a module is not accessible from other module
func TestPlaceHolder_PerModule(t *testing.T) {
	con1 := NewContainer()
	RegisterPlaceholderToContainer[*m.Trader](con1)

	con2 := NewContainer()
	RegisterPlaceholderToContainer[*m.Trader](con2)

	ctx := ProvideScopedValueToContainer(con1, context.Background(), &m.Trader{Name: "John"})
	trader, ctx := GetFromContainer[*m.Trader](con1, ctx)
	assert.Equal(t, "John", trader.Name)

	assert2.PanicsWithError(t, assert2.ErrorStartsWith("no value has been provided for this placeholder"), func() {
		trader, ctx = GetFromContainer[*m.Trader](con2, ctx)
	})
}

func tradersListContainsName(p []*m.Trader, name string) bool {
	for _, v := range p {
		if v.Name == name {
			return true
		}
	}
	return false
}
