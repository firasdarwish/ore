# Ore: A Lightweight Dependency Injection Container for Go

[![Go Reference](https://pkg.go.dev/badge/github.com/firasdarwish/ore.svg)](https://pkg.go.dev/github.com/firasdarwish/ore)
[![Go Report Card](https://goreportcard.com/badge/github.com/firasdarwish/ore)](https://goreportcard.com/report/github.com/firasdarwish/ore)
[![Go Reference](https://pkg.go.dev/badge/github.com/firasdarwish/ore.svg)](https://pkg.go.dev/github.com/firasdarwish/ore)
[![Mentioned in Awesome Go](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go?tab=readme-ov-file#dependency-injection)
[![codecov](https://codecov.io/gh/firasdarwish/ore/graph/badge.svg?token=ISZVCCYGCR)](https://codecov.io/gh/firasdarwish/ore)

**Ore** is a powerful and flexible **Dependency Injection (DI)** library for Go, designed to simplify complex
application structures while maintaining performance and modularity.

# 🌐 **[Full Documentation: ore.lilury.com](https://ore.lilury.com)**

---

## Key Features

### 1. **Flexible Lifetime Management**

- **Singletons**: Lifetime spans the entire application.
- **Scoped**: Lifetime is tied to a specific context.
- **Transient**: New instance created on every resolution.

---

### 2. **Alias Support**

- Link multiple implementations to the same interface.
- Easily resolve the preferred implementation or retrieve all registered options.

---

### 3. **Graceful Termination**

- **Application Termination**: Shutdown all resolved singletons implementing `Shutdowner` in proper dependency order.
- **Context Termination**: Dispose all resolved scoped instances implementing `Disposer` when the context ends.

---

### 4. **Placeholder Registration**

- Register incomplete dependencies at application setup.
- Provide runtime values dynamically when requests or events occur.

---

### 5. **Multiple Containers (Modules)**

- Create isolated dependency graphs for modular applications.
- Enforce module boundaries by separating service registrations and resolutions per container.

---

### 6. **Advanced Dependency Validation**

- Detect and prevent common pitfalls like:
    - **Missing Dependencies**: Ensure all resolvers are registered.
    - **Circular Dependencies**: Avoid infinite loops in dependency resolution.
    - **Lifetime Misalignment**: Catch improper lifetime dependencies (e.g., singleton depending on transient).

---

### 7. **Keyed and Keyless Registration**

- Support for multiple instances of the same type, scoped by keys.
- Easily resolve services using keys to manage module-specific configurations.

---

### 8. **Runtime Value Injection**

- Inject dependencies dynamically during runtime, tailored to request context or user data.

---

### 9. **Context-Based Dependency Resolution**

- Dependencies can be tied to a specific `context.Context`.
- Automatic cleanup of scoped services when the context ends.

---

## Installation

```bash
go get -u github.com/firasdarwish/ore
```

### Quick Example

```go
package main

import (
  "context"
  "fmt"
  "github.com/firasdarwish/ore"
)

// Define a service
type MyService struct {
  Message string
}

func main() {
  // Register a singleton service
  ore.RegisterSingleton(&MyService{Message: "Hello, Ore!"})
  
  // Resolve the service
  service, _ := ore.Get[*MyService](context.Background())
  fmt.Println(service.Message) // Output: Hello, Ore!
}
```

For complete usage examples and advanced configurations, visit the **[Ore Documentation](https://ore.lilury.com)**.

<br />

## 👤 Contributors

![Contributors](https://contrib.rocks/image?repo=firasdarwish/ore)

<br />

## Contributing

Feel free to contribute by opening issues, suggesting features, or submitting pull requests. We welcome your feedback
and contributions.

<br />

## License

This project is licensed under the MIT License - see the LICENSE file for details.