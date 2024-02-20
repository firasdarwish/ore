# ore - Generic Dependency Injection Container for Go

[![Go Report Card](https://goreportcard.com/badge/github.com/firasdarwish/ore)](https://goreportcard.com/report/github.com/firasdarwish/ore)
[![Go Reference](https://pkg.go.dev/badge/github.com/firasdarwish/ore.svg)](https://pkg.go.dev/github.com/firasdarwish/ore)
[![Mentioned in Awesome Go](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go?tab=readme-ov-file#dependency-injection)
[![Maintainability](https://api.codeclimate.com/v1/badges/3bd6f2fa4390af7c8faa/maintainability)](https://codeclimate.com/github/firasdarwish/ore/maintainability)
[![Test Coverage](https://api.codeclimate.com/v1/badges/3bd6f2fa4390af7c8faa/test_coverage)](https://codeclimate.com/github/firasdarwish/ore/test_coverage)


![ore](https://github.com/firasdarwish/ore/assets/1930361/c1426ba1-777a-43f5-8a9a-7520caa45516)


___
`ore` is a lightweight, generic & simple dependency injection (DI) container for Go.

Inspired by the principles of ASP.NET Dependency Injection, designed to facilitate
the management of object lifetimes and the inversion of control in your applications.


<br />

# Features

- **Singletons**: Register components as singletons, ensuring that there's only one instance throughout the entire
  application.


- **Transients**: Register components as transients, creating a new instance each time it is requested.


- **Scoped Instances**: Register components as scoped, tying them to a specific context or scope. Scoped components are
  created once per scope and reused within that scope.


- **Lazy Initialization**: Support for lazy initialization of components, allowing for efficient resource utilization.


- **Multiple Implementations of the Same Interface**: Register and retrieve several implementations of the same
  interface type, allowing for flexible and modular design.


- **Keyed Services Injection**: Support for injecting services based on a key, allowing you to differentiate between
  multiple implementations of the same interface or type.


- **Concurrency-Safe**: Utilizes a mutex to ensure safe concurrent access to the container.

<br />

# Installation

```bash
go get -u github.com/firasdarwish/ore
```

<br />

# Usage

### Import

```go
import "github.com/firasdarwish/ore"
```

### Example Service

```go
// interface
type Counter interface {
  AddOne()
  GetCount() int
}

// implementation
type simpleCounter struct {
  counter int
}

func (c *simpleCounter) AddOne()  {
  c.counter++
}

func (c *simpleCounter) GetCount() int {
  return c.counter
}

func (c *simpleCounter) New(ctx context.Context) Counter {
  return &simpleCounter{}
}
```

<br />

### Eager Singleton

```go
package main

import (
  "context"
  "github.com/firasdarwish/ore"
)

func main() {
  var c Counter
  c = &simpleCounter{}

  // register
  ore.RegisterEagerSingleton[Counter](c)

  ctx := context.Background()

  // retrieve
  c, ctx := ore.Get[Counter](ctx)
  c.AddOne()
  c.AddOne()
}
```

<br />

### Lazy (using Creator[T] interface)

```go
package main

import (
  "context"
  "fmt"
  "github.com/firasdarwish/ore"
)

func main() {
  // register
  ore.RegisterLazyCreator[Counter](ore.Scoped, &simpleCounter{})

  // OR
  //ore.RegisterLazyCreator[Counter](ore.Transient, &simpleCounter{})
  //ore.RegisterLazyCreator[Counter](ore.Singleton, &simpleCounter{})

  ctx := context.Background()

  // retrieve
  c, ctx := ore.Get[Counter](ctx)
  c.AddOne()
  c.AddOne()

  // retrieve again
  c, ctx = ore.Get[Counter](ctx)
  c.AddOne()

  // prints out: `TOTAL: 3`
  fmt.Println("TOTAL: ", c.GetCount())
}
```

<br />

### Lazy (using anonymous func)

```go
package main

import (
  "context"
  "fmt"
  "github.com/firasdarwish/ore"
)

func main() {
  // register
  ore.RegisterLazyFunc[Counter](ore.Scoped, func(ctx context.Context) Counter {
    return &simpleCounter{}
  })

  // OR
  //ore.RegisterLazyFunc[Counter](ore.Transient, func(ctx context.Context) Counter {
  //  return &simpleCounter{}
  //})

  // Keyed service registration
  //ore.RegisterLazyFunc[Counter](ore.Singleton, func(ctx context.Context) Counter {
  // return &simpleCounter{}
  //}, "name here", 1234)

  ctx := context.Background()

  // retrieve
  c, ctx := ore.Get[Counter](ctx)
  c.AddOne()
  c.AddOne()

  // Keyed service retrieval
  //c, ctx := ore.Get[Counter](ctx, "name here", 1234)

  // retrieve again
  c, ctx = ore.Get[Counter](ctx)
  c.AddOne()

  // prints out: `TOTAL: 3`
  fmt.Println("TOTAL: ", c.GetCount())
}
```

<br />

### Several Implementations

```go
package main

import (
  "context"
  "github.com/firasdarwish/ore"
)

func main() {
  // register
  ore.RegisterLazyCreator[Counter](ore.Scoped, &simpleCounter{})

  ore.RegisterLazyCreator[Counter](ore.Scoped, &yetAnotherCounter{})

  ore.RegisterLazyFunc[Counter](ore.Transient, func(ctx context.Context) Counter {
    return &simpleCounter{}
  })

  ore.RegisterLazyCreator[Counter](ore.Singleton, &yetAnotherCounter{})

  ctx := context.Background()

  // returns a slice of `Counter` implementations
  counters, ctx := ore.GetList[Counter](ctx)

  // to retrieve a slice of keyed services
  //counters, ctx := ore.GetList[Counter](ctx, "my integer counters")

  for _, c := range counters {
    c.AddOne()
  }

  // It will always return the LAST registered implementation
  defaultImplementation, ctx := ore.Get[Counter](ctx) // simpleCounter
  defaultImplementation.AddOne()
}

```

<br />

### Keyed Services Retrieval Example

```go
package main

import (
  "context"
  "fmt"
  "github.com/firasdarwish/ore"
)

func main() {
  // register
  ore.RegisterLazyFunc[Counter](ore.Singleton, func(ctx context.Context) Counter {
    return &simpleCounter{}
  }, "name here", 1234)

  //ore.RegisterLazyCreator[Counter](ore.Scoped, &simpleCounter{}, "name here", 1234)

  //ore.RegisterEagerSingleton[Counter](&simpleCounter{}, "name here", 1234)

  ctx := context.Background()

  // Keyed service retrieval
  c, ctx := ore.Get[Counter](ctx, "name here", 1234)
  c.AddOne()

  // prints out: `TOTAL: 1`
  fmt.Println("TOTAL: ", c.GetCount())
}

```

## More Complex Example
```go

type Numeric interface {
	int
}

type GenericCounter[T Numeric] interface {
  Add(number T)
  GetCount() T
}

type genericCounter[T Numeric] struct {
  counter T
}

func (gc *genericCounter[T]) Add(number T) {
  gc.counter += number
}

func (gc *genericCounter[T]) GetCount(ctx context.Context) T {
  return gc.counter
}
```

```go
package main

import (
  "context"
  "github.com/firasdarwish/ore"
)

func main() {

  // register
  ore.RegisterLazyFunc[GenericCounter[int]](ore.Scoped, func(ctx context.Context) GenericCounter[int] {
    return &genericCounter[int]{}
  })

  // retrieve
  c, ctx := ore.Get[GenericCounter[int]](ctx)
}

```

<br />


# Contributing


Feel free to contribute by opening issues, suggesting features, or submitting pull requests. We welcome your feedback and contributions.


# License


This project is licensed under the MIT License - see the LICENSE file for details.
