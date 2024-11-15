package ore

import "context"

// RegisterKeyedCreator Registers a lazily initialized value using a `Creator[T]` interface
func RegisterKeyedCreator[T any](lifetime Lifetime, creator Creator[T], key KeyStringer) {
	if key == nil {
		panic(nilKey)
	}
	registerCreatorToContainer[T](DefaultContainer, lifetime, creator, key)
}

// RegisterKeyedSingleton Registers an eagerly instantiated singleton value
// To register an eagerly instantiated scoped value use [ProvideScopedValue]
func RegisterKeyedSingleton[T any](impl T, key KeyStringer) {
	if key == nil {
		panic(nilKey)
	}
	registerSingletonToContainer[T](DefaultContainer, impl, key)
}

// RegisterKeyedFunc Registers a lazily initialized value using an `Initializer[T]` function signature
func RegisterKeyedFunc[T any](lifetime Lifetime, initializer Initializer[T], key KeyStringer) {
	if key == nil {
		panic(nilKey)
	}
	registerFuncToContainer(DefaultContainer, lifetime, initializer, key)
}

// RegisterKeyedPlaceholder registers a future value with Scoped lifetime.
// This value will be injected in runtime using the [ProvideScopedValue] function.
// Resolving objects which depend on this value will panic if the value has not been provided.
// Placeholder with the same type and key can be registered only once.
func RegisterKeyedPlaceholder[T any](key KeyStringer) {
	if key == nil {
		panic(nilKey)
	}
	registerPlaceholderToContainer[T](DefaultContainer, key)
}

// ProvideKeyedScopedValue injects a concrete value into the given context.
// This value will be available only to the default container. And the container can only resolve this value if
// it has the matching (type and key's) Placeholder registered. Checkout the [RegisterPlaceholder] function for more info.
func ProvideKeyedScopedValue[T any](ctx context.Context, value T, key KeyStringer) context.Context {
	if key == nil {
		panic(nilKey)
	}
	return provideScopedValueToContainer(DefaultContainer, ctx, value, key)
}
