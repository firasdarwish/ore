package main

import (
	"context"
	"fmt"
)

type counter struct {
	count int
}

func (c *counter) AddOne() {
	c.count++
}

func (c *counter) Total() int {
	return c.count
}

func (*counter) New(ctx context.Context) Counter {
	fmt.Println("NEWLY INITIALIZED")
	return &counter{}
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

func (*genCounter[T]) New(ctx context.Context) GenericCounter[T] {
	return &genCounter[T]{}
}
