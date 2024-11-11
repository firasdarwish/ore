package main

import (
	"context"
	"fmt"

	"github.com/firasdarwish/ore"
)

func main() {
	ore.RegisterLazyFunc[Counter](ore.Singleton, func(ctx context.Context) (Counter, context.Context) {
		fmt.Println("NEWLY INITIALIZED FROM FUNC")
		return &mycounter{}, ctx
	}, "firas")

	ore.RegisterLazyFunc[Counter](ore.Singleton, func(ctx context.Context) (Counter, context.Context) {
		fmt.Println("NEWLY INITIALIZED FROM FUNC")
		return &mycounter{}, ctx
	}, "darwish")

	ore.RegisterLazyCreator[Counter](ore.Singleton, &mycounter{})

	cc := &mycounter{}
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

	gc, _ = ore.Get[GenericCounter[uint]](ctx)
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

type mycounter struct {
	count int
}

var _ Counter = (*mycounter)(nil)

func (c *mycounter) AddOne() {
	c.count++
}

func (c *mycounter) Total() int {
	return c.count
}

func (*mycounter) New(ctx context.Context) (Counter, context.Context) {
	fmt.Println("NEWLY INITIALIZED")
	return &mycounter{}, ctx
}

type numeric interface {
	uint
}

type genCounter[T numeric] struct {
	count T
}

func (t *genCounter[T]) Add(num T) {
	t.count += num
}

func (t *genCounter[T]) Total() T {
	return t.count
}

func (*genCounter[T]) New(ctx context.Context) (GenericCounter[T], context.Context) {
	return &genCounter[T]{}, ctx
}
