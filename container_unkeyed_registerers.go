package ore

import (
	"context"
)

// RegisterCreatorToContainer Registers a lazily initialized value to the given container using a `Creator[T]` interface
func RegisterCreatorToContainer[T any](con *Container, lifetime Lifetime, creator Creator[T]) {
	registerCreatorToContainer(con, lifetime, creator)
}

// RegisterSingletonToContainer Registers an eagerly instantiated singleton value to the given container.
// To register an eagerly instantiated scoped value use [ProvideScopedValueToContainer]
func RegisterSingletonToContainer[T comparable](con *Container, impl T) {
	registerSingletonToContainer(con, impl)
}

// RegisterFuncToContainer Registers a lazily initialized value to the given container using an `Initializer[T]` function signature
func RegisterFuncToContainer[T any](con *Container, lifetime Lifetime, initializer Initializer[T]) {
	registerFuncToContainer(con, lifetime, initializer)
}

// RegisterPlaceholderToContainer registers a future value with Scoped lifetime to the given container.
// This value will be injected in runtime using the [ProvideScopedValue] function.
// Resolving objects which depend on this value will panic if the value has not been provided.
// Placeholder with the same type and key can be registered only once.
func RegisterPlaceholderToContainer[T comparable](con *Container) {
	registerPlaceholderToContainer[T](con)
}

// ProvideScopedValueToContainer injects a concrete value into the given context.
// This value will be available only to the given container. And the container can only resolve this value if
// it has the matching (type and key's) Placeholder registered. Checkout the [RegisterPlaceholderToContainer] function for more info.
func ProvideScopedValueToContainer[T comparable](con *Container, ctx context.Context, value T) context.Context {
	return provideScopedValueToContainer(con, ctx, value)
}
