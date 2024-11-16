package ore

import "context"

// RegisterCreator Registers a lazily initialized value using a `Creator[T]` interface
func RegisterCreator[T any](lifetime Lifetime, creator Creator[T]) {
	registerCreatorToContainer[T](DefaultContainer, lifetime, creator, nilKey)
}

// RegisterSingleton Registers an eagerly instantiated singleton value
// To register an eagerly instantiated scoped value use [ProvideScopedValue]
func RegisterSingleton[T any](impl T) {
	registerSingletonToContainer[T](DefaultContainer, impl, nilKey)
}

// RegisterFunc Registers a lazily initialized value using an `Initializer[T]` function signature
func RegisterFunc[T any](lifetime Lifetime, initializer Initializer[T]) {
	registerFuncToContainer(DefaultContainer, lifetime, initializer, nilKey)
}

// RegisterPlaceholder registers a future value with Scoped lifetime.
// This value will be injected in runtime using the [ProvideScopedValue] function.
// Resolving objects which depend on this value will panic if the value has not been provided.
// Placeholder with the same type and key can be registered only once.
func RegisterPlaceholder[T any]() {
	registerPlaceholderToContainer[T](DefaultContainer, nilKey)
}

// ProvideScopedValue injects a concrete value into the given context.
// This value will be available only to the default container. And the container can only resolve this value if
// it has the matching (type and key's) Placeholder registered. Checkout the [RegisterPlaceholder] function for more info.
func ProvideScopedValue[T any](ctx context.Context, value T) context.Context {
	return provideScopedValueToContainer(DefaultContainer, ctx, value, nilKey)
}
