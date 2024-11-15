package ore

import (
	"context"
)

// RegisterKeyedCreatorToContainer Registers a lazily initialized value to the given container using a `Creator[T]` interface
func RegisterKeyedCreatorToContainer[T any](con *Container, lifetime Lifetime, creator Creator[T], key KeyStringer) {
	if key == nil {
		panic(nilKey)
	}
	registerCreatorToContainer(con, lifetime, creator, key)
}

// RegisterKeyedSingletonToContainer Registers an eagerly instantiated singleton value to the given container.
// To register an eagerly instantiated scoped value use [ProvideScopedValueToContainer]
func RegisterKeyedSingletonToContainer[T comparable](con *Container, impl T, key KeyStringer) {
	if key == nil {
		panic(nilKey)
	}
	registerSingletonToContainer(con, impl, key)
}

// RegisterKeyedFuncToContainer Registers a lazily initialized value to the given container using an `Initializer[T]` function signature
func RegisterKeyedFuncToContainer[T any](con *Container, lifetime Lifetime, initializer Initializer[T], key KeyStringer) {
	if key == nil {
		panic(nilKey)
	}
	registerFuncToContainer(con, lifetime, initializer, key)
}

// RegisterKeyedPlaceholderToContainer registers a future value with Scoped lifetime to the given container.
// This value will be injected in runtime using the [ProvideScopedValue] function.
// Resolving objects which depend on this value will panic if the value has not been provided.
// Placeholder with the same type and key can be registered only once.
func RegisterKeyedPlaceholderToContainer[T comparable](con *Container, key KeyStringer) {
	if key == nil {
		panic(nilKey)
	}
	registerPlaceholderToContainer[T](con, key)
}

// ProvideKeyedScopedValueToContainer injects a concrete value into the given context.
// This value will be available only to the given container. And the container can only resolve this value if
// it has the matching (type and key's) Placeholder registered. Checkout the [RegisterPlaceholderToContainer] function for more info.
func ProvideKeyedScopedValueToContainer[T comparable](con *Container, ctx context.Context, value T, key KeyStringer) context.Context {
	if key == nil {
		panic(nilKey)
	}
	return provideScopedValueToContainer(con, ctx, value, key)
}
