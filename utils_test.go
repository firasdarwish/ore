package ore

import (
	"context"
	"testing"
)

var types = []Lifetime{Singleton, Transient, Scoped}

func mustHavePanicked(t *testing.T) {
	if r := recover(); r == nil {
		t.Errorf("Expected panic")
	}
}

type someCounter interface {
	AddOne()
	GetCount() int
}

type numeric interface {
	uint
}

type someCounterGeneric[T numeric] interface {
	Add(number T)
	GetCount() T
}

type simpleCounter struct {
	counter int
}

func (c *simpleCounter) AddOne() {
	c.counter++
}

func (c *simpleCounter) GetCount() int {
	return c.counter
}

func (c *simpleCounter) New(ctx context.Context) (someCounter, context.Context) {
	return &simpleCounter{}, ctx
}

type simpleCounter2 struct {
	counter int
}

func (c *simpleCounter2) AddOne() {
	c.counter++
}

func (c *simpleCounter2) GetCount() int {
	return c.counter
}

func (c *simpleCounter2) New(ctx context.Context) (someCounter, context.Context) {
	return &simpleCounter2{}, ctx
}

type counterGeneric[T numeric] struct {
	counter T
}

func (c *counterGeneric[T]) Add(number T) {
	c.counter += number
}

func (c *counterGeneric[T]) GetCount() T {
	return c.counter
}

func (c *counterGeneric[T]) New(ctx context.Context) (someCounterGeneric[T], context.Context) {
	return &counterGeneric[T]{}, ctx
}

type c struct {
	Counter int
}

type a struct {
	C *c
}

func (*a) New(ctx context.Context) (*a, context.Context) {
	ccc, ctx := Get[*c](ctx)

	return &a{
		C: ccc,
	}, ctx
}

func (*c) New(ctx context.Context) (*c, context.Context) {
	return &c{}, ctx
}
