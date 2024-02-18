package ore

import (
	"context"
	"io"
	"strconv"
	"testing"
)

var types = []RegistrationType{Singleton, Transient, Scoped}

func mustHavePanicked(t *testing.T) {
	if r := recover(); r == nil {
		t.Errorf("Expected panic")
	}
}

type Counter interface {
	AddOne()
	GetCount() int
}

type CounterWriter interface {
	Add(number int)
	GetCount() int
}

type numeric interface {
	uint
}

type CounterGeneric[T numeric] interface {
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

func (c *simpleCounter) New(ctx context.Context) Counter {
	return &simpleCounter{}
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

func (c *simpleCounter2) New(ctx context.Context) Counter {
	return &simpleCounter2{}
}

type counterWriter struct {
	counter int
	writer  io.Writer
}

func (c *counterWriter) Add(number int) {
	_, _ = c.writer.Write([]byte("New Number Added: " + strconv.Itoa(number)))
	c.counter += number
}

func (c *counterWriter) GetCount() int {
	_, _ = c.writer.Write([]byte("Total Count: " + strconv.Itoa(c.counter)))
	return c.counter
}

func (c *counterWriter) New(ctx context.Context) CounterWriter {

	writer, _ := Get[io.Writer](ctx)

	return &counterWriter{
		writer: writer,
	}
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

func (c *counterGeneric[T]) New(ctx context.Context) CounterGeneric[T] {
	return &counterGeneric[T]{}
}
