package models

import (
	"context"
	"github.com/firasdarwish/ore/internal/interfaces"
)

type SimpleCounter struct {
	Counter int
}

func (c *SimpleCounter) AddOne() {
	c.Counter++
}

func (c *SimpleCounter) GetCount() int {
	return c.Counter
}

func (c *SimpleCounter) New(ctx context.Context) (interfaces.SomeCounter, context.Context) {
	return &SimpleCounter{}, ctx
}
