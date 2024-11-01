# ore - Generic Dependency Injection Container for Go

[![Go Report Card](https://goreportcard.com/badge/github.com/firasdarwish/ore)](https://goreportcard.com/report/github.com/firasdarwish/ore)
[![Go Reference](https://pkg.go.dev/badge/github.com/firasdarwish/ore.svg)](https://pkg.go.dev/github.com/firasdarwish/ore)
[![Mentioned in Awesome Go](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go?tab=readme-ov-file#dependency-injection)
[![Maintainability](https://api.codeclimate.com/v1/badges/3bd6f2fa4390af7c8faa/maintainability)](https://codeclimate.com/github/firasdarwish/ore/maintainability)
[![codecov](https://codecov.io/gh/firasdarwish/ore/graph/badge.svg?token=ISZVCCYGCR)](https://codecov.io/gh/firasdarwish/ore)

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

func (c *simpleCounter) New(ctx context.Context) (Counter, context.Context) {
  return &simpleCounter{}, ctx
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
  ore.RegisterLazyFunc[Counter](ore.Scoped, func(ctx context.Context) (Counter, context.Context) {
    return &simpleCounter{}, ctx
  })

  // OR
  //ore.RegisterLazyFunc[Counter](ore.Transient, func(ctx context.Context) (Counter, context.Context) {
  //  return &simpleCounter{}, ctx
  //})

  // Keyed service registration
  //ore.RegisterLazyFunc[Counter](ore.Singleton, func(ctx context.Context) (Counter, context.Context) {
  // return &simpleCounter{}, ctx
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

  ore.RegisterLazyFunc[Counter](ore.Transient, func(ctx context.Context) (Counter, context.Context) {
    return &simpleCounter{}, ctx
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

#### Injecting Mocks in Tests

The last registered implementation takes precedence, so you can register a mock implementation in the test, which will override the real implementation.

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
  ore.RegisterLazyFunc[Counter](ore.Singleton, func(ctx context.Context) (Counter, context.Context) {
    return &simpleCounter{}, ctx
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

### Alias: Register struct, get interface

```go
type IPerson interface{}
type Broker struct {
  Name string
} //implements IPerson

type Trader struct {
  Name string
} //implements IPerson

func TestGetInterfaceAlias(t *testing.T) {
  ore.RegisterLazyFunc(ore.Scoped, func(ctx context.Context) (*Broker, context.Context) {
    return &Broker{Name: "Peter"}, ctx
  })
  ore.RegisterLazyFunc(ore.Scoped, func(ctx context.Context) (*Broker, context.Context) {
    return &Broker{Name: "John"}, ctx
  })
  ore.RegisterLazyFunc(ore.Scoped, func(ctx context.Context) (*Trader, context.Context) {
    return &Trader{Name: "Mary"}, ctx
  })

  ore.RegisterAlias[IPerson, *Trader]() //link IPerson to *Trader
  ore.RegisterAlias[IPerson, *Broker]() //link IPerson to *Broker

  //no IPerson was registered to the container, but we can still `Get` it out of the container.
  //(1) IPerson is alias to both *Broker and *Trader. *Broker takes precedence because it's the last one linked to IPerson.
  //(2) multiple *Borker (Peter and John) are registered to the container, the last registered (John) takes precedence.
  person, _ := ore.Get[IPerson](context.Background()) // will return the broker John

  personList, _ := ore.GetList[IPerson](context.Background()) // will return all registered broker and trader
}
```

Alias is also scoped by key. When you "Get" an alias with keys for eg: `ore.Get[IPerson](ctx, "module1")` then Ore would return only Services registered under this key ("module1") and panic if no service found.

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
  ore.RegisterLazyFunc[GenericCounter[int]](ore.Scoped, func(ctx context.Context) (GenericCounter[int], context.Context) {
    return &genericCounter[int]{}, ctx
  })

  // retrieve
  c, ctx := ore.Get[GenericCounter[int]](ctx)
}

```

<br />

# Benchmarks

```bash
goos: windows
goarch: amd64
pkg: github.com/firasdarwish/ore
cpu: 13th Gen Intel(R) Core(TM) i9-13900H
BenchmarkRegisterLazyFunc
BenchmarkRegisterLazyFunc-20             4953412               233.5 ns/op
BenchmarkRegisterLazyCreator
BenchmarkRegisterLazyCreator-20          5468863               231.3 ns/op
BenchmarkRegisterEagerSingleton
BenchmarkRegisterEagerSingleton-20       4634733               267.4 ns/op
BenchmarkGet
BenchmarkGet-20                          3766730               321.9 ns/op
BenchmarkGetList
BenchmarkGetList-20                      1852132               637.0 ns/op
```

# Contributing

Feel free to contribute by opening issues, suggesting features, or submitting pull requests. We welcome your feedback
and contributions.

# License

This project is licensed under the MIT License - see the LICENSE file for details.
