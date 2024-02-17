package main

import (
	"context"
	"fmt"
	"github.com/firasdarwish/ore"
)

func main() {

	ore.RegisterLazyFunc[Counter](ore.Singleton, func(ctx context.Context) Counter {
		fmt.Println("NEWLY INITIALIZED FROM FUNC")
		return &counter{}
	}, "firas")

	ore.RegisterLazyFunc[Counter](ore.Singleton, func(ctx context.Context) Counter {
		fmt.Println("NEWLY INITIALIZED FROM FUNC")
		return &counter{}
	}, "darwish")

	ore.RegisterLazyCreator[Counter](ore.Singleton, &counter{})

	var cc Counter
	cc = &counter{}
	ore.RegisterEagerSingleton[Counter](cc)

	ctx := context.Background()

	fmt.Println("STARTED ...")

	c, ctx := ore.Get[Counter](ctx, "firas")
	c.AddOne()
	c.AddOne()

	c, ctx = ore.Get[Counter](ctx, "darwish")
	c.AddOne()
	c.AddOne()

	fmt.Printf("Total Count: %v", c.Total())

	ore.RegisterLazyCreator[GenericCounter[uint]](ore.Scoped, &genCounter[uint]{})

	gc, ctx := ore.Get[GenericCounter[uint]](ctx)
	gc.Add(1)
	gc.Add(1)
	gc.Add(1)
	gc.Add(1)

	gc, ctx = ore.Get[GenericCounter[uint]](ctx)
	gc.Add(1)

	fmt.Println(gc.Total())
}

type Counter interface {
	AddOne()
	Total() int
}

type GenericCounter[T numeric] interface {
	Add(num T)
	Total() T
}
