package ore

import "context"

// RegisterCreator Registers a lazily initialized value using a `Creator[T]` interface
func RegisterCreator[T any](lifetime Lifetime, creator Creator[T], key ...KeyStringer) {
	registerCreatorToContainer[T](DefaultContainer, lifetime, creator, key...)
}

// RegisterSingleton Registers an eagerly instantiated singleton value
// To register an eagerly instantiated scoped value use [ProvideScopedValue]
func RegisterSingleton[T comparable](impl T, key ...KeyStringer) {
	registerSingletonToContainer[T](DefaultContainer, impl, key...)
}

// RegisterFunc Registers a lazily initialized value using an `Initializer[T]` function signature
func RegisterFunc[T any](lifetime Lifetime, initializer Initializer[T], key ...KeyStringer) {
	registerFuncToContainer(DefaultContainer, lifetime, initializer, key...)
}

// RegisterAlias Registers an interface type to a concrete implementation.
// Allowing you to register the concrete implementation to the default container and later get the interface from it.
func RegisterAlias[TInterface, TImpl any]() {
	registerAliasToContainer[TInterface, TImpl](DefaultContainer)
}

// RegisterPlaceholder registers a future value with Scoped lifetime.
// This value will be injected in runtime using the [ProvideScopedValue] function.
// Resolving objects which depend on this value will panic if the value has not been provided.
// Placeholder with the same type and key can be registered only once.
func RegisterPlaceholder[T comparable](key ...KeyStringer) {
	registerPlaceholderToContainer[T](DefaultContainer, key...)
}

// ProvideScopedValue injects a concrete value into the given context.
// This value will be available only to the default container. And the container can only resolve this value if
// it has the matching (type and key's) Placeholder registered. Checkout the [RegisterPlaceholder] function for more info.
func ProvideScopedValue[T comparable](ctx context.Context, value T, key ...KeyStringer) context.Context {
	return provideScopedValueToContainer(DefaultContainer, ctx, value, key...)
}
