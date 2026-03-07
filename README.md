# Ore — Dependency Injection for Go

[![Go Reference](https://pkg.go.dev/badge/github.com/firasdarwish/ore.svg)](https://pkg.go.dev/github.com/firasdarwish/ore)
[![Go Report Card](https://goreportcard.com/badge/github.com/firasdarwish/ore)](https://goreportcard.com/report/github.com/firasdarwish/ore)
[![Mentioned in Awesome Go](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go)
[![codecov](https://codecov.io/gh/firasdarwish/ore/graph/badge.svg?token=ISZVCCYGCR)](https://codecov.io/gh/firasdarwish/ore)

**Ore** is a lightweight, type-safe dependency injection (DI) container for Go. Inspired by ASP.NET's DI model, it gives you clean lifetime management, lazy initialization, runtime value injection, and modular containers — without magic or reflection soup.

---

## Table of Contents

1. [Why Ore?](#why-ore)
2. [Installation](#installation)
3. [Core Concepts](#core-concepts)
4. [Lifetimes](#lifetimes)
  - [Singleton](#singleton)
  - [Scoped](#scoped)
  - [Transient](#transient)
5. [Registering Services](#registering-services)
  - [Eager Singleton](#eager-singleton)
  - [Anonymous Functions](#anonymous-functions-registerfunc)
  - [Creator\[T\] Interface](#creatort-interface-registercreator)
6. [Resolving Services](#resolving-services)
  - [Get](#get)
  - [GetList](#getlist)
7. [Keyed Services](#keyed-services)
8. [Aliases](#aliases)
9. [Placeholder Services](#placeholder-services)
10. [Isolated Containers](#isolated-containers)
11. [Validation](#validation)
12. [Graceful Termination](#graceful-termination)
13. [Recommended Startup Pattern](#recommended-startup-pattern)
14. [Real-World Usage Patterns](#real-world-usage-patterns)
15. [API Reference](#api-reference)

---

## Why Ore?

Go encourages simplicity, and many projects wire dependencies by hand. That works — until it doesn't. As applications grow, manual wiring becomes error-prone, lifetime bugs sneak in, and test setup becomes tedious.

Ore solves this without going overboard:

- **Type-safe** — powered by Go generics, no `interface{}` casting
- **Lifetime-aware** — Singleton, Scoped, and Transient, correctly enforced
- **Context-native** — scoped instances live and die with `context.Context`
- **Validates your graph** — catches circular dependencies, missing registrations, and lifetime mismatches at startup
- **Modular** — isolated containers per module for clean monolith architecture
- **Non-invasive** — your structs don't need to embed anything from Ore

---

## Installation

```bash
go get -u github.com/firasdarwish/ore
```

```go
import "github.com/firasdarwish/ore"
```

---

## Core Concepts

Ore works in two phases:

**1. Registration** — at startup, you tell Ore how to build each service and what lifetime it should have.

**2. Resolution** — at runtime, you ask Ore for a service. It constructs it (or returns a cached instance, depending on lifetime), injects its dependencies, and returns it to you.

The key insight: **you never call `new(MyService)` scattered throughout your app**. Ore centralizes construction, so lifetimes are predictable and dependencies are explicit.

---

## Lifetimes

Every service registered with Ore has a **lifetime** that controls when instances are created and how long they live.

### Singleton

A singleton is created **once** and reused for the entire application lifetime. Use this for stateless services, configuration, database connection pools, loggers, and anything that's safe to share globally.

```
App starts → instance created → reused forever → app shuts down
```

There are two flavors:

**Eager singleton** — created immediately at registration time. Best for critical services where you want startup failures to surface early.

**Lazy singleton** — created the first time it's resolved. Best for services that may never be used, or that are expensive to initialize.

### Scoped

A scoped service is created **once per context** (`context.Context`). Every call to `Get` with the same context returns the same instance. A new context gets a fresh instance.

```
Request A context → instance A (shared within request A)
Request B context → instance B (shared within request B)
```

This is ideal for HTTP request handlers, database transactions, and anything that should be consistent within a single unit of work but isolated from other units.

### Transient

A transient service is created **fresh on every resolution**. No caching, no sharing.

```
Get() → new instance
Get() → another new instance
```

Use transients for lightweight, stateful objects where sharing would cause bugs.

---

## Registering Services

### Eager Singleton

Pass an already-constructed instance directly. Ore stores it and returns it on every `Get`.

```go
type Logger interface {
    Log(msg string)
}

type zapLogger struct{}
func (z *zapLogger) Log(msg string) { /* ... */ }

// Registered immediately — no construction function needed
ore.RegisterSingleton[Logger](&zapLogger{})
```

### Anonymous Functions (`RegisterFunc`)

Pass a constructor function. Ore calls it when the service is first needed (or on every `Get` for transients). The function receives a `context.Context` and returns the service plus the (potentially updated) context.

```go
ore.RegisterFunc[Logger](ore.Singleton, func(ctx context.Context) (Logger, context.Context) {
    return &zapLogger{}, ctx
})
```

**Injecting dependencies inside a constructor:**

When one service depends on another, call `ore.Get` inside the constructor and pass the context through. This is what threads scoped state correctly across a dependency chain.

```go
type UserService interface {
    GetUser(id string) (*User, error)
}

type userServiceImpl struct {
    db DB
}

ore.RegisterFunc[UserService](ore.Scoped, func(ctx context.Context) (UserService, context.Context) {
    // Resolve DB — passes ctx through so scoped instances are shared
    db, ctx := ore.Get[DB](ctx)
    return &userServiceImpl{db: db}, ctx
})
```

> **Always pass `ctx` through the dependency chain.** Scoped services store their instances in the context. If you don't thread the returned `ctx` forward, scoped dependencies won't be shared correctly within the same scope.

### `Creator[T]` Interface (`RegisterCreator`)

An alternative to anonymous functions. Implement the `Creator[T]` interface on your struct, and Ore calls `New` to construct it.

```go
type Creator[T any] interface {
    New(ctx context.Context) (T, context.Context)
}
```

This is useful when you want the struct itself to own its construction logic, keeping it colocated with the type definition.

```go
type simpleCounter struct {
    count int
}

// Implementing Creator[Counter]
func (c *simpleCounter) New(ctx context.Context) (Counter, context.Context) {
    return &simpleCounter{count: 0}, ctx
}

func (c *simpleCounter) AddOne()       { c.count++ }
func (c *simpleCounter) GetCount() int { return c.count }

// Register using RegisterCreator
ore.RegisterCreator[Counter](ore.Scoped, &simpleCounter{})
```

**Anonymous func vs Creator[T] — when to use which:**

| | `RegisterFunc` | `RegisterCreator` |
|---|---|---|
| Constructor location | Inline at registration site | On the struct itself |
| Good for | External types, simple wiring | Types that own their construction |
| Coupling to Ore | None | Struct knows about `context.Context` |
| Verbosity | Slightly more boilerplate | Cleaner registration call |

---

## Resolving Services

### `Get`

Resolves a single service. Returns the service and an updated context (important for scoped services).

```go
ctx := context.Background()

logger, ctx := ore.Get[Logger](ctx)
logger.Log("hello")
```

Always use the returned `ctx` for subsequent resolutions if any services in the chain are scoped.

### `GetList`

Resolves **all registered implementations** of a type. Returns a slice.

```go
// If you registered three different Counter implementations:
counters, ctx := ore.GetList[Counter](ctx)
for _, c := range counters {
    fmt.Println(c.GetCount())
}
```

`GetList` is useful for plugin-style architectures or when you genuinely need all implementations (e.g., firing all event handlers, running all validators).

`GetList` never panics if nothing is registered — it returns an empty slice.

---

## Keyed Services

Keys let you register **multiple implementations of the same type** and select among them by name at resolution time.

```go
type Greeter interface {
    Greet() string
}

type FriendlyGreeter struct{}
func (g *FriendlyGreeter) Greet() string { return "Hey! Great to see you!" }

type FormalGreeter struct{}
func (g *FormalGreeter) Greet() string { return "Good day. How may I assist you?" }

// Register with keys
ore.RegisterKeyedFunc[Greeter](ore.Scoped, func(ctx context.Context) (Greeter, context.Context) {
    return &FriendlyGreeter{}, ctx
}, "friendly")

ore.RegisterKeyedFunc[Greeter](ore.Transient, func(ctx context.Context) (Greeter, context.Context) {
    return &FormalGreeter{}, ctx
}, "formal")

// Resolve by key
friendly, ctx := ore.GetKeyed[Greeter](ctx, "friendly")
fmt.Println(friendly.Greet()) // Hey! Great to see you!

formal, ctx := ore.GetKeyed[Greeter](ctx, "formal")
fmt.Println(formal.Greet()) // Good day. How may I assist you?
```

**Get all implementations under a key:**

```go
greeters, ctx := ore.GetKeyedList[Greeter](ctx, "friendly")
```

**Common use cases for keyed services:**

- Multiple payment providers (`"stripe"`, `"paypal"`)
- Multiple notification channels (`"email"`, `"sms"`, `"push"`)
- Module-specific service overrides
- Feature flags driving different implementations at runtime

---

## Aliases

Aliases let you link a **concrete type** (struct pointer) to an **interface**, without registering the interface directly. This is powerful when you want to resolve by a broad interface but your implementations are registered under their concrete types.

```go
type IPerson interface {
    GetName() string
}

type Broker struct{ Name string }
func (b *Broker) GetName() string { return b.Name }

type Trader struct{ Name string }
func (t *Trader) GetName() string { return t.Name }

// Register concrete types
ore.RegisterFunc[*Broker](ore.Scoped, func(ctx context.Context) (*Broker, context.Context) {
    return &Broker{Name: "Alice"}, ctx
})
ore.RegisterFunc[*Trader](ore.Scoped, func(ctx context.Context) (*Trader, context.Context) {
    return &Trader{Name: "Bob"}, ctx
})

// Link them to IPerson
ore.RegisterAlias[IPerson, *Broker]()
ore.RegisterAlias[IPerson, *Trader]() // last-linked takes precedence for Get

// Resolve — returns Trader (last linked)
person, ctx := ore.Get[IPerson](ctx)
fmt.Println(person.GetName()) // Bob

// Resolve all — returns both
people, ctx := ore.GetList[IPerson](ctx)
fmt.Println(len(people)) // 2
```

**Precedence rules:**

- `Get` returns the **most recently linked** alias, unless a direct resolver for the interface exists — then the direct resolver always wins.
- `GetList` returns **all** linked implementations plus any direct resolvers.
- Alias-of-alias is not supported and will panic.
- Ore validates at registration that the concrete type actually implements the interface.

---

## Placeholder Services

Placeholders let you declare a dependency that **doesn't exist at registration time** but will be provided at runtime. This is the right pattern for request-scoped values like authenticated users, tenant IDs, or feature flag evaluations.

### How it works

1. **Declare** a placeholder at startup — this tells Ore "something of this type will be provided later."
2. **Provide** the value at request time by injecting it into the context.
3. **Resolve** it normally — any service that depends on the placeholder resolves correctly once the value is set.

```go
// 1. At startup: declare that a *User will be injected per-request
ore.RegisterPlaceholder[*User]()

// 2. Register a service that depends on *User
ore.RegisterFunc[UserDashboard](ore.Scoped, func(ctx context.Context) (UserDashboard, context.Context) {
    user, ctx := ore.Get[*User](ctx) // depends on the placeholder
    return &userDashboardImpl{user: user}, ctx
})

// 3. At request time: inject the actual user
func handleRequest(w http.ResponseWriter, r *http.Request) {
    user := getUserFromJWT(r)
    ctx := ore.ProvideScopedValue[*User](r.Context(), user)

    dashboard, ctx := ore.Get[UserDashboard](ctx)
    // dashboard.user is set correctly
}
```

**What happens if you forget to provide the value?** Ore panics when trying to resolve a service that depends on an unfulfilled placeholder. This is intentional — a silent nil would be far worse.

**Keyed placeholders** follow the same pattern using `RegisterKeyedPlaceholder` and `ProvideKeyedScopedValue`.

**Placeholders and real resolvers coexist:**

If you later register a real resolver for the same type+key, it takes precedence over the placeholder for `Get`. Both appear in `GetList`. This lets you gradually migrate from runtime injection to proper construction logic.

---

## Isolated Containers

By default, all registrations go into Ore's **default container**. For larger applications, you can create **isolated containers** — each with its own independent dependency graph.

This is the foundation for **modular monolith architecture**: each module owns its container, preventing accidental cross-module coupling.

```go
// Create isolated containers per module
brokerContainer := ore.NewContainer()
traderContainer := ore.NewContainer()

// Register into specific containers
ore.RegisterFuncToContainer(brokerContainer, ore.Singleton, func(ctx context.Context) (BrokerService, context.Context) {
    return &brokerServiceImpl{}, ctx
})

ore.RegisterFuncToContainer(traderContainer, ore.Scoped, func(ctx context.Context) (TraderService, context.Context) {
    return &traderServiceImpl{}, ctx
})

// Resolve from specific containers
broker, ctx := ore.GetFromContainer[BrokerService](brokerContainer, ctx)
trader, ctx := ore.GetFromContainer[TraderService](traderContainer, ctx)
```

Services registered in one container are completely invisible to another. There is no accidental leakage between modules.

**Using containers in tests:**

Isolated containers are excellent for unit tests. Each test gets its own clean container with mocked dependencies — no shared global state, no test order dependencies.

```go
func TestUserService(t *testing.T) {
    c := ore.NewContainer()
    ore.RegisterSingletonToContainer[DB](c, &mockDB{})
    ore.RegisterFuncToContainer(c, ore.Scoped, func(ctx context.Context) (UserService, context.Context) {
        db, ctx := ore.GetFromContainer[DB](c, ctx)
        return &userServiceImpl{db: db}, ctx
    })

    svc, _ := ore.GetFromContainer[UserService](c, context.Background())
    // test svc with mockDB
}
```

All container-scoped registration functions follow the naming convention `XxxToContainer` (e.g., `RegisterFuncToContainer`, `RegisterSingletonToContainer`, `RegisterPlaceholderToContainer`).

---

## Validation

Ore validates your dependency graph to catch bugs at startup rather than in production.

### What Ore validates

**Missing dependencies** — a service depends on a type that was never registered.

**Circular dependencies** — Service A depends on B, B depends on A. Would cause infinite recursion.

**Lifetime misalignment** — a long-lived service depends on a shorter-lived one. For example, a Singleton depending on a Scoped service is a bug: the Singleton gets created once, captures a scoped instance, and that instance outlives its intended scope. Ore catches this.

### How to use validation

By default, Ore validates on every `Get` call. For production, this per-call overhead can be eliminated by disabling it after a one-time startup check:

```go
func main() {
    // Register all services...
    ore.RegisterSingleton[DB](&dbImpl{})
    ore.RegisterFunc[UserService](ore.Scoped, NewUserService)
    // ...

    // Lock the container — no new registrations allowed after this
    ore.Seal()

    // Validate the full graph once at startup
    // This resolves everything, checks all dependencies, then clears instances
    ore.Validate()

    // Disable per-call validation in production for best performance
    ore.DisableValidation = true

    // Start your server...
}
```

`ore.Seal()` causes Ore to panic if any code tries to register a new service after the fact — useful for preventing accidental late registrations in large codebases.

`ore.Validate()` tries to resolve all registered services (including their full dependency chains), verifies correctness, then clears the instances so the app starts fresh.

> **Constructor purity matters.** Since `Validate()` actually runs your constructors, they should be deterministic and side-effect-free. Don't make network calls, open files, or start goroutines inside constructors.

---

## Graceful Termination

When your application shuts down, resources held by Singleton services (DB connections, open files, background workers) need to be cleaned up. Ore helps coordinate this without prescribing a specific interface.

### Application shutdown (Singletons)

Define whatever cleanup interface fits your app:

```go
type Shutdowner interface {
    Shutdown() error
}

type Closer interface {
    Close() error
}
```

At shutdown, ask Ore for all resolved Singletons that implement your interface. Ore returns them in **reverse resolution order** — dependencies are cleaned up before their dependents.

```go
// In your shutdown handler
shutdownables := ore.GetResolvedSingletons[Shutdowner]()
for _, s := range shutdownables {
    if err := s.Shutdown(); err != nil {
        log.Printf("shutdown error: %v", err)
    }
}
```

Only Singletons that were **actually resolved** during the app's lifetime are returned. Lazily registered singletons that were never used are excluded.

### Request/context shutdown (Scoped)

For scoped services, cleanup happens when the context ends:

```go
type Disposer interface {
    Dispose(ctx context.Context) error
}

// In your request handler or middleware
ctx, cancel := context.WithCancel(r.Context())
defer func() {
    cancel()
    // Clean up scoped instances that implemented Disposer
    disposables := ore.GetResolvedScopedInstances[Disposer](ctx)
    for _, d := range disposables {
        _ = d.Dispose(ctx)
    }
}()

// Handle request using ctx...
```

Full example showing both patterns together:

```go
func main() {
    ore.RegisterSingleton[*GlobalRepo](&GlobalRepo{})
    ore.RegisterCreator(ore.Scoped, &ScopedRepo{})

    ore.Seal()
    ore.Validate()

    // Simulate a request
    ctx, cancel := context.WithCancel(context.Background())
    _, ctx = ore.Get[*ScopedRepo](ctx)

    // End request — clean up scoped resources
    cancel()
    disposables := ore.GetResolvedScopedInstances[Disposer](ctx)
    for _, d := range disposables {
        _ = d.Dispose(ctx)
    }

    // Shut down app — clean up singletons
    shutdownables := ore.GetResolvedSingletons[Shutdowner]()
    for _, s := range shutdownables {
        _ = s.Shutdown()
    }
}
```

---

## Recommended Startup Pattern

Here is the battle-tested pattern for setting up Ore in a production Go application:

```go
package main

import (
    "context"
    "github.com/firasdarwish/ore"
)

func registerServices() {
    // Eager singletons for critical infrastructure
    ore.RegisterSingleton[Config](loadConfig())

    // Lazy singletons for app-wide services
    ore.RegisterFunc[DB](ore.Singleton, func(ctx context.Context) (DB, context.Context) {
        cfg, ctx := ore.Get[Config](ctx)
        return connectDB(cfg.DatabaseURL), ctx
    })

    // Scoped services for per-request work
    ore.RegisterFunc[UserRepository](ore.Scoped, func(ctx context.Context) (UserRepository, context.Context) {
        db, ctx := ore.Get[DB](ctx)
        return &userRepo{db: db}, ctx
    })

    ore.RegisterFunc[UserService](ore.Scoped, func(ctx context.Context) (UserService, context.Context) {
        repo, ctx := ore.Get[UserRepository](ctx)
        return &userService{repo: repo}, ctx
    })

    // Placeholders for runtime-injected values
    ore.RegisterPlaceholder[*AuthUser]()
}

func main() {
    registerServices()

    // Lock and validate at startup
    ore.Seal()
    ore.Validate()
    ore.DisableValidation = true // skip per-call overhead in prod

    startServer()
}
```

---

## Real-World Usage Patterns

### HTTP middleware injecting auth user

```go
func authMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        user, err := parseJWT(r.Header.Get("Authorization"))
        if err != nil {
            http.Error(w, "unauthorized", 401)
            return
        }
        // Inject the authenticated user into the request context
        ctx := ore.ProvideScopedValue[*AuthUser](r.Context(), user)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

func handleGetProfile(w http.ResponseWriter, r *http.Request) {
    svc, _ := ore.Get[UserService](r.Context())
    // UserService was built with the *AuthUser from this request's context
    profile := svc.GetCurrentUserProfile()
    json.NewEncoder(w).Encode(profile)
}
```

### Modular monolith with isolated containers

```go
// broker/module.go
var BrokerContainer = ore.NewContainer()

func init() {
    ore.RegisterFuncToContainer(BrokerContainer, ore.Singleton, NewBrokerConfig)
    ore.RegisterFuncToContainer(BrokerContainer, ore.Scoped, NewOrderService)
    ore.RegisterFuncToContainer(BrokerContainer, ore.Scoped, NewPositionService)
}

// trader/module.go
var TraderContainer = ore.NewContainer()

func init() {
    ore.RegisterFuncToContainer(TraderContainer, ore.Singleton, NewTraderConfig)
    ore.RegisterFuncToContainer(TraderContainer, ore.Scoped, NewTradeExecutor)
}

// main.go
func main() {
    BrokerContainer.Seal()
    TraderContainer.Seal()
    BrokerContainer.Validate()
    TraderContainer.Validate()

    // Each module resolves from its own container
    // No accidental cross-module dependency is possible
}
```

### Plugin-style multi-implementation with GetList

```go
type Validator interface {
    Validate(input Input) error
}

// Register multiple validators
ore.RegisterFunc[Validator](ore.Singleton, NewEmailValidator)
ore.RegisterFunc[Validator](ore.Singleton, NewPhoneValidator)
ore.RegisterFunc[Validator](ore.Singleton, NewAgeValidator)

// Run all validators
validators, _ := ore.GetList[Validator](context.Background())
for _, v := range validators {
    if err := v.Validate(input); err != nil {
        return err
    }
}
```

### Switching implementations with keyed services

```go
type PaymentProvider interface {
    Charge(amount float64, card string) error
}

ore.RegisterKeyedFunc[PaymentProvider](ore.Singleton, NewStripeProvider, "stripe")
ore.RegisterKeyedFunc[PaymentProvider](ore.Singleton, NewPayPalProvider, "paypal")

// Choose at runtime based on user preference
func processPayment(ctx context.Context, method string, amount float64, card string) error {
    provider, ctx := ore.GetKeyed[PaymentProvider](ctx, method)
    return provider.Charge(amount, card)
}
```

---

## API Reference

### Registration

| Function | Description |
|---|---|
| `RegisterSingleton[T](impl T)` | Eager singleton — instance provided directly |
| `RegisterFunc[T](lifetime, fn)` | Lazy registration via anonymous constructor function |
| `RegisterCreator[T](lifetime, creator)` | Lazy registration via `Creator[T]` interface |
| `RegisterPlaceholder[T]()` | Declare a future runtime-injected value |
| `RegisterAlias[TInterface, TConcrete]()` | Link a concrete type to an interface |
| `RegisterKeyedFunc[T](lifetime, fn, key)` | Keyed variant of `RegisterFunc` |
| `RegisterKeyedSingleton[T](impl, key)` | Keyed eager singleton |
| `RegisterKeyedCreator[T](lifetime, creator, key)` | Keyed variant of `RegisterCreator` |
| `RegisterKeyedPlaceholder[T](key)` | Keyed placeholder |

All registration functions have a `ToContainer` variant (e.g., `RegisterFuncToContainer`) for isolated containers.

### Resolution

| Function | Description |
|---|---|
| `Get[T](ctx)` | Resolve a single service |
| `GetList[T](ctx)` | Resolve all registered implementations of T |
| `GetKeyed[T](ctx, key)` | Resolve a single keyed service |
| `GetKeyedList[T](ctx, key)` | Resolve all keyed implementations |
| `GetFromContainer[T](container, ctx)` | Resolve from a specific container |
| `GetListFromContainer[T](container, ctx)` | Resolve all from a specific container |

### Runtime Injection

| Function | Description |
|---|---|
| `ProvideScopedValue[T](ctx, value)` | Inject a value into a placeholder via context |
| `ProvideKeyedScopedValue[T](ctx, value, key)` | Keyed variant of `ProvideScopedValue` |
| `ProvideScopedValueToContainer(container, ctx, value)` | Inject into a specific container |

### Lifecycle

| Function | Description |
|---|---|
| `Seal()` | Lock the default container — no further registrations |
| `Validate()` | Validate the full dependency graph of the default container |
| `GetResolvedSingletons[T]()` | Get all resolved singletons implementing T (for shutdown) |
| `GetResolvedScopedInstances[T](ctx)` | Get all resolved scoped instances implementing T (for disposal) |
| `DisableValidation = true` | Disable per-call validation (use after startup `Validate()`) |

### Container

| Function | Description |
|---|---|
| `NewContainer()` | Create a new isolated container |
| `container.Seal()` | Lock an isolated container |
| `container.Validate()` | Validate an isolated container's dependency graph |
| `container.DisableValidation` | Per-container validation toggle |

---

> Full documentation and examples: **[ore.lilury.com](https://ore.lilury.com)**
>
> GitHub: **[github.com/firasdarwish/ore](https://github.com/firasdarwish/ore)**