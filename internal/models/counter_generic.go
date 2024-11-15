package models

import (
	"context"
	"github.com/firasdarwish/ore/internal/interfaces"
)

type CounterGeneric[T interfaces.Numeric] struct {
	Counter T
}

func (c *CounterGeneric[T]) Add(number T) {
	c.Counter += number
}

func (c *CounterGeneric[T]) GetCount() T {
	return c.Counter
}

func (c *CounterGeneric[T]) New(ctx context.Context) (interfaces.SomeCounterGeneric[T], context.Context) {
	return &CounterGeneric[T]{}, ctx
}
