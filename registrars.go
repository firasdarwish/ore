package ore

import (
	"context"
	"fmt"
	"reflect"
	"time"
)

// RegisterLazyCreatorToContainer Registers a lazily initialized value to the given container using a `Creator[T]` interface
func RegisterLazyCreatorToContainer[T any](con *Container, lifetime Lifetime, creator Creator[T], key ...KeyStringer) {
	if creator == nil {
		panic(nilVal[T]())
	}

	e := serviceResolverImpl[T]{
		resolverMetadata: resolverMetadata{
			lifetime: lifetime,
		},
		creatorInstance: creator,
	}
	addResolver[T](con, e, key...)
}

// RegisterLazyCreator Registers a lazily initialized value using a `Creator[T]` interface
func RegisterLazyCreator[T any](lifetime Lifetime, creator Creator[T], key ...KeyStringer) {
	RegisterLazyCreatorToContainer[T](DefaultContainer, lifetime, creator, key...)
}

// RegisterEagerSingletonToContainer Registers an eagerly instantiated singleton value to the given container.
// To register an eagerly instantiated scoped value use [ProvideScopedValueToContainer]
func RegisterEagerSingletonToContainer[T comparable](con *Container, impl T, key ...KeyStringer) {
	if isNil[T](impl) {
		panic(nilVal[T]())
	}

	e := serviceResolverImpl[T]{
		resolverMetadata: resolverMetadata{
			lifetime: Singleton,
		},
		singletonConcrete: &concrete{
			value:          impl,
			lifetime:       Singleton,
			invocationTime: time.Now(),
		},
	}
	addResolver[T](con, e, key...)
}

// RegisterEagerSingleton Registers an eagerly instantiated singleton value
// To register an eagerly instantiated scoped value use [ProvideScopedValue]
func RegisterEagerSingleton[T comparable](impl T, key ...KeyStringer) {
	RegisterEagerSingletonToContainer[T](DefaultContainer, impl, key...)
}

// RegisterLazyFuncToContainer Registers a lazily initialized value to the given container using an `Initializer[T]` function signature
func RegisterLazyFuncToContainer[T any](con *Container, lifetime Lifetime, initializer Initializer[T], key ...KeyStringer) {
	if initializer == nil {
		panic(nilVal[T]())
	}

	e := serviceResolverImpl[T]{
		resolverMetadata: resolverMetadata{
			lifetime: lifetime,
		},
		anonymousInitializer: &initializer,
	}
	addResolver[T](con, e, key...)
}

// RegisterLazyFunc Registers a lazily initialized value using an `Initializer[T]` function signature
func RegisterLazyFunc[T any](lifetime Lifetime, initializer Initializer[T], key ...KeyStringer) {
	RegisterLazyFuncToContainer(DefaultContainer, lifetime, initializer, key...)
}

// RegisterAliasToContainer Registers an interface type to a concrete implementation in the given container.
// Allowing you to register the concrete implementation to the container and later get the interface from it.
func RegisterAliasToContainer[TInterface, TImpl any](con *Container) {
	interfaceType := reflect.TypeFor[TInterface]()
	implType := reflect.TypeFor[TImpl]()

	if !implType.Implements(interfaceType) {
		panic(fmt.Errorf("%s does not implements %s", implType, interfaceType))
	}

	addAliases[TInterface, TImpl](con)
}

// RegisterAlias Registers an interface type to a concrete implementation.
// Allowing you to register the concrete implementation to the default container and later get the interface from it.
func RegisterAlias[TInterface, TImpl any]() {
	RegisterAliasToContainer[TInterface, TImpl](DefaultContainer)
}

// RegisterPlaceHolderToContainer registers a future value with Scoped lifetime to the given container.
// This value will be injected in runtime using the [ProvideScopedValue] function.
// Resolving objects which depend on this value will panic if the value has not been provided.
// Placeholder with the same type and key can be registered only once.
func RegisterPlaceHolderToContainer[T comparable](con *Container, key ...KeyStringer) {
	e := serviceResolverImpl[T]{
		resolverMetadata: resolverMetadata{
			lifetime: Scoped,
		},
	}
	addResolver[T](con, e, key...)
}

// RegisterPlaceHolderToContainer registers a future value with Scoped lifetime.
// This value will be injected in runtime using the [ProvideScopedValue] function.
// Resolving objects which depend on this value will panic if the value has not been provided.
// Placeholder with the same type and key can be registered only once.
func RegisterPlaceHolder[T comparable](key ...KeyStringer) {
	RegisterPlaceHolderToContainer[T](DefaultContainer, key...)
}

// ProvideScopedValueToContainer injects a concrete value into the given context.
// This value will be available only to the given container. And the container can only resolve this value if
// it has the matching (type and key's) Placeholder registered. Checkout the [RegisterPlaceHolderToContainer] function for more info.
func ProvideScopedValueToContainer[T comparable](con *Container, ctx context.Context, value T, key ...KeyStringer) context.Context {
	concreteValue := &concrete{
		value:           value,
		lifetime:        Scoped,
		invocationTime:  time.Now(),
		invocationLevel: 0,
	}
	id := contextKey{
		containerID: con.containerID,
		typeID:      typeIdentifier[T](key...),
		resolverID:  placeHolderResolverID,
	}
	return addScopedConcreteToContext(ctx, id, concreteValue)
}

// ProvideScopedValue injects a concrete value into the given context.
// This value will be available only to the default container. And the container can only resolve this value if
// it has the matching (type and key's) Placeholder registered. Checkout the [RegisterPlaceHolder] function for more info.
func ProvideScopedValue[T comparable](ctx context.Context, value T, key ...KeyStringer) context.Context {
	return ProvideScopedValueToContainer[T](DefaultContainer, ctx, value, key...)
}
