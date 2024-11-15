# ore - Generic Dependency Injection Container for Go

[![Go Report Card](https://goreportcard.com/badge/github.com/firasdarwish/ore)](https://goreportcard.com/report/github.com/firasdarwish/ore)
[![Go Reference](https://pkg.go.dev/badge/github.com/firasdarwish/ore.svg)](https://pkg.go.dev/github.com/firasdarwish/ore)
[![Mentioned in Awesome Go](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go?tab=readme-ov-file#dependency-injection)
[![Maintainability](https://api.codeclimate.com/v1/badges/3bd6f2fa4390af7c8faa/maintainability)](https://codeclimate.com/github/firasdarwish/ore/maintainability)
[![codecov](https://codecov.io/gh/firasdarwish/ore/graph/badge.svg?token=ISZVCCYGCR)](https://codecov.io/gh/firasdarwish/ore)
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Ffirasdarwish%2Fore.svg?type=shield&issueType=license)](https://app.fossa.com/projects/git%2Bgithub.com%2Ffirasdarwish%2Fore?ref=badge_shield&issueType=license)
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Ffirasdarwish%2Fore.svg?type=shield&issueType=security)](https://app.fossa.com/projects/git%2Bgithub.com%2Ffirasdarwish%2Fore?ref=badge_shield&issueType=security)
![ore](https://github.com/firasdarwish/ore/assets/1930361/c1426ba1-777a-43f5-8a9a-7520caa45516)


___
`ore` is a lightweight, generic & simple dependency injection (DI) container for Go.

Inspired by the principles of ASP.NET Dependency Injection, designed to facilitate
the management of object lifetimes and the inversion of control in your applications.


<br />

## Features

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

## Installation

```bash
go get -u github.com/firasdarwish/ore
```

<br />

## Usage

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
  return &models.SimpleCounter{}, ctx
}
```

<br />

### Eager Singleton

```go
var c Counter
c = &models.SimpleCounter{}

// register
ore.RegisterEagerSingleton[Counter](c)

ctx := context.Background()

// retrieve
c, ctx := ore.Get[Counter](ctx)
c.AddOne()
c.AddOne()
```

<br />

### Lazy (using Creator[T] interface)

```go
// register
ore.RegisterLazyCreator[Counter](ore.Scoped, &models.SimpleCounter{})

// OR
//ore.RegisterLazyCreator[Counter](ore.Transient, &models.SimpleCounter{})
//ore.RegisterLazyCreator[Counter](ore.Singleton, &models.SimpleCounter{})

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
```

<br />

### Lazy (using anonymous func)

```go
  // register
ore.RegisterLazyFunc[Counter](ore.Scoped, func(ctx context.Context) (Counter, context.Context) {
    return &models.SimpleCounter{}, ctx
})

// OR
//ore.RegisterLazyFunc[Counter](ore.Transient, func(ctx context.Context) (Counter, context.Context) {
//  return &models.SimpleCounter{}, ctx
//})

// Keyed service registration
//ore.RegisterLazyFunc[Counter](ore.Singleton, func(ctx context.Context) (Counter, context.Context) {
// return &models.SimpleCounter{}, ctx
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
```

<br />

### Several Implementations

```go
  // register
ore.RegisterLazyCreator[Counter](ore.Scoped, &models.SimpleCounter{})

ore.RegisterLazyCreator[Counter](ore.Scoped, &yetAnotherCounter{})

ore.RegisterLazyFunc[Counter](ore.Transient, func(ctx context.Context) (Counter, context.Context) {
    return &models.SimpleCounter{}, ctx
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
```

#### Injecting Mocks in Tests

The last registered implementation takes precedence, so you can register a mock implementation in the test, which will override the real implementation.

<br />

### Keyed Services Retrieval Example

```go
  // register
ore.RegisterLazyFunc[Counter](ore.Singleton, func(ctx context.Context) (Counter, context.Context) {
    return &models.SimpleCounter{}, ctx
}, "name here", 1234)

//ore.RegisterLazyCreator[Counter](ore.Scoped, &models.SimpleCounter{}, "name here", 1234)

//ore.RegisterEagerSingleton[Counter](&models.SimpleCounter{}, "name here", 1234)

ctx := context.Background()

// Keyed service retrieval
c, ctx := ore.Get[Counter](ctx, "name here", 1234)
c.AddOne()

// prints out: `TOTAL: 1`
fmt.Println("TOTAL: ", c.GetCount())
```
<br />

### Alias: Register `struct`, Get `interface`

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

<br />

### Registration Validation

Once you're done with registering all the services, it is recommended to call `ore.Seal()`, then `ore.Validate()`, then finally `ore.DisableValidation=true`.

`ore.Validate()` invokes ALL your registered resolvers. The purpose is to panic early if your registrations were in bad shape:

- Missing Dependency: you forgot to register certain resolvers.
- Circular Dependency: A depends on B which depends on A.
- Lifetime Misalignment: a longer lifetime service (eg. Singleton) depends on a shorter one (eg Transient).

<br />

### Registration Recommendation

(1) You should call `ore.Validate()`

- In a test which is automatically run on your CI/CD pipeline (option 1)
- On application start, right after all the registrations (option 2)

Option 1 (run `ore.Validate` on test) is usually a better choice.

(2) It is recommended to seal your container `ore.Seal()` (which seals the container) on application start => Please don't call `ore.RegisterXX` all over the place.

(3) A combination of `ore.Buile()` and then `ore.Validate()` and then `ore.DisabledValidation=true` ensures no more new resolvers will be registered AND all registered resolvers are validated, this will 
prevent any further validation each time a resolver is invoked (`ore.Get`) which greatly enhances performance.

(4) Keep the object creation function (a.k.a resolvers) simple. Their only responsibility should be **object creation**.

- They should not spawn new goroutine
- They should not open database connection
- They should not contain any "if" statement or other business logic

<br />

### Graceful application termination

On application termination, you want to call `Shutdown()` on all the "Singletons" objects which have been created during the application lifetime.

Here how Ore can help you:

```go
// Assuming that the Application provides certain instances with Singleton lifetime.
// Some of these singletons implement a custom `Shutdowner` interface (defined within the application)
type Shutdowner interface {
  Shutdown()
}
ore.RegisterEagerSingleton(&Logger{}) //*Logger implements Shutdowner
ore.RegisterEagerSingleton(&SomeRepository{}) //*SomeRepository implements Shutdowner
ore.RegisterEagerSingleton(&SomeService{}, "some_module") //*SomeService implements Shutdowner

//On application termination, Ore can help to retrieve all the singletons implementation 
//of the `Shutdowner` interface.
//There might be other `Shutdowner`'s implementation which were lazily registered but 
//have never been created.
//Ore will ignore them, and return only the concrete instances which can be Shutdown()
shutdowables := ore.GetResolvedSingletons[Shutdowner]() 

//Now we can Shutdown() them all and gracefully terminate our application.
//The most recently invoked instance will be Shutdown() first
for _, instance := range disposables {
   instance.Shutdown()
}
```

In resume, the `ore.GetResolvedSingletons[TInterface]()` function returns a list of Singleton implementations of the `[TInterface]`.

- It returns only the instances which had been invoked (a.k.a resolved).
- All the implementations (including "keyed" ones) will be returned.
- The returned instances are sorted by the invocation order, the first one being latest invoked one.
  - if "A" depends on "B", "C", Ore will make sure to return "B" and "C" first in the list so that they would be shutdowned before "A".

<br />

### Graceful context termination

On context termination, you want to call `Dispose()` on all the "Scoped" objects which have been created during the context lifetime.

Here how Ore can help you:

```go
//Assuming that your Application provides certain instances with Scoped lifetime.
//Some of them implements a "Disposer" interface (defined within the application).
type Disposer interface {
  Dispose()
}
ore.RegisterLazyCreator(ore.Scoped, &SomeDisposableService{}) //*SomeDisposableService implements Disposer

//a new request arrive
ctx, cancel := context.WithCancel(context.Background())

//start a go routine that will clean up resources when the context is canceled
go func() {
  <-ctx.Done() // Wait for the context to be canceled
  // Perform your cleanup tasks here
  disposables := ore.GetResolvedScopedInstances[Disposer](ctx)
  //The most recently invoked instance will be Dispose() first
  for _, d := range disposables {
    _ = d.Dispose(ctx)
  }
}()
...
ore.Get[*SomeDisposableService](ctx) //invoke some scoped services
cancel() //cancel the ctx

```

The `ore.GetResolvedScopedInstances[TInterface](context)` function returns a list of implementations of the `[TInterface]` which are Scoped in the input context:

- It returns only the instances which had been invoked (a.k.a resolved) during the context lifetime.
- All the implementations (of all modules) including "keyed" one will be returned.
- The returned instances are sorted by invocation order, the first one being the latest invoked one.
  - if "A" depends on "B", "C", Ore will make sure to return "B" and "C" first in the list so that they would be Disposed before "A".

<br />

### Multiple Containers (a.k.a Modules)

| DefaultContainer | Custom container |
|------------------|------------------|
| Get | GetFromContainer |
| GetList | GetListFromContainer |
| GetResolvedSingletons | GetResolvedSingletonsFromContainer |
| RegisterAlias | RegisterAliasToContainer |
| RegisterEagerSingleton | RegisterEagerSingletonToContainer |
| RegisterLazyCreator | RegisterLazyCreatorToContainer |
| RegisterLazyFunc | RegisterLazyFuncToContainer |
| RegisterPlaceHolder | RegisterPlaceHolderToContainer |
| ProvideScopedValue | ProvideScopedValueToContainer |

Most of time you only need the Default Container. In rare use case such as the Modular Monolith Architecture, you might want to use multiple containers, one per module. Ore provides minimum support for "module" in this case:

```go
//broker module
brokerContainer := ore.NewContainer()
ore.RegisterLazyFuncToContainer(brokerContainer, ore.Singleton, func(ctx context.Context) (*Broker, context.Context) {
  brs, ctx = ore.GetFromContainer[*BrokerageSystem](brokerContainer, ctx)
  return &Broker{brs}, ctx
})
// brokerContainer.Seal() //prevent further registration
// brokerContainer.Validate() //check the dependency graph
// brokerContainer.DisableValidation = true //disable check when resolve new object
broker, _ := ore.GetFromContainer[*Broker](brokerContainer, context.Background())

//trader module
traderContainer := ore.NewContainer()
ore.RegisterLazyFuncToContainer(traderContainer, ore.Singleton, func(ctx context.Context) (*Trader, context.Context) {
  mkp, ctx = ore.GetFromContainer[*MarketPlace](traderContainer, ctx)
  return &Trader{mkp}, ctx
})
trader, _ := ore.GetFromContainer[*Trader](traderContainer, context.Background())
```

Important: You will have to prevent cross modules access to the containers by yourself. For eg, don't let your "Broker
module" to have access to the `traderContainer` of the "Trader module".

<br />

### Injecting value at Runtime

A common scenario is that your "Service" depends on something which you couldn't provide on registration time. You can provide this dependency only when certain requests or events arrive later. Ore allows you to build an "incomplete" dependency graph using the "placeholder".

```go
//register SomeService which depends on "someConfig"
ore.RegisterLazyFunc[*SomeService](ore.Scoped, func(ctx context.Context) (*SomeService, context.Context) {
  someConfig, ctx := ore.Get[string](ctx, "someConfig")
  return &SomeService{someConfig}, ctx
})

//someConfig is unknow at registration time because 
//this value depends on the future user's request
ore.RegisterPlaceHolder[string]("someConfig")

//a new request arrive
ctx := context.Background()
//suppose that the request is sent by "admin"
ctx = context.WithValue(ctx, "role", "admin")

//inject a different somConfig value depending on the request's content
userRole := ctx.Value("role").(string)
if userRole == "admin" {
  ctx = ore.ProvideScopedValue(ctx, "Admin config", "someConfig")
} else if userRole == "supervisor" {
  ctx = ore.ProvideScopedValue(ctx, "Supervisor config", "someConfig")
} else if userRole == "user" {
  if (isAuthenticatedUser) {
    ctx = ore.ProvideScopedValue(ctx, "Public user config", "someConfig")
  } else {
    ctx = ore.ProvideScopedValue(ctx, "Private user config", "someConfig")
  }
}

//Get the service to handle this request
service, ctx := ore.Get[*SomeService](ctx)
fmt.Println(service.someConfig) //"Admin config"
```

([See full codes here](./examples/placeholderdemo/main.go))

- `ore.RegisterPlaceHolder[T](key...)` registers a future value with Scoped lifetime.
  - This value will be injected in runtime using the `ProvideScopedValue` function.
  - Resolving objects which depend on this value will panic if the value has not been provided.

- `ore.ProvideScopedValue[T](context, value T, key...)` injects a concrete value into the given context
  - `ore` can access (`Get()` or `GetList()`) to this value only if the corresponding placeholder (which matches the type and keys) is registered.

- A value provided to a placeholder would never replace value returned by other resolvers. It's the opposite, if a type (and key) could be resolved by a real resolver (such as `RegisterLazyFunc`, `RegisterLazyCreator`...), then the later would take precedent.

<br/>

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

// register
ore.RegisterLazyFunc[GenericCounter[int]](ore.Scoped, func(ctx context.Context) (GenericCounter[int], context.Context) {
    return &genericCounter[int]{}, ctx
})

// retrieve
c, ctx := ore.Get[GenericCounter[int]](ctx)
```

<br />

## Benchmarks

```bash
goos: windows
goarch: amd64
pkg: github.com/firasdarwish/ore
cpu: 13th Gen Intel(R) Core(TM) i9-13900H
BenchmarkRegisterLazyFunc-20             5706694               196.9 ns/op
BenchmarkRegisterLazyCreator-20          6283534               184.5 ns/op
BenchmarkRegisterEagerSingleton-20       5146953               211.5 ns/op
BenchmarkInitialGet-20                   3440072               352.1 ns/op
BenchmarkGet-20                          9806043               121.8 ns/op
BenchmarkInitialGetList-20               1601787               747.9 ns/op
BenchmarkGetList-20                      4237449               282.1 ns/op
PASS
ok      github.com/firasdarwish/ore     11.427s
```

Checkout also [examples/benchperf/README.md](examples/benchperf/README.md)

<br />

## 👤 Contributors

![Contributors](https://contrib.rocks/image?repo=firasdarwish/ore)


## Contributing

Feel free to contribute by opening issues, suggesting features, or submitting pull requests. We welcome your feedback
and contributions.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
