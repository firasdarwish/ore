package aliasdemo

import (
	"context"
	"testing"

	"github.com/firasdarwish/ore"
)

func TestGetInterfaceAliasWithKeys(t *testing.T) {
	ore.RegisterLazyFunc(ore.Scoped, func(ctx context.Context) (*broker, context.Context) {
		return &broker{Name: "Peter1"}, ctx
	}, "module1")
	ore.RegisterLazyFunc(ore.Scoped, func(ctx context.Context) (*broker, context.Context) {
		return &broker{Name: "John1"}, ctx
	}, "module1")
	ore.RegisterLazyFunc(ore.Scoped, func(ctx context.Context) (*trader, context.Context) {
		return &trader{Name: "Mary1"}, ctx
	}, "module1")

	ore.RegisterLazyFunc(ore.Scoped, func(ctx context.Context) (*broker, context.Context) {
		return &broker{Name: "John2"}, ctx
	}, "module2")
	ore.RegisterLazyFunc(ore.Scoped, func(ctx context.Context) (*trader, context.Context) {
		return &trader{Name: "Mary2"}, ctx
	}, "module2")

	ore.RegisterLazyFunc(ore.Scoped, func(ctx context.Context) (*trader, context.Context) {
		return &trader{Name: "Mary3"}, ctx
	}, "module3")

	ore.RegisterAlias[iPerson, *trader]() //link IPerson to *Trader
	ore.RegisterAlias[iPerson, *broker]() //link IPerson to *Broker

	ctx := context.Background()

	//no IPerson was registered to the container, but we can still `Get` it.
	//(1) IPerson is alias to both *Broker and *Trader. *Broker takes precedence because it's the last one linked to IPerson.
	//(2) multiple *Borker (Peter and John) are registered to the container, the last registered (John) takes precedence.
	person1, ctx := ore.Get[iPerson](ctx, "module1") // will return the broker John
	switch person := person1.(type) {
	case *broker:
		if person.Name != "John1" {
			t.Errorf("got %v, expected %v", person.Name, "John1")
		}
	case *trader:
		t.Errorf("got Trader, expected Broker")
	}

	personList1, ctx := ore.GetList[iPerson](ctx, "module1") // will return all registered broker and trader
	if len(personList1) != 3 {
		t.Errorf("got %v, expected %v", len(personList1), 3)
	}

	person2, ctx := ore.Get[iPerson](ctx, "module2") // will return the broker John
	if person2.(*broker).Name != "John2" {
		t.Errorf("got %v, expected %v", person2.(*broker).Name, "John2")
	}

	personList2, ctx := ore.GetList[iPerson](ctx, "module2") // will return all registered broker and trader
	if len(personList2) != 2 {
		t.Errorf("got %v, expected %v", len(personList2), 2)
	}

	person3, ctx := ore.Get[iPerson](ctx, "module3") // will return the trader Mary
	if person3.(*trader).Name != "Mary3" {
		t.Errorf("got %v, expected %v", person3.(*trader).Name, "Mary3")
	}

	personList3, ctx := ore.GetList[iPerson](ctx, "module3") // will return all registered broker and trader
	if len(personList3) != 1 {
		t.Errorf("got %v, expected %v", len(personList3), 1)
	}

	personListNoModule, _ := ore.GetList[iPerson](ctx) // will return all registered broker and trader without keys
	if len(personListNoModule) != 0 {
		t.Errorf("got %v, expected %v", len(personListNoModule), 0)
	}
}

type iPerson interface{}
type broker struct {
	Name string
} //implements IPerson

type trader struct {
	Name string
} //implements IPerson
