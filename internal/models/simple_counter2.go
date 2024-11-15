package models

import (
	"context"
	"github.com/firasdarwish/ore/internal/interfaces"
)

type SimpleCounter2 struct {
	Counter int
}

func (c *SimpleCounter2) AddOne() {
	c.Counter++
}

func (c *SimpleCounter2) GetCount() int {
	return c.Counter
}

func (c *SimpleCounter2) New(ctx context.Context) (interfaces.SomeCounter, context.Context) {
	return &SimpleCounter2{}, ctx
}
