package ore

import (
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

// RegisterEagerSingleton Registers an eagerly instantiated singleton value to the given container
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
// Allowing you to register the concrete implementation to the container and later get the interface from it.
func RegisterAlias[TInterface, TImpl any]() {
	RegisterAliasToContainer[TInterface, TImpl](DefaultContainer)
}
