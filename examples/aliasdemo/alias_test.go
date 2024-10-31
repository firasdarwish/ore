package aliasdemo

import (
	"context"
	"testing"

	"github.com/firasdarwish/ore"
)

func TestGetInterfaceAliasWithKeys(t *testing.T) {
	ore.RegisterLazyFunc(ore.Scoped, func(ctx context.Context) (*Broker, context.Context) {
		return &Broker{Name: "Peter1"}, ctx
	}, "module1")
	ore.RegisterLazyFunc(ore.Scoped, func(ctx context.Context) (*Broker, context.Context) {
		return &Broker{Name: "John1"}, ctx
	}, "module1")
	ore.RegisterLazyFunc(ore.Scoped, func(ctx context.Context) (*Trader, context.Context) {
		return &Trader{Name: "Mary1"}, ctx
	}, "module1")

	ore.RegisterLazyFunc(ore.Scoped, func(ctx context.Context) (*Broker, context.Context) {
		return &Broker{Name: "John2"}, ctx
	}, "module2")
	ore.RegisterLazyFunc(ore.Scoped, func(ctx context.Context) (*Trader, context.Context) {
		return &Trader{Name: "Mary2"}, ctx
	}, "module2")

	ore.RegisterLazyFunc(ore.Scoped, func(ctx context.Context) (*Trader, context.Context) {
		return &Trader{Name: "Mary3"}, ctx
	}, "module3")

	ore.RegisterAlias[IPerson, *Trader]() //link IPerson to *Trader
	ore.RegisterAlias[IPerson, *Broker]() //link IPerson to *Broker

	ctx := context.Background()

	//no IPerson was registered to the container, but we can still `Get` it.
	//(1) IPerson is alias to both *Broker and *Trader. *Broker takes precedence because it's the last one linked to IPerson.
	//(2) multiple *Borker (Peter and John) are registered to the container, the last registered (John) takes precedence.
	person1, ctx := ore.Get[IPerson](ctx, "module1") // will return the broker John
	switch person := person1.(type) {
	case *Broker:
		if person.Name != "John1" {
			t.Errorf("got %v, expected %v", person.Name, "John1")
		}
	case *Trader:
		t.Errorf("got Trader, expected Broker")
	}

	personList1, ctx := ore.GetList[IPerson](ctx, "module1") // will return all registered broker and trader
	if len(personList1) != 3 {
		t.Errorf("got %v, expected %v", len(personList1), 3)
	}

	person2, ctx := ore.Get[IPerson](ctx, "module2") // will return the broker John
	if person2.(*Broker).Name != "John2" {
		t.Errorf("got %v, expected %v", person2.(*Broker).Name, "John2")
	}

	personList2, ctx := ore.GetList[IPerson](ctx, "module2") // will return all registered broker and trader
	if len(personList2) != 2 {
		t.Errorf("got %v, expected %v", len(personList2), 2)
	}

	person3, ctx := ore.Get[IPerson](ctx, "module3") // will return the trader Mary
	if person3.(*Trader).Name != "Mary3" {
		t.Errorf("got %v, expected %v", person3.(*Trader).Name, "Mary3")
	}

	personList3, ctx := ore.GetList[IPerson](ctx, "module3") // will return all registered broker and trader
	if len(personList3) != 1 {
		t.Errorf("got %v, expected %v", len(personList3), 1)
	}

	personListNoModule, _ := ore.GetList[IPerson](ctx) // will return all registered broker and trader without keys
	if len(personListNoModule) != 0 {
		t.Errorf("got %v, expected %v", len(personListNoModule), 0)
	}
}

// func TestGetInterfaceAliasWithDifferentScope(t *testing.T) {
// 	module := "TestGetInterfaceAliasWithDifferentScope"
// 	ore.RegisterLazyFunc(ore.Transient, func(ctx context.Context) (*Broker, context.Context) {
// 		return &Broker{Name: "Transient"}, ctx
// 	}, module)
// 	ore.RegisterLazyFunc(ore.Singleton, func(ctx context.Context) (*Broker, context.Context) {
// 		return &Broker{Name: "Singleton"}, ctx
// 	}, module)
// 	ore.RegisterLazyFunc(ore.Scoped, func(ctx context.Context) (*Broker, context.Context) {
// 		return &Broker{Name: "Scoped"}, ctx
// 	}, module)
// 	ore.RegisterAlias[IPerson, *Broker]() //link IPerson to *Broker

// 	ctx := context.Background()

// 	person, ctx := ore.Get[IPerson](ctx, module)
// 	if person.(*Broker).Name != "Scoped" {
// 		t.Errorf("got %v, expected %v", person.(*Broker).Name, "Scoped")
// 	}

// 	personList, _ := ore.GetList[IPerson](ctx, module)
// 	if len(personList) != 2 {
// 		t.Errorf("got %v, expected %v", len(personList), 2)
// 	}
// }

type IPerson interface{}
type Broker struct {
	Name string
} //implements IPerson

type Trader struct {
	Name string
} //implements IPerson
